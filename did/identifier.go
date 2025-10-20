package did

import (
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/btcsuite/btcutil/base58"
	crypto "github.com/wang900115/LCA/crypt"
)

var (
	DIDVersion = 1
)

// PeerDID defines the interface for a Peer DID.
type PeerDID interface {
	DID() DID
	ToDocument() *DIDDocument
	SignDocument() ([]byte, error)
	VerifyDocument([]byte) (bool, error)
}

// DIDMetadata holds metadata for a DID.
type Metadata struct {
	Controller string
	Version    int
}

// ServiceEndpointType defines the type of service endpoint.
type ServiceEndpoint struct {
	ID   string
	Type string
	URL  string
}

// DID represents a Decentralized Identifier.
type DID struct {
	ID       string
	Address  string
	KeyPair  *PeerKeyPair
	Metadata Metadata
	Services []ServiceEndpoint
}

// NewDID creates a new PeerDID instance.
func NewDID(services []ServiceEndpoint) PeerDID {
	var did DID
	pair, err := NewPeerKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}
	did.ID = pair.generateDID()
	did.Metadata = Metadata{
		Controller: did.ID,
		Version:    DIDVersion,
	}
	did.Address = pair.generateAddr()
	did.KeyPair = pair
	did.Services = services
	return &did
}

// DIDDocument represents a DID Document.
type DIDDocument struct {
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

// DID returns the DID information.
func (d *DID) DID() DID {
	return *d
}

// ToDocument converts the DID to a DID Document.
func (d *DID) ToDocument() *DIDDocument {
	id := d.ID
	return &DIDDocument{
		Context: "https://www.w3.org/ns/did/v1",
		ID:      id,
		VerificationMethod: []VerificationMethod{{
			ID:              id + "#keys-1",
			Type:            "Ed25519VerificationKey2018",
			Controller:      d.Metadata.Controller,
			PublicKeyBase58: base58.Encode(d.KeyPair.EdPublic),
		}},
		KeyAgreement: []VerificationMethod{{
			ID:              id + "#keys-2",
			Type:            "X25519KeyAgreementKey2019",
			Controller:      d.Metadata.Controller,
			PublicKeyBase58: base58.Encode(d.KeyPair.XPublic.Bytes()),
		}},
		Service:   d.Services,
		CreatedAt: time.Now().UTC().Unix(),
		Version:   d.Metadata.Version,
	}
}

// SignDocument signs the DID Document.
func (d *DID) SignDocument() ([]byte, error) {
	doc := d.ToDocument()
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	signature, err := crypto.ED25519Sign(d.KeyPair.EdPrivate, data)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// VerifyDocument verifies the signature of the DID Document.
func (d *DID) VerifyDocument(signature []byte) (bool, error) {
	doc := d.ToDocument()
	data, err := json.Marshal(doc)
	if err != nil {
		return false, err
	}
	return crypto.ED25519Verify(d.KeyPair.EdPublic, data, signature)
}
