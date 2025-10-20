package network

import (
	"github.com/wang900115/LCA/did"
)

type HandShakeContent struct {
	DIDDocument *did.DIDDocument `json:"did_document"`
	Signature   []byte           `json:"signature"` // Signature of the DID Document
	Challenge   []byte           `json:"challenge"` // Random challenge for replay protection
	Version     string           `json:"version"`
}

func NewHandShakeContent(didDoc *did.DIDDocument, signature, challenge []byte, version string) *HandShakeContent {
	return &HandShakeContent{
		DIDDocument: didDoc,
		Signature:   signature,
		Challenge:   challenge,
		Version:     version,
	}
}
