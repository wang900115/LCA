package crypto

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
)

func TestHMACSHA256(t *testing.T) {
	Key := []byte("PRIVATE")
	data := []byte("MESSAGE")
	signature := HMACSign(sha256.New, Key, data)
	assert.NotEmpty(t, signature)
	ok := HMACVerify(sha256.New, Key, data, signature)
	assert.True(t, ok)
}

func TestHMACSHA3(t *testing.T) {
	Key := []byte("PRIVATE")
	data := []byte("MESSAGE")
	signature := HMACSign(sha3.New256, Key, data)
	assert.NotEmpty(t, signature)
	ok := HMACVerify(sha3.New256, Key, data, signature)
	assert.True(t, ok)
}
