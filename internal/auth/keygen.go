package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecureAPIKey generates a cryptographically secure random API key
func GenerateSecureAPIKey() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to base64 URL-safe string
	apiKey := base64.URLEncoding.EncodeToString(bytes)
	return apiKey, nil
}

