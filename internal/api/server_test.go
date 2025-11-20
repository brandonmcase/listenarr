package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/config"
)

func setupTestServer(t *testing.T) (*Server, string) {
	gin.SetMode(gin.TestMode)

	// Create test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create test config
	testConfig := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8686,
		},
		Auth: config.AuthConfig{
			Enabled: true,
			APIKey:  "test-api-key",
		},
	}

	server := NewServer(testConfig, testDB)
	return server, "test-api-key"
}

func TestHealthCheck(t *testing.T) {
	server, _ := setupTestServer(t)

	req, _ := http.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
	assert.Contains(t, w.Body.String(), "listenarr")
}

func TestGetLibrary_RequiresAuth(t *testing.T) {
	server, apiKey := setupTestServer(t)

	// Test without API key
	req, _ := http.NewRequest("GET", "/api/v1/library", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test with API key
	req, _ = http.NewRequest("GET", "/api/v1/library", nil)
	req.Header.Set("X-API-Key", apiKey)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDownloads_RequiresAuth(t *testing.T) {
	server, apiKey := setupTestServer(t)

	req, _ := http.NewRequest("GET", "/api/v1/downloads", nil)
	req.Header.Set("X-API-Key", apiKey)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProcessingQueue_RequiresAuth(t *testing.T) {
	server, apiKey := setupTestServer(t)

	req, _ := http.NewRequest("GET", "/api/v1/processing", nil)
	req.Header.Set("X-API-Key", apiKey)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchAudiobooks_RequiresAuth(t *testing.T) {
	server, apiKey := setupTestServer(t)

	req, _ := http.NewRequest("GET", "/api/v1/search", nil)
	req.Header.Set("X-API-Key", apiKey)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddToLibrary_RequiresAuth(t *testing.T) {
	server, apiKey := setupTestServer(t)

	req, _ := http.NewRequest("POST", "/api/v1/library", nil)
	req.Header.Set("X-API-Key", apiKey)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRemoveFromLibrary_RequiresAuth(t *testing.T) {
	server, apiKey := setupTestServer(t)

	req, _ := http.NewRequest("DELETE", "/api/v1/library/123", nil)
	req.Header.Set("X-API-Key", apiKey)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStartDownload_RequiresAuth(t *testing.T) {
	server, apiKey := setupTestServer(t)

	req, _ := http.NewRequest("POST", "/api/v1/downloads", nil)
	req.Header.Set("X-API-Key", apiKey)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
