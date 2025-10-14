package crypto

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"errors"
	"io"
)

// generate x25519 pub/pri pair
func X25519GenerateKey(r io.Reader) (*ecdh.PublicKey, *ecdh.PrivateKey, error) {
	privateKey, err := ecdh.X25519().GenerateKey(r)
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.PublicKey()
	return publicKey, privateKey, nil
}

// generate ed25519 pub/pri pair
func ED25519GenerateKey(r io.Reader) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(r)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

// using private key to signature data
func ED25519Sign(privateKey ed25519.PrivateKey, data []byte) (signature []byte) {
	return ed25519.Sign(privateKey, data)
}

// using public key to verify data is correct
func ED25519Verify(publicKey ed25519.PublicKey, data []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, data, signature)
}

// using ed25519 private key sign x25519 public key
func SignX25519PublicKey(edPriv ed25519.PrivateKey, xPub *ecdh.PublicKey) []byte {
	return ED25519Sign(edPriv, xPub.Bytes())
}

// using ed25519 public key to verify x25519 public signature
func VerifyX25519PublicKeySignature(edPub ed25519.PublicKey, xPub *ecdh.PublicKey, signature []byte) bool {
	return ED25519Verify(edPub, xPub.Bytes(), signature)
}

// compute shared X25519 secret key
func ComputeX25519SharedKey(privateKey *ecdh.PrivateKey, peerPublicKey *ecdh.PublicKey) ([]byte, error) {
	if privateKey == nil || peerPublicKey == nil {
		return nil, errors.New("invalid key")
	}
	shared, err := privateKey.ECDH(peerPublicKey)
	if err != nil {
		return nil, err
	}
	return shared, nil
}
