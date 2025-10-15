package did

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"io"

	c "github.com/wang900115/LCA/crypto"
	"github.com/wang900115/LCA/pkg/util/encode"
)

type PeerKeyPair struct {
	EdPublic  ed25519.PublicKey
	EdPrivate ed25519.PrivateKey
	XPublic   *ecdh.PublicKey
	XPrivate  *ecdh.PrivateKey
}

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

func (k *PeerKeyPair) generateDID() string {
	header := []byte{0xed, 0x01}
	payload := append(header, k.EdPublic...)
	return "did:key:z" + encode.Base58Encode(payload)
}
