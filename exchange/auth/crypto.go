package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// deriveKey derives a 32-byte AES key from password using SHA-256.
func deriveKey(password string) []byte {
	h := sha256.Sum256([]byte(password))
	return h[:]
}

// EncryptPrivateKey encrypts privKeyHex with AES-256-GCM using a key derived
// from password via SHA-256.  The returned string is hex-encoded
// nonce (12 bytes) + ciphertext.
func EncryptPrivateKey(privKeyHex, password string) (string, error) {
	key := deriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(privKeyHex), nil)

	combined := append(nonce, ciphertext...) //nolint:gocritic
	return hex.EncodeToString(combined), nil
}

// DecryptPrivateKey reverses EncryptPrivateKey, returning the original
// private key hex string.
func DecryptPrivateKey(encryptedHex, password string) (string, error) {
	combined, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", fmt.Errorf("decode hex: %w", err)
	}

	key := deriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(combined) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := combined[:nonceSize], combined[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}
