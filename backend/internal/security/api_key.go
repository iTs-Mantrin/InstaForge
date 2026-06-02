package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// APIKeyGenerator handles API key creation and hashing.
type APIKeyGenerator struct{}

// NewAPIKeyGenerator creates a new API key generator.
func NewAPIKeyGenerator() *APIKeyGenerator {
	return &APIKeyGenerator{}
}

// Generate creates a new API key (if-key format) and its hash.
func (g *APIKeyGenerator) Generate() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("failed to generate key: %w", err)
	}

	rawKey := fmt.Sprintf("if_%s", hex.EncodeToString(bytes))
	hash := g.Hash(rawKey)

	return rawKey, hash, nil
}

// Hash computes a SHA-256 hash of an API key for secure storage.
func (g *APIKeyGenerator) Hash(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

// Validate compares a raw key against a stored hash.
func (g *APIKeyGenerator) Validate(rawKey, storedHash string) bool {
	return g.Hash(rawKey) == storedHash
}

// GenerateShortToken creates a short random token (for CSRF, etc.).
func GenerateShortToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
