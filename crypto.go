package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
)

const (
	// BADKEY Error for bad key length
	BADKEY = "Unsupported key size"
)

// AESCipher Cipher object for AES GCM mode
type AESCipher struct {
	ciphertext string
	nonce      string
}

// Check key length support
func allowedLength(length int) bool {
	KEYSIZES := map[int]bool{
		16: true,
		32: true,
	}
	return KEYSIZES[length]
}

// Generate AES key from password and context
func generateKey(key string, context string, length int) ([]byte, error) {
	if !allowedLength(length) {
		return nil, errors.New(BADKEY)
	}
	hash := sha512.Sum512([]byte(key + context))
	return hash[:length], nil
}

// Encrypt string using AES in GCM mode, return object with base64 encoding or ciphertext and nonce
// Supported key sizes: 128bits, 256bits
func encryptAESGCM(plain string, key string, context string, length int) (AESCipher, error) {
	var res AESCipher
	genkey, err := generateKey(key, context, length)
	if err != nil {
		return res, err
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return res, err
	}
	block, err := aes.NewCipher(genkey)
	if err != nil {
		return res, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return res, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plain), nil)
	res = AESCipher{
		base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(nonce),
	}
	return res, nil
}

// Decrypt payload encrypted with AES in GCM mode, accepts object with base64 encoding or ciphertext and nonce
// Supported key sizes: 128bits, 256bits
func decryptAESGCM(cipherobj *AESCipher, key string, context string, length int) (string, error) {
	genkey, err := generateKey(key, context, length)
	if err != nil {
		return "", err
	}
	cipherbytes, err1 := base64.StdEncoding.DecodeString(cipherobj.ciphertext)
	nonce, err2 := base64.StdEncoding.DecodeString(cipherobj.nonce)
	if err1 != nil || err2 != nil {
		return "", errors.New(err1.Error() + "; " + err2.Error())
	}
	block, err := aes.NewCipher(genkey)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plain, err := aesgcm.Open(nil, nonce, cipherbytes, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// Evaluate SHA256 digest of a string in base64 encoding
func sha256sum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}
