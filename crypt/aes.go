/*
	AES Module (Outer Layer Encryption in Transport Layer)
	------------------------------------------------------
	This module is used to encrypt packets at the outermost layer
	of the transport protocol.

	Purpose:
	1. Protect packet content from eavesdropping.
	2. Prevent tampering during transmission.

	Key (key) Explanation:
	------------------------------------------------------------
	The `key` parameter is a "shared key" between sender and receiver.
	This shared key is usually established through a secure key exchange
	mechanism (e.g., Diffie-Hellman, ECDH, or RSA key exchange) and
	should NOT be hard-coded.
*/

package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// KCS#7 Padding: pad plaintext to a multiple of blockSize
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS#7 Unpadding: remove padding after decryption
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// AESCBCEncrypt: encrypt plaintext using AES-CBC with shared key
func AESCBCEncrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plainText = pkcs7Padding(plainText, block.BlockSize())
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)
	return cipherText, nil
}

// AESCBCDecrypt: decrypt ciphertext using AES-CBC with shared key
func AESCBCDecrypto(cipherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText = pkcs7UnPadding(cipherText)
	return cipherText, nil
}

// derivedKey = KDF(sharedKey || senderPublicKey || receiverPublicKey)
func DeriveAESKey(sharedKey []byte, senderPub, receiverPub ed25519.PublicKey) ([]byte, []byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}
	info := append(senderPub, receiverPub...)
	hkdf := hkdf.New(sha256.New, sharedKey, salt, info)
	key := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, key); err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}
