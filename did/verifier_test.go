package did

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDIDIdentifierVerifierIntegration(t *testing.T) {
	t.Run("Full DiD create and verification", func(t *testing.T) {
		services := []ServiceEndpoint{
			{
				ID:              "messaging",
				Type:            "MessagingService",
				ServiceEndpoint: "https://example.com/messaging",
			},
		}
		peerDID := NewDIDIdentifier(services)
		assert.NotNil(t, peerDID)

		doc := peerDID.Document()
		assert.NotNil(t, doc)
		assert.NotEmpty(t, doc.ID)
		assert.NotEmpty(t, doc.VerificationMethod)

		signature, err := peerDID.SignDocument()
		require.NoError(t, err)
		assert.NotEmpty(t, signature)

		verifier := NewDefaultDIDVerifier()
		isValid, err := verifier.VerifyDocument(doc, signature)
		assert.NoError(t, err)
		assert.True(t, isValid)

		stats := verifier.GetStats()
		assert.Equal(t, int64(1), stats.TotalVerifications)
		assert.Equal(t, int64(1), stats.SuccessfulVerifications)
		assert.Equal(t, int64(0), stats.FailedVerifications)
	})

	t.Run("Cross-DID verification", func(t *testing.T) {
		verifier := NewDefaultDIDVerifier()

		did1 := NewDIDIdentifier([]ServiceEndpoint{})
		did2 := NewDIDIdentifier([]ServiceEndpoint{})
		did3 := NewDIDIdentifier([]ServiceEndpoint{})

		doc1 := did1.Document()
		sig1, _ := did1.SignDocument()

		doc2 := did2.Document()
		sig2, _ := did2.SignDocument()

		doc3 := did3.Document()
		sig3, _ := did3.SignDocument()

		valid1, err1 := verifier.VerifyDocument(doc1, sig1)
		assert.NoError(t, err1)
		assert.True(t, valid1)

		valid2, err2 := verifier.VerifyDocument(doc2, sig2)
		assert.NoError(t, err2)
		assert.True(t, valid2)

		valid3, err3 := verifier.VerifyDocument(doc3, sig3)
		assert.NoError(t, err3)
		assert.True(t, valid3)

		stats := verifier.GetStats()
		assert.Equal(t, int64(3), stats.TotalVerifications)
		assert.Equal(t, int64(3), stats.SuccessfulVerifications)
		assert.Equal(t, int64(0), stats.FailedVerifications)
	})

	t.Run("Cross-DID signature verification failure", func(t *testing.T) {
		verifier := NewDefaultDIDVerifier()

		did1 := NewDIDIdentifier([]ServiceEndpoint{})
		did2 := NewDIDIdentifier([]ServiceEndpoint{})

		doc1 := did1.Document()
		sig2, _ := did2.SignDocument()

		isValid, err := verifier.VerifyDocument(doc1, sig2)
		assert.NoError(t, err)
		assert.False(t, isValid)

		stats := verifier.GetStats()
		assert.Equal(t, int64(1), stats.TotalVerifications)
		assert.Equal(t, int64(0), stats.SuccessfulVerifications)
		assert.Equal(t, int64(1), stats.FailedVerifications)
	})
}

func TestDIDLifecycleWithVerifier(t *testing.T) {
	t.Run("DID Lifecycle Management", func(t *testing.T) {
		verifier := NewDefaultDIDVerifier()

		did := NewDIDIdentifier([]ServiceEndpoint{
			{ID: "service1", Type: "Type1", ServiceEndpoint: "https://service1.com"},
		})

		doc := did.Document()
		signature, _ := did.SignDocument()

		isValid, err := verifier.VerifyDocument(doc, signature)
		assert.NoError(t, err)
		assert.True(t, isValid)

		updatedDID := NewDIDIdentifier([]ServiceEndpoint{
			{ID: "service1", Type: "Type1", ServiceEndpoint: "https://service1.com"},
			{ID: "service2", Type: "Type2", ServiceEndpoint: "https://service2.com"},
		})

		updatedDoc := updatedDID.Document()
		updatedSignature, _ := updatedDID.SignDocument()

		isValidUpdated, errUpdated := verifier.VerifyDocument(updatedDoc, updatedSignature)
		assert.NoError(t, errUpdated)
		assert.True(t, isValidUpdated)

		stats := verifier.GetStats()
		assert.Equal(t, int64(2), stats.TotalVerifications)
		assert.Equal(t, int64(2), stats.SuccessfulVerifications)
	})
}

