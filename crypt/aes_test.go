package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCBC(t *testing.T) {
	plaintText := []byte("AESTEST")
	key := []byte("1234567890ABCDEF")

	cipherText, err := AESCBCEncrypt(plaintText, key)
	t.Logf("cipherText: %v", cipherText)
	assert.Nil(t, err)

	resText, err := AESCBCDecrypto(cipherText, key)
	t.Logf("plaintText: %v", resText)
	assert.Nil(t, err)

	assert.Equal(t, plaintText, resText)
}
