package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAPIKeyMiddleware_ValidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	validKey := "test-api-key-12345"
	router.Use(APIKeyMiddleware(validKey))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", validKey)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestAPIKeyMiddleware_InvalidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	validKey := "test-api-key-12345"
	router.Use(APIKeyMiddleware(validKey))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or missing API key")
}

func TestAPIKeyMiddleware_MissingKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	validKey := "test-api-key-12345"
	router.Use(APIKeyMiddleware(validKey))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or missing API key")
}

func TestAPIKeyMiddleware_QueryParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	validKey := "test-api-key-12345"
	router.Use(APIKeyMiddleware(validKey))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test?apikey=test-api-key-12345", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestAPIKeyMiddleware_HealthCheckBypass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	validKey := "test-api-key-12345"
	router.Use(APIKeyMiddleware(validKey))
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	req, _ := http.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}

func TestValidateAPIKeyFormat(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		want   bool
	}{
		{"valid key", "abcdefghijklmnop", true},
		{"valid key with numbers", "abc123DEF4567890", true}, // 16 chars
		{"valid key with dashes", "abc-def-ghi-jkl-mn", true}, // 20 chars
		{"valid key with underscores", "abc_def_ghi_jkl_mn", true}, // 21 chars
		{"valid base64 key", "dGVzdC1rZXktZm9yLWJhc2U2NC1lbmNvZGluZw==", true}, // base64 can have = and /
		{"valid base64 with slash", "dGVzdC9rZXkvd2l0aC9zbGFzaA==", true},
		{"too short", "short", false},
		{"empty", "", false},
		{"invalid character", "abcdefghijklmnop@", false},
		{"valid long key", "abcdefghijklmnopqrstuvwxyz123456", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAPIKeyFormat(tt.apiKey)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateSecureAPIKey(t *testing.T) {
	key1, err := GenerateSecureAPIKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, key1)
	assert.True(t, ValidateAPIKeyFormat(key1))

	// Generate another key to ensure uniqueness
	key2, err := GenerateSecureAPIKey()
	assert.NoError(t, err)
	assert.NotEqual(t, key1, key2, "Generated keys should be unique")
}
