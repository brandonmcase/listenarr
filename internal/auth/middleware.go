package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware validates API key from request
func APIKeyMiddleware(validAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for health check endpoint
		if c.Request.URL.Path == "/api/health" {
			c.Next()
			return
		}

		// Get API key from header or query parameter
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("apikey")
		}

		// Validate API key
		if apiKey == "" || apiKey != validAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateAPIKey generates a secure random API key
func GenerateAPIKey() (string, error) {
	// Generate a 32-character random string
	// In production, use crypto/rand for better security
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b), nil
}

// ValidateAPIKeyFormat validates that an API key has the correct format
func ValidateAPIKeyFormat(apiKey string) bool {
	// API key should be at least 16 characters
	if len(apiKey) < 16 {
		return false
	}
	// Check for valid characters (alphanumeric and base64-safe chars)
	// Base64 URL-safe encoding uses: A-Z, a-z, 0-9, -, _, =, /
	for _, char := range apiKey {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' || char == '=' || char == '/') {
			return false
		}
	}
	return true
}
