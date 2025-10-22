package did

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/sha3"
	"io"

	c "github.com/wang900115/LCA/crypt"
	"github.com/wang900115/LCA/pkg/util/encode"
)

// KeyPair defines the interface for key pair operations.
type KeyPair interface {
	GenerateDID() string
	GenerateAddr() string
	GetEd25519PublicKey() []byte
	GetX25519PublicKey() []byte
	SignData(data []byte) ([]byte, error)
	VerifyData(data []byte, signature []byte) (bool, error)
	Shake(peerPublicKey *ecdh.PublicKey) ([]byte, error)
	Unshake(peerPublicKey *ecdh.PublicKey, signature []byte, peerEdPublicKey ed25519.PublicKey) ([]byte, error)
}

// PeerKeyPair holds the key pairs for a Peer DID.
type PeerKeyPair struct {
	EdPublic  ed25519.PublicKey
	EdPrivate ed25519.PrivateKey
	XPublic   *ecdh.PublicKey
	XPrivate  *ecdh.PrivateKey
}

// NewPeerKeyPair generates a new PeerKeyPair.
func NewPeerKeyPair(r io.Reader) (KeyPair, error) {
	edPub, edPriv, err := c.ED25519GenerateKey(r)
	if err != nil {
		return nil, err
	}
	xPub, xPriv, err := c.X25519GenerateKey(r)
	if err != nil {
		return nil, err
	}
	newKeyPair := &PeerKeyPair{
		EdPublic:  edPub,
		EdPrivate: edPriv,
		XPublic:   xPub,
		XPrivate:  xPriv,
	}
	return newKeyPair, nil
}

// GetEd25519PublicKey returns the Ed25519 public key.
func (k *PeerKeyPair) GetEd25519PublicKey() []byte {
	return k.EdPublic
}

// GetX25519PublicKey returns the X25519 public key.
func (k *PeerKeyPair) GetX25519PublicKey() []byte {
	return k.XPublic.Bytes()
}

// using Ed25519 private key sign data
func (k *PeerKeyPair) SignData(data []byte) ([]byte, error) {
	return c.ED25519Sign(k.EdPrivate, data)
}

// using Ed25519 public key verify data signature
func (k *PeerKeyPair) VerifyData(data []byte, signature []byte) (bool, error) {
	return c.ED25519Verify(k.EdPublic, data, signature)
}

// GenerateDID generates a DID from the Ed25519 public key.
func (k *PeerKeyPair) GenerateDID() string {
	header := []byte{0xed, 0x01}
	payload := append(header, k.EdPublic...)
	return "did:key:z" + encode.Base58Encode(payload)
}

// GenerateAddr generates an address from the ED25519 public key.
func (k *PeerKeyPair) GenerateAddr() string {
	hash := sha3.Sum256(k.EdPublic)
	addrBytes := hash[:]
	return "addr:" + encode.Base58Encode(addrBytes)
}

// Shake sign own X25519 public key using own Ed25519 private key
func (k *PeerKeyPair) Shake(peerPublicKey *ecdh.PublicKey) ([]byte, error) {
	return c.SignX25519PublicKey(k.EdPrivate, k.XPublic)
}

// Unshake verify peer's X25519 public key signature using peer's Ed25519 public key and generate shared secret
func (k *PeerKeyPair) Unshake(peerPublicKey *ecdh.PublicKey, signature []byte, peerEdPublicKey ed25519.PublicKey) ([]byte, error) {
	valid, err := c.ED25519Verify(peerEdPublicKey, peerPublicKey.Bytes(), signature)
	if err != nil || !valid {
		return nil, err
	}
	return c.ComputeX25519SharedKey(k.XPrivate, peerPublicKey)
}
