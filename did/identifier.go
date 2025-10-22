package did

import (
	"crypto/ed25519"
	"crypto/rand"
	"time"

	"github.com/btcsuite/btcutil/base58"
	crypto "github.com/wang900115/LCA/crypt"
)

var (
	DIDVersion = 1
)

// IdentifierDID defines the interface for a Identifier DID.
type IdentifierDID interface {
	Addr() string
	Document() *Document
	SignDocument() ([]byte, error)
	Sign(data []byte) ([]byte, error)
}

// Metadata holds metadata for a DID.
type Metadata struct {
	Controller string
	Version    int
}

// ServiceEndpoint represents a service endpoint in the DID
type ServiceEndpoint struct {
	ID   string
	Type string
	URL  string
}

// DIDIdentifier represents a Decentralized Identifier.
type DIDIdentifier struct {
	ID       string
	Address  string
	KeyPair  KeyPair
	Metadata Metadata
	Services []ServiceEndpoint
}

// NewDID creates a new IdentifierDID instance.
func NewDIDIdentifier(services []ServiceEndpoint) IdentifierDID {
	var did DIDIdentifier
	var err error
	did.KeyPair, err = NewPeerKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}
	did.ID = did.KeyPair.GenerateDID()
	did.Metadata = Metadata{
		Controller: did.ID,
		Version:    DIDVersion,
	}
	did.Address = did.KeyPair.GenerateAddr()
	did.Services = services
	return &did
}

// Addr returns the DID address.
func (d *DIDIdentifier) Addr() string {
	return d.Address
}

// Document converts the DID to a DID Document.
func (d *DIDIdentifier) Document() *Document {
	return NewDocument(*d, time.Now().Unix())
}

// SignDocument signs the DID Document.
func (d *DIDIdentifier) SignDocument() ([]byte, error) {
	doc := d.Document()
	data, err := doc.JSONMarshal()
	if err != nil {
		return nil, err
	}
	signature, err := d.KeyPair.SignData(data)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// SignMessage signs a message using the DID's key pair.
func (d *DIDIdentifier) Sign(data []byte) ([]byte, error) {
	signature, err := d.KeyPair.SignData(data)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// extract extracts the Ed25519 public key from the DID Document.
func extract(doc *Document) (ed25519.PublicKey, error) {
	for _, vm := range doc.VerificationMethod {
		if vm.Type == VerificationType {
			return base58.Decode(vm.PublicKeyBase58), nil
		}
	}
	return nil, crypto.ErrED25519PublicKeyMissing
}
