package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/listenarr/listenarr/internal/models"
)

func TestSearchAudiobooks(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{
		Title:    "Test Book",
		AuthorID: author.ID,
		ISBN:     "1234567890",
	}
	db.Create(&book)

	router := gin.New()
	router.GET("/api/v1/search", server.searchAudiobooks)

	t.Run("Search with query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/search?q=Test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Search without query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/search", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Search by ISBN", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/search?q=1234567890", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
