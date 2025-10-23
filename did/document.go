package did

import (
	"encoding/json"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

const (
	VerificationID   = "#keys-1"
	KeyAgreementID   = "#keys-2"
	VerificationType = "Ed25519VerificationKey2018"
	KeyAgreementType = "X25519KeyAgreementKey2019"
	DIDContext       = "https://www.w3.org/ns/did/v1"
)

// Document represents the DID's Document structure following W3C DID spec.
type Document struct {
	Context              []string             `json:"@context"`
	ID                   string               `json:"id"`
	VerificationMethod   []VerificationMethod `json:"verificationMethod"`
	Authentication       []string             `json:"authentication"`
	AssertionMethod      []string             `json:"assertionMethod"`
	KeyAgreement         []string             `json:"keyAgreement"`
	CapabilityInvocation []string             `json:"capabilityInvocation"`
	CapabilityDelegation []string             `json:"capabilityDelegation"`
	Service              []ServiceEndpoint    `json:"service,omitempty"`
	Created              string               `json:"created,omitempty"`
	Updated              string               `json:"updated,omitempty"`
}

// VerificationMethod represents a verification method in the DID Document.
type VerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

// ServiceEndpoint represents a service endpoint in the DID Document.
type ServiceEndpoint struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	ServiceEndpoint interface{} `json:"serviceEndpoint"`
}

// NewDocument creates a new DID Document based on the provided DIDIdentifier and creation time.
func NewDocument(did DIDIdentifier, createdAt time.Time) *Document {
	vmId := composeID(did.ID, VerificationID)
	kaId := composeID(did.ID, KeyAgreementID)

	doc := &Document{
		Context: []string{
			DIDContext,
		},
		ID: did.ID,
		VerificationMethod: []VerificationMethod{
			newVerificationMethod(vmId, did.Metadata.Controller, VerificationType, did.KeyPair.GetEd25519PublicKey()),
			newKeyAgreementMethod(kaId, did.Metadata.Controller, KeyAgreementType, did.KeyPair.GetX25519PublicKey()),
		},
		Authentication:       []string{vmId},
		AssertionMethod:      []string{vmId},
		KeyAgreement:         []string{kaId},
		CapabilityInvocation: []string{vmId},
		CapabilityDelegation: []string{vmId},
		Service:              convertToW3CServices(did.Services),
		Created:              createdAt.Format(time.RFC3339),
		Updated:              createdAt.Format(time.RFC3339),
	}
	return doc
}

func newVerificationMethod(id, controller, keyType string, publicKey []byte) VerificationMethod {
	return VerificationMethod{
		ID:                 id,
		Type:               keyType,
		Controller:         controller,
		PublicKeyMultibase: "z" + base58.Encode(publicKey), // multibase with base58btc prefix
	}
}

func newKeyAgreementMethod(id, controller, keyType string, publicKey []byte) VerificationMethod {
	return VerificationMethod{
		ID:                 id,
		Type:               keyType,
		Controller:         controller,
		PublicKeyMultibase: "z" + base58.Encode(publicKey), // multibase with base58btc prefix
	}
}

func convertToW3CServices(services []ServiceEndpoint) []ServiceEndpoint {
	w3cServices := make([]ServiceEndpoint, len(services))
	for i, service := range services {
		w3cServices[i] = ServiceEndpoint{
			ID:              service.ID,
			Type:            service.Type,
			ServiceEndpoint: service.ServiceEndpoint,
		}
	}
	return w3cServices
}

const (
	// New
	Ed25519VerificationKey2020 = "Ed25519VerificationKey2020"
	X25519KeyAgreementKey2020  = "X25519KeyAgreementKey2020"

	// Old
	Ed25519VerificationKey2018 = "Ed25519VerificationKey2018"
	X25519KeyAgreementKey2019  = "X25519KeyAgreementKey2019"
)

// NewDocumentWithNewStandards creates a new DID Document following the latest W3C DID standards.
func NewDocumentWithNewStandards(did DIDIdentifier, createdAt time.Time) *Document {
	vmId := composeID(did.ID, VerificationID)
	kaId := composeID(did.ID, KeyAgreementID)

	return &Document{
		Context: []string{
			"https://www.w3.org/ns/did/v1",
			"https://w3id.org/security/suites/ed25519-2020/v1",
			"https://w3id.org/security/suites/x25519-2020/v1",
		},
		ID: did.ID,
		VerificationMethod: []VerificationMethod{
			{
				ID:                 vmId,
				Type:               Ed25519VerificationKey2020,
				Controller:         did.Metadata.Controller,
				PublicKeyMultibase: "z" + base58.Encode(did.KeyPair.GetEd25519PublicKey()),
			},
			{
				ID:                 kaId,
				Type:               X25519KeyAgreementKey2020,
				Controller:         did.Metadata.Controller,
				PublicKeyMultibase: "z" + base58.Encode(did.KeyPair.GetX25519PublicKey()),
			},
		},
		Authentication:       []string{vmId},
		AssertionMethod:      []string{vmId},
		KeyAgreement:         []string{kaId},
		CapabilityInvocation: []string{vmId},
		CapabilityDelegation: []string{vmId},
		Service:              convertToW3CServices(did.Services),
		Created:              createdAt.Format(time.RFC3339),
		Updated:              createdAt.Format(time.RFC3339),
	}
}

func (d *Document) JSONMarshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Document) JSONUnmarshal(data []byte) error {
	return json.Unmarshal(data, d)
}

func composeID(id, fragment string) string {
	return id + fragment
}
