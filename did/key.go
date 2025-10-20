package did

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/sha3"
	"io"

	c "github.com/wang900115/LCA/crypt"
	"github.com/wang900115/LCA/pkg/util/encode"
)

// PeerKeyPair holds the key pairs for a Peer DID.
type PeerKeyPair struct {
	EdPublic  ed25519.PublicKey
	EdPrivate ed25519.PrivateKey
	XPublic   *ecdh.PublicKey
	XPrivate  *ecdh.PrivateKey
}

// NewPeerKeyPair generates a new PeerKeyPair.
func NewPeerKeyPair(r io.Reader) (*PeerKeyPair, error) {
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

// generateDID generates a DID from the Ed25519 public key.
func (k *PeerKeyPair) generateDID() string {
	header := []byte{0xed, 0x01}
	payload := append(header, k.EdPublic...)
	return "did:key:z" + encode.Base58Encode(payload)
}

// generateAddr generates an address from the ED25519 public key.
func (k *PeerKeyPair) generateAddr() string {
	hash := sha3.Sum256(k.EdPublic)
	addrBytes := hash[:]
	return "addr:" + encode.Base58Encode(addrBytes)
}
