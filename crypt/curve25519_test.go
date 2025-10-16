package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestED25519(t *testing.T) {
	data := []byte("MESSAGE")
	pubKey, priKey, err := ED25519GenerateKey(rand.Reader)
	assert.Nil(t, err)
	t.Logf("ED25519 Private Key: %v \nED25519 Public Key: %v", priKey, pubKey)
	signature, err := ED25519Sign(priKey, data)
	assert.Nil(t, err)
	ok, err := ED25519Verify(pubKey, data, signature)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestSharedKey(t *testing.T) {
	pubKeyAlice, priKeyAlice, err := ED25519GenerateKey(rand.Reader)
	assert.Nil(t, err)
	t.Logf("Alice ED25519 Private Key: %v \nAlice ED25519 Public Key: %v\n", priKeyAlice, pubKeyAlice)
	xpubKeyAlice, xpriKeyAlice, err := X25519GenerateKey(rand.Reader)
	assert.Nil(t, err)
	t.Logf("Alice X25519 Private Key:  %v \nAlice X25519 Public Key:  %v\n", xpriKeyAlice.Bytes(), xpubKeyAlice.Bytes())
	signature1, err := SignX25519PublicKey(priKeyAlice, xpubKeyAlice)
	assert.Nil(t, err)
	ok, err := VerifyX25519PublicKeySignature(pubKeyAlice, xpubKeyAlice, signature1)
	assert.Nil(t, err)
	assert.True(t, ok)

	pubKeyCindy, priKeyCindy, err := ED25519GenerateKey(rand.Reader)
	assert.Nil(t, err)
	t.Logf("Cindy ED25519 Private Key: %v \nCindy ED25519 Public Key: %v\n", priKeyCindy, pubKeyCindy)
	xpubKeyCindy, xpriKeyCindy, err := X25519GenerateKey(rand.Reader)
	assert.Nil(t, err)
	t.Logf("Cindy X25519 Private Key:  %v \nCindy X25519 Public Key:  %v\n", xpriKeyCindy.Bytes(), xpubKeyCindy.Bytes())
	signature2, err := SignX25519PublicKey(priKeyCindy, xpubKeyCindy)
	assert.Nil(t, err)
	ok, err = VerifyX25519PublicKeySignature(pubKeyCindy, xpubKeyCindy, signature2)
	assert.Nil(t, err)
	assert.True(t, ok)

	sharedKey, err := ComputeX25519SharedKey(xpriKeyAlice, xpubKeyCindy)
	assert.Nil(t, err)
	sharedKey2, err := ComputeX25519SharedKey(xpriKeyCindy, xpubKeyAlice)
	assert.Nil(t, err)
	assert.Equal(t, sharedKey, sharedKey2, "Alice and Cindy should derive the same shared key")
}
