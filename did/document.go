package did

import (
	"encoding/json"

	"github.com/btcsuite/btcutil/base58"
)

const (
	VerificationID   = "#keys-1"
	KeyAgreementID   = "#keys-2"
	VerificationType = "Ed25519VerificationKey2018"
	KeyAgreementType = "X25519KeyAgreementKey2019"
)

// Document represents the DID's Document structure.
type Document struct {
	Context            string               `json:"@context"`
	ID                 string               `json:"id"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
	KeyAgreement       []VerificationMethod `json:"keyAgreement"`
	Service            []ServiceEndpoint    `json:"service,omitempty"`
	CreatedAt          int64                `json:"created,omitempty"`
	Version            int                  `json:"version,omitempty"`
}

// VerificationMethod represents a verification method in the DID Document.
type VerificationMethod struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	Controller      string `json:"controller"`
	PublicKeyBase58 string `json:"publicKeyBase58"`
}

// NewDocument converts the DID to a DID Document.
func NewDocument(did DIDIdentifier, createdAt int64) *Document {
	doc := &Document{
		Context:            "https://www.w3.org/ns/did/v1",
		ID:                 did.ID,
		VerificationMethod: []VerificationMethod{newVerificationMethod(did.ID, did.Metadata.Controller, VerificationType, did.KeyPair.GetEd25519PublicKey())},
		KeyAgreement:       []VerificationMethod{newKeyAgreementMethod(did.ID, did.Metadata.Controller, KeyAgreementType, did.KeyPair.GetX25519PublicKey())},
		Service:            did.Services,
		CreatedAt:          createdAt,
		Version:            did.Metadata.Version,
	}
	return doc
}

func newVerificationMethod(id, controller string, keyType string, publicKey []byte) VerificationMethod {
	return VerificationMethod{
		ID:              composeID(id, VerificationID),
		Type:            keyType,
		Controller:      controller,
		PublicKeyBase58: base58.Encode(publicKey),
	}
}

func newKeyAgreementMethod(id, controller string, keyType string, publicKey []byte) VerificationMethod {
	return VerificationMethod{
		ID:              composeID(id, KeyAgreementID),
		Type:            keyType,
		Controller:      controller,
		PublicKeyBase58: base58.Encode(publicKey),
	}
}

func composeID(id, fragment string) string {
	return id + fragment
}

func (d *Document) JSONMarshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Document) JSONUnmarshal(data []byte) error {
	return json.Unmarshal(data, d)
}
