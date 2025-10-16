package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRC32(t *testing.T) {
	data := []byte("Hello World")
	n := CRC32(data)
	t.Logf("CRC32 CheckSum: %d\n", n)
	assert.NotZero(t, n)
	ok := VerifyCRC32(data, n)
	assert.True(t, ok)
}

func TestCRC64(t *testing.T) {
	data := []byte("Hello World")
	n := CRC64(data)
	t.Logf("CRC64 CheckSum: %d\n", n)
	assert.NotZero(t, n)
	ok := VerifyCRC64(data, n)
	assert.True(t, ok)
}
