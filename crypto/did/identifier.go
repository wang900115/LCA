package did

import (
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/wang900115/LCA/crypto"
)

var (
	DIDVersion = 1
)

type PeerDID interface {
	ToDocument() *DIDDocument
	SignDocument() ([]byte, error)
	VerifyDocument([]byte) (bool, error)
}

type DIDMetadata struct {
	CreatedAt  int64
	Controller string
	Version    int
}

type ServiceEndpoint struct {
	ID   string
	Type string
	URL  string
}

type DID struct {
	DID      string
	KeyPair  *PeerKeyPair
	Metadata DIDMetadata
	Services []ServiceEndpoint
}

func NewPeerDID(services []ServiceEndpoint) (*DID, error) {
	var did DID
	pair, err := NewPeerKeyPair(rand.Reader)
	if err != nil {
		return nil, err
	}
	did.DID = pair.generateDID()
	did.Metadata = DIDMetadata{
		CreatedAt:  time.Now().UTC().Unix(),
		Controller: did.DID,
		Version:    DIDVersion,
	}
	did.Services = services
	return &did, nil
}

type DIDDocument struct {
	Context            string               `json:"@context"`
	ID                 string               `json:"id"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
	KeyAgreement       []VerificationMethod `json:"keyAgreement"`
	Service            []ServiceEndpoint    `json:"service,omitempty"`
	CreatedAt          int64                `json:"created,omitempty"`
	Version            int                  `json:"version,omitempty"`
}

type VerificationMethod struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	Controller      string `json:"controller"`
	PublicKeyBase58 string `json:"publicKeyBase58"`
}

func (d *DID) ToDocument() *DIDDocument {
	id := d.DID
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
		CreatedAt: d.Metadata.CreatedAt,
		Version:   d.Metadata.Version,
	}
}

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

func (d *DID) VerifyDocument(signature []byte) (bool, error) {
	doc := d.ToDocument()
	data, err := json.Marshal(doc)
	if err != nil {
		return false, err
	}
	return crypto.ED25519Verify(d.KeyPair.EdPublic, data, signature)
}