func TestTrustedRootScenarios(t *testing.T) {
	t.Run("Trusted Root Verification", func(t *testing.T) {
		config := VerifierConfig{
			EnableCache:        true,
			CacheTTL:           30 * time.Minute,
			MaxCacheSize:       1000,
			ValidateTimestamp:  true,
			TimestampTolerance: 5 * time.Minute,
			RequireTrustedRoot: true,
		}
		verifier := NewDIDVerifier(config)

		authorityDID := NewDIDIdentifier([]ServiceEndpoint{})
		authorityDoc := authorityDID.Document()

		verifier.AddTrustedRoot(authorityDoc.ID)

		authoritySignature, _ := authorityDID.SignDocument()
		isValid, err := verifier.VerifyDocument(authorityDoc, authoritySignature)
		assert.NoError(t, err)
		assert.True(t, isValid)

		untrustedDID := NewDIDIdentifier([]ServiceEndpoint{})
		untrustedDoc := untrustedDID.Document()
		untrustedSignature, _ := untrustedDID.SignDocument()

		isValidUntrusted, errUntrusted := verifier.VerifyDocument(untrustedDoc, untrustedSignature)
		assert.Error(t, errUntrusted)
		assert.False(t, isValidUntrusted)
		assert.Contains(t, errUntrusted.Error(), ErrDocNotController.Error())
	})
}

func TestCacheEffectivenessWithRealDIDs(t *testing.T) {
	t.Run("Cache Effectiveness", func(t *testing.T) {
		config := VerifierConfig{
			EnableCache:        true,
			CacheTTL:           1 * time.Hour,
			MaxCacheSize:       100,
			ValidateTimestamp:  false,
			RequireTrustedRoot: false,
		}
		verifier := NewDIDVerifier(config)

		did := NewDIDIdentifier([]ServiceEndpoint{})
		doc := did.Document()
		signature, _ := did.SignDocument()

		start := time.Now()
		isValid1, err1 := verifier.VerifyDocument(doc, signature)
		duration1 := time.Since(start)

		assert.NoError(t, err1)
		assert.True(t, isValid1)

		start = time.Now()
		isValid2, err2 := verifier.VerifyDocument(doc, signature)
		duration2 := time.Since(start)

		assert.NoError(t, err2)
		assert.True(t, isValid2)

		t.Logf("First verification: %v, Second verification: %v", duration1, duration2)

		stats := verifier.GetStats()
		assert.Equal(t, int64(2), stats.TotalVerifications)
		assert.Equal(t, int64(1), stats.CacheMisses)
		assert.Equal(t, int64(1), stats.CacheHits)
		assert.Equal(t, int64(2), stats.SuccessfulVerifications)
	})
}

func TestTimestampValidationWithRealDIDs(t *testing.T) {
	t.Run("Timestamp Validation", func(t *testing.T) {

		config := VerifierConfig{
			EnableCache:        false,
			ValidateTimestamp:  true,
			TimestampTolerance: 1 * time.Minute,
			RequireTrustedRoot: false,
		}
		verifier := NewDIDVerifier(config)

		did := NewDIDIdentifier([]ServiceEndpoint{})
		doc := did.Document()

		doc.Created = time.Now().Add(-5 * time.Minute).Format(time.RFC3339)

		signature, _ := did.SignDocument()

		isValid, err := verifier.VerifyDocument(doc, signature)
		assert.Error(t, err)
		assert.False(t, isValid)
		assert.Contains(t, err.Error(), ErrTimestampInvalid.Error())
		doc.Created = time.Now().Format(time.RFC3339)
		signature, _ = did.SignDocument()
		isValid, err = verifier.VerifyDocument(doc, signature)
		assert.NoError(t, err)
		assert.True(t, isValid)
	})
}

func TestPeerToPeerVerification(t *testing.T) {
	t.Run("P2P Verification Scenario", func(t *testing.T) {
		// Simulate mutual verification between two peers
		verifier := NewDefaultDIDVerifier()

		// Alice and Bob each create a DID
		aliceDID := NewDIDIdentifier([]ServiceEndpoint{
			{ID: "alice-messaging", Type: "MessagingService", ServiceEndpoint: "https://alice.com/msg"},
		})

		bobDID := NewDIDIdentifier([]ServiceEndpoint{
			{ID: "bob-messaging", Type: "MessagingService", ServiceEndpoint: "https://bob.com/msg"},
		})

		aliceDoc := aliceDID.Document()
		aliceSignature, _ := aliceDID.SignDocument()

		bobDoc := bobDID.Document()
		bobSignature, _ := bobDID.SignDocument()

		aliceValidFromBob, err := verifier.VerifyDocument(aliceDoc, aliceSignature)
		assert.NoError(t, err)
		assert.True(t, aliceValidFromBob)

		bobValidFromAlice, err := verifier.VerifyDocument(bobDoc, bobSignature)
		assert.NoError(t, err)
		assert.True(t, bobValidFromAlice)

		assert.NotEqual(t, aliceDoc.ID, bobDoc.ID)

		aliceKey, _ := extract(aliceDoc)
		bobKey, _ := extract(bobDoc)
		assert.NotEqual(t, aliceKey, bobKey)

		t.Logf("Alice DID: %s", aliceDoc.ID)
		t.Logf("Bob DID: %s", bobDoc.ID)
	})
}
