/*
Hash-based Message Authentication Code with x52119 shared key
is used in the transmit integrity layer.
Encrypts packet payload.
*/
package crypto

import (
	"crypto/hmac"
	"hash"
)

// generate HMAC
func HMACSign(f func() hash.Hash, key, data []byte) []byte {
	mac := hmac.New(f, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// verify HMAC
func HMACVerify(f func() hash.Hash, key, data, expected []byte) bool {
	mac := HMACSign(f, key, data)
	return hmac.Equal(mac, expected)
}
