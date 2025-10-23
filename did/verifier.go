package did

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"sync"
	"time"

	crypto "github.com/wang900115/LCA/crypt"
)

var (
	ErrExpired            = errors.New("verification result expired")
	ErrNotFound           = errors.New("verification result not found")
	ErrMissingCreatedAt   = errors.New("document missing createdAt field")
	ErrMissingTrustedRoot = errors.New("document not signed by trusted root")
	ErrDocNotController   = errors.New("document controller not in trusted roots")
	ErrTimestampInvalid   = errors.New("document timestamp is invalid")
)

// VerifierDID defines the interface for verifying a DID Document.
type VerifierDID interface {
	VerifyDocument(doc *Document, signature []byte) (bool, error)
	GetStats() VerificationStats
	AddTrustedRoot(did string)
	ClearCache()
}

// DIDVerifier implements the VerifierDID interface.
type DIDVerifier struct {
	verifiedCache map[string]*VerificationResult
	cacheMutex    sync.RWMutex
	config        VerifierConfig
	stats         VerificationStats
	trustedRoots  []string
	statsMutex    sync.Mutex
}

// VerifierConfig holds configuration for the DID verifier.
type VerifierConfig struct {
	EnableCache        bool
	CacheTTL           time.Duration
	MaxCacheSize       int
	ValidateTimestamp  bool
	TimestampTolerance time.Duration
	RequireTrustedRoot bool
}

// VerificationResult holds the result of a DID verification attempt.
type VerificationResult struct {
	IsValid    bool
	DID        string
	VerifiedAt time.Time
	ExpiresAt  time.Time
	ErrorMsg   string
	Signature  []byte
	PublicKey  ed25519.PublicKey
}

// VerificationStats holds statistics about the verification process.
type VerificationStats struct {
	TotalVerifications      int64
	SuccessfulVerifications int64
	FailedVerifications     int64
	CacheHits               int64
	CacheMisses             int64
}

func NewDIDVerifier(config VerifierConfig) VerifierDID {
	return &DIDVerifier{
		verifiedCache: make(map[string]*VerificationResult),
		config:        config,
		stats:         VerificationStats{},
	}
}

func NewDefaultDIDVerifier() VerifierDID {
	config := VerifierConfig{
		EnableCache:        true,
		CacheTTL:           30 * time.Minute,
		MaxCacheSize:       1000,
		ValidateTimestamp:  true,
		TimestampTolerance: 5 * time.Minute,
		RequireTrustedRoot: false,
	}
	return NewDIDVerifier(config)
}

// VerifyDocument verifies the DID Document using the provided signature.
func (v *DIDVerifier) VerifyDocument(doc *Document, signature []byte) (bool, error) {
	v.statsMutex.Lock()
	v.stats.TotalVerifications++
	v.statsMutex.Unlock()

	cacheKey := v.cacheKey(doc, signature)
	if v.config.EnableCache {
		if cachedResult, err := v.getCachedResult(cacheKey); err == nil {
			v.statsMutex.Lock()
			v.stats.CacheHits++
			v.statsMutex.Unlock()

			if cachedResult.IsValid {
				v.recordSuccess()
			} else {
				v.recordFailure()
			}
			return cachedResult.IsValid, nil
		} else {
			v.statsMutex.Lock()
			v.stats.CacheMisses++
			v.statsMutex.Unlock()
		}
	}
	if v.config.ValidateTimestamp {
		if err := v.validateTimestamp(doc); err != nil {
			v.recordFailure()
			return false, err
		}
	}
	if v.config.RequireTrustedRoot {
		if err := v.validateTrustedRoot(doc); err != nil {
			v.recordFailure()
			return false, err
		}
	}
	publicKey, err := extract(doc)
	if err != nil {
		return false, err
	}
	isValid, err := verifyDocumentWithKey(doc, signature, publicKey)
	if err != nil {
		v.recordFailure()
		return false, err
	}
	result := &VerificationResult{
		IsValid:    isValid,
		DID:        doc.ID,
		VerifiedAt: time.Now(),
		ExpiresAt:  time.Now().Add(v.config.CacheTTL),
		Signature:  signature,
		PublicKey:  publicKey,
	}
	if !isValid {
		result.ErrorMsg = "signature verification failed"
		v.recordFailure()
	} else {
		v.recordSuccess()
	}
	if v.config.EnableCache {
		v.cacheResult(cacheKey, result)
	}
	return isValid, nil
}

