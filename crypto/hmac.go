package crypto

import (
	"crypto/hmac"
	"crypto/sha256"

	"golang.org/x/crypto/sha3"
)

// generate HMAC-SHA2-256
func HMACSignSHA2(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// generate HMAC-SHA3-256
func HMACSignSHA3(key, data []byte) []byte {
	mac := hmac.New(sha3.New256, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// verify HMAC-SHA2-256
func HMACVerifySHA2(key, data, expected []byte) bool {
	mac := HMACSignSHA2(key, data)
	return hmac.Equal(mac, expected)
}

// verify HMAC-SHA3-256
func HMACVerifySHA3(key, data, expected []byte) bool {
	mac := HMACSignSHA3(key, data)
	return hmac.Equal(mac, expected)
}
