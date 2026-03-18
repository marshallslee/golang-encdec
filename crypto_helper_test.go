package main

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "testSecretKey"
	plainText := "Hello, World!"

	encrypted, err := Encrypt(key, plainText)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plainText {
		t.Errorf("expected %q, got %q", plainText, decrypted)
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := "testSecretKey"
	plainText := "Hello, World!"

	encrypted, err := Encrypt(key, plainText)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt("wrongKey", encrypted)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}
