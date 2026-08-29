package main

import (
	"encoding/base64"
	"testing"
)

func TestValidateKEK(t *testing.T) {
	t.Run("Empty Key", func(t *testing.T) {
		_, err := validateKEK("")
		if err == nil {
			t.Error("Expected error for empty key, got nil")
		}
	})

	t.Run("Valid Base64 Key", func(t *testing.T) {
		// 32 bytes valid base64 key
		key := base64.StdEncoding.EncodeToString([]byte("12345678901234567890123456789012"))
		decoded, err := validateKEK(key)
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if len(decoded) != 32 {
			t.Errorf("Expected 32 bytes, got %d", len(decoded))
		}
	})

	t.Run("Invalid Key Length", func(t *testing.T) {
		// 16 bytes valid base64 key
		key := base64.StdEncoding.EncodeToString([]byte("1234567890123456"))
		_, err := validateKEK(key)
		if err == nil {
			t.Error("Expected error for invalid key length, got nil")
		}
	})

	t.Run("Plaintext 32 bytes Key", func(t *testing.T) {
		key := "12345678901234567890123456789012"
		decoded, err := validateKEK(key)
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if len(decoded) != 32 {
			t.Errorf("Expected 32 bytes, got %d", len(decoded))
		}
	})
}
