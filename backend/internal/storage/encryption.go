package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// EncryptionKey is the master key for encryption (loaded from env or generated)
var encryptionKey []byte

// InitEncryption initializes the encryption key from environment, file, or generates one
func InitEncryption(dataDir string) error {
	// Try to load from environment variable first
	keyStr := os.Getenv("QUOTIO_ENCRYPTION_KEY")
	if keyStr != "" {
		// Key should be base64-encoded 32-byte key (256-bit)
		key, err := base64.StdEncoding.DecodeString(keyStr)
		if err != nil {
			return fmt.Errorf("failed to decode encryption key from environment: %w", err)
		}
		if len(key) != 32 {
			return fmt.Errorf("encryption key must be 32 bytes, got %d", len(key))
		}
		encryptionKey = key
		return nil
	}

	// Try to load from file
	keyFilePath := filepath.Join(dataDir, ".encryption.key")
	keyData, err := os.ReadFile(keyFilePath)
	if err == nil {
		// File exists, decode the key
		key, err := base64.StdEncoding.DecodeString(string(keyData))
		if err != nil {
			return fmt.Errorf("failed to decode encryption key from file: %w", err)
		}
		if len(key) != 32 {
			return fmt.Errorf("encryption key must be 32 bytes, got %d", len(key))
		}
		encryptionKey = key
		return nil
	}

	// Generate a new key if not provided
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Persist the key to file
	encodedKey := base64.StdEncoding.EncodeToString(key)
	if err := os.WriteFile(keyFilePath, []byte(encodedKey), 0600); err != nil {
		return fmt.Errorf("failed to persist encryption key: %w", err)
	}

	encryptionKey = key
	return nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	if len(encryptionKey) != 32 {
		return "", fmt.Errorf("encryption key not initialized")
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce (should be 12 bytes for GCM)
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := aead.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	if len(encryptionKey) != 32 {
		return "", fmt.Errorf("encryption key not initialized")
	}

	// Decode from base64
	encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		// If decoding fails, assume plaintext (for backwards compatibility)
		return ciphertext, nil
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(encrypted) < nonceSize {
		// If too short, assume plaintext (for backwards compatibility)
		return ciphertext, nil
	}

	nonce := encrypted[:nonceSize]
	ciphertext_only := encrypted[nonceSize:]

	// Decrypt
	plaintext, err := aead.Open(nil, nonce, ciphertext_only, nil)
	if err != nil {
		// If decryption fails, assume plaintext (for backwards compatibility with legacy data)
		return ciphertext, nil
	}

	return string(plaintext), nil
}