// AddTrustedRoot adds a trusted root DID to the verifier.
func (v *DIDVerifier) AddTrustedRoot(rootDID string) {
	v.cacheMutex.Lock()
	defer v.cacheMutex.Unlock()
	v.trustedRoots = append(v.trustedRoots, rootDID)
}

// GetStats returns the current verification statistics.
func (v *DIDVerifier) GetStats() VerificationStats {
	v.statsMutex.Lock()
	defer v.statsMutex.Unlock()
	return v.stats
}

// ClearCache clears the verification result cache.
func (v *DIDVerifier) ClearCache() {
	v.cacheMutex.Lock()
	defer v.cacheMutex.Unlock()
	v.verifiedCache = make(map[string]*VerificationResult)
}

func (v *DIDVerifier) cacheResult(key string, result *VerificationResult) {
	v.cacheMutex.Lock()
	defer v.cacheMutex.Unlock()
	if len(v.verifiedCache) >= v.config.MaxCacheSize {
		v.cleanupCache()
	}
	v.verifiedCache[key] = result
}

func (v *DIDVerifier) cacheKey(doc *Document, signature []byte) string {
	return fmt.Sprintf("%s:%x", doc.ID, signature)
}

func (v *DIDVerifier) getCachedResult(key string) (*VerificationResult, error) {
	v.cacheMutex.RLock()
	defer v.cacheMutex.RUnlock()
	result, exists := v.verifiedCache[key]
	if !exists {
		return nil, ErrNotFound
	}
	if time.Now().After(result.ExpiresAt) {
		return nil, ErrExpired
	}
	return result, nil
}

func (v *DIDVerifier) validateTimestamp(doc *Document) error {
	if doc.Created == "" {
		return ErrMissingCreatedAt
	}
	docTime, err := time.Parse(time.RFC3339, doc.Created)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	if now.Sub(docTime) > v.config.TimestampTolerance {
		return ErrTimestampInvalid
	}
	if docTime.Sub(now) > v.config.TimestampTolerance {
		return ErrTimestampInvalid
	}
	return nil
}

func (v *DIDVerifier) validateTrustedRoot(doc *Document) error {
	if len(v.trustedRoots) == 0 {
		return ErrMissingTrustedRoot
	}
	for _, vm := range doc.VerificationMethod {
		for _, root := range v.trustedRoots {
			if vm.Controller == root {
				return nil
			}
		}
	}
	for _, root := range v.trustedRoots {
		if doc.ID == root {
			return nil
		}
	}
	return ErrDocNotController
}

func (v *DIDVerifier) cleanupCache() {
	now := time.Now()
	for key, result := range v.verifiedCache {
		if now.After(result.ExpiresAt) {
			delete(v.verifiedCache, key)
		}
	}
}

func (v *DIDVerifier) recordSuccess() {
	v.statsMutex.Lock()
	defer v.statsMutex.Unlock()
	v.stats.SuccessfulVerifications++
}

func (v *DIDVerifier) recordFailure() {
	v.statsMutex.Lock()
	defer v.statsMutex.Unlock()
	v.stats.FailedVerifications++
}

// verifyDocumentWithKey verifies the DID Document using the provided public key.
func verifyDocumentWithKey(doc *Document, signature []byte, publicKey ed25519.PublicKey) (bool, error) {
	data, err := doc.JSONMarshal()
	if err != nil {
		return false, err
	}
	return crypto.ED25519Verify(publicKey, data, signature)
}
