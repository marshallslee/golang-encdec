package main

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	keyObtentionIterations = 1000
	saltSize               = 8
)

// deriveKeyAndIV derives DES key and IV from password and salt using MD5
// This implements the OpenSSL/PKCS#5 PBEWithMD5AndDES key derivation.
func deriveKeyAndIV(password []byte, salt []byte, iterations int) (key []byte, iv []byte) {
	hash := md5.New()
	hash.Write(password)
	hash.Write(salt)
	result := hash.Sum(nil)

	for i := 1; i < iterations; i++ {
		hash.Reset()
		hash.Write(result)
		result = hash.Sum(nil)
	}

	// MD5 produces 16 bytes: first 8 = DES key, next 8 = IV
	key = result[:8]
	iv = result[8:]
	return
}

// pkcs5Pad pads data to block size using PKCS5/PKCS7 padding.
func pkcs5Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

// pkcs5Unpad removes PKCS5/PKCS7 padding.
func pkcs5Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding > des.BlockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-padding], nil
}

// Encrypt encrypts plainText using PBEWithMD5AndDES with the given secretKey.
// Output format: base64(salt + ciphertext), compatible with Jasypt.
func Encrypt(secretKey, plainText string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	password := []byte(secretKey)
	key, iv := deriveKeyAndIV(password, salt, keyObtentionIterations)

	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create DES cipher: %w", err)
	}

	padded := pkcs5Pad([]byte(plainText), des.BlockSize)
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	// Jasypt format: salt (8 bytes) + ciphertext
	result := append(salt, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts encryptedText (base64 encoded) using PBEWithMD5AndDES with the given secretKey.
// Expects Jasypt format: base64(salt + ciphertext).
func Decrypt(secretKey, encryptedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(data) < saltSize+des.BlockSize {
		return "", fmt.Errorf("encrypted data too short")
	}

	salt := data[:saltSize]
	ciphertext := data[saltSize:]

	password := []byte(secretKey)
	key, iv := deriveKeyAndIV(password, salt, keyObtentionIterations)

	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create DES cipher: %w", err)
	}

	if len(ciphertext)%des.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	unpadded, err := pkcs5Unpad(plaintext)
	if err != nil {
		return "", fmt.Errorf("failed to unpad: %w", err)
	}

	return string(unpadded), nil
}
