/*
	Curve25519 Module (Key Management and Signatures)
	------------------------------------------------------------
	This package provides cryptographic utilities for key management,
	signatures, and shared secret computation using Ed25519 and X25519.

	Main Features:
	1. Generate Ed25519 public/private key pairs (for signing/verification)
	2. Generate X25519 public/private key pairs (for ECDH shared secret)
	3. Sign arbitrary data with Ed25519 private key
	4. Verify signatures with Ed25519 public key
	5. Sign X25519 public keys for authenticity
	6. Compute X25519 shared secret key for symmetric encryption
	7. Must* variants panic on error for convenience in trusted contexts

	Usage Notes:
	- X25519 keys are used for Diffie-Hellman key exchange.
	- Ed25519 keys are used for digital signatures.
	- Errors are returned if keys are missing or signatures are absent.
	- Must* functions panic instead of returning errors; use only when errors are impossible or should halt execution.
*/

package crypto

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"io"
)

type errCrypto struct{ msg string }

func (e errCrypto) Error() string { return e.msg }

var (
	ErrX25519PrivateKeyMissing      = &errCrypto{"x25519 private key is missing"}
	ErrX25519RemotePublicKeyMissing = &errCrypto{"x25519 remote public key is missing"}
	ErrED25519PrivateKeyMissing     = &errCrypto{"ed25519 private key is missing"}
	ErrED25519PublicKeyMissing      = &errCrypto{"ed25519 public key is missing"}
	ErrSignatureMissing             = &errCrypto{"signature missing"}
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
func ED25519Sign(privateKey ed25519.PrivateKey, data []byte) ([]byte, error) {
	if len(privateKey) > 0 {
		return ed25519.Sign(privateKey, data), nil
	}
	return nil, ErrED25519PrivateKeyMissing
}

// using public key to verify data is correct
func ED25519Verify(publicKey ed25519.PublicKey, data []byte, signature []byte) (bool, error) {
	lenSignature := len(signature)
	if lenSignature == 0 {
		return false, ErrSignatureMissing
	}
	if len(publicKey) > 0 {
		return ed25519.Verify(publicKey, data, signature), nil
	}
	return false, ErrED25519PublicKeyMissing
}

// using ed25519 private key sign x25519 public key
func SignX25519PublicKey(edPriv ed25519.PrivateKey, xPub *ecdh.PublicKey) ([]byte, error) {
	return ED25519Sign(edPriv, xPub.Bytes())
}

// using ed25519 private key sign x25519 public key if not will panic err
func MustSignX25519PublicKey(edPriv ed25519.PrivateKey, xPub *ecdh.PublicKey) []byte {
	signature, err := ED25519Sign(edPriv, xPub.Bytes())
	if err != nil {
		panic(err)
	}
	return signature
}

// using ed25519 public key to verify x25519 public signature
func VerifyX25519PublicKeySignature(edPub ed25519.PublicKey, xPub *ecdh.PublicKey, signature []byte) (bool, error) {
	return ED25519Verify(edPub, xPub.Bytes(), signature)
}

// using ed25519 public key to verify x25519 public signature if not will panic err
func MustVerifyX25519PublicKeySignature(edPub ed25519.PublicKey, xPub *ecdh.PublicKey, signature []byte) bool {
	res, err := ED25519Verify(edPub, xPub.Bytes(), signature)
	if err != nil {
		panic(err)
	}
	return res
}

// compute shared X25519 secret key
func ComputeX25519SharedKey(privateKey *ecdh.PrivateKey, peerPublicKey *ecdh.PublicKey) ([]byte, error) {
	if privateKey == nil {
		return nil, ErrX25519PrivateKeyMissing
	}
	if peerPublicKey == nil {
		return nil, ErrX25519RemotePublicKeyMissing
	}
	shared, err := privateKey.ECDH(peerPublicKey)
	if err != nil {
		return nil, err
	}
	return shared, nil
}

// compute shared X25519 secret key if not will panic err
func MustComputeX25519SharedKey(privateKey *ecdh.PrivateKey, peerPublicKey *ecdh.PublicKey) []byte {
	sharedKey, err := ComputeX25519SharedKey(privateKey, peerPublicKey)
	if err != nil {
		panic(err)
	}
	return sharedKey
}
