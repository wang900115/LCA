package crypto

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
)

func TestHMACSHA256(t *testing.T) {
	privKey := []byte("PRIVATE")
	data := []byte("MESSAGE")
	signature := HMACSign(sha256.New, privKey, data)
	assert.NotEmpty(t, signature)
	ok := HMACVerify(sha256.New, privKey, data, signature)
	assert.True(t, ok)
}

func TestHMACSHA3(t *testing.T) {
	privKey := []byte("PRIVATE")
	data := []byte("MESSAGE")
	signature := HMACSign(sha3.New256, privKey, data)
	assert.NotEmpty(t, signature)
	ok := HMACVerify(sha3.New256, privKey, data, signature)
	assert.True(t, ok)
}
