package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/config"
	"github.com/listenarr/listenarr/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate models
	err = db.AutoMigrate(
		&models.Author{},
		&models.Series{},
		&models.Book{},
		&models.Audiobook{},
		&models.LibraryItem{},
		&models.Release{},
		&models.Download{},
		&models.ProcessingTask{},
	)
	assert.NoError(t, err)

	return db
}

func setupLibraryTestServer(db *gorm.DB) *Server {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8686,
		},
		Auth: config.AuthConfig{
			Enabled: false, // Disable auth for tests
		},
	}
	return NewServer(cfg, db)
}

func TestGetLibrary(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	item := models.LibraryItem{
		BookID:    book.ID,
		Status:    models.LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	db.Create(&item)

	router := gin.New()
	router.GET("/api/v1/library", server.getLibrary)

	t.Run("Get library items", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/library", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
		assert.Equal(t, 1, response.Pagination.Total)
	})

	t.Run("Get library items with pagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/library?page=1&limit=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 10, response.Pagination.Limit)
	})

	t.Run("Get library items with status filter", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/library?status=wanted", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Get library items with author filter", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/library?author_id=1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetLibraryItem(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	item := models.LibraryItem{
		BookID:    book.ID,
		Status:    models.LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	db.Create(&item)

	router := gin.New()
	router.GET("/api/v1/library/:id", server.getLibraryItem)

	t.Run("Get existing library item", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/library/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	t.Run("Get non-existent library item", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/library/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Get library item with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/library/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAddToLibrary(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	router := gin.New()
	router.POST("/api/v1/library", server.addToLibrary)

	t.Run("Add book to library", func(t *testing.T) {
		reqBody := AddToLibraryRequest{
			Title:      "New Book",
			AuthorName: "New Author",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/library", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)

		// Verify item was created
		var count int64
		db.Model(&models.LibraryItem{}).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Add book with missing required fields", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "Book without author",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/library", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Add book with series", func(t *testing.T) {
		seriesName := "Test Series"
		reqBody := AddToLibraryRequest{
			Title:          "Series Book",
			AuthorName:     "Series Author",
			SeriesName:     &seriesName,
			SeriesPosition: func() *int { pos := 1; return &pos }(),
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/library", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify series was created
		var series models.Series
		err := db.Where("name = ?", seriesName).First(&series).Error
		assert.NoError(t, err)
	})

	t.Run("Add duplicate book", func(t *testing.T) {
		// First add
		reqBody := AddToLibraryRequest{
			Title:      "Duplicate Book",
			AuthorName: "Duplicate Author",
		}
		body, _ := json.Marshal(reqBody)

		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("POST", "/api/v1/library", bytes.NewBuffer(body))
		req1.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		// Try to add again
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/api/v1/library", bytes.NewBuffer(body))
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusConflict, w2.Code)
	})
}

func TestRemoveFromLibrary(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	item := models.LibraryItem{
		BookID:    book.ID,
		Status:    models.LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	db.Create(&item)

	router := gin.New()
	router.DELETE("/api/v1/library/:id", server.removeFromLibrary)

	t.Run("Remove existing library item", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/library/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify item was soft deleted
		var count int64
		db.Model(&models.LibraryItem{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Remove non-existent library item", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/library/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Remove library item with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/library/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestToLibraryItemResponse(t *testing.T) {
	item := &models.LibraryItem{
		ID:        1,
		BookID:    1,
		Status:    models.LibraryItemStatusWanted,
		AddedDate: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	item.Book = models.Book{
		ID:    1,
		Title: "Test Book",
		Author: models.Author{
			ID:   1,
			Name: "Test Author",
		},
	}

	response := toLibraryItemResponse(item)
	assert.NotNil(t, response)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "wanted", response.Status)
	assert.NotNil(t, response.Book)
	assert.Equal(t, "Test Book", response.Book.Title)
	assert.NotNil(t, response.Book.Author)
	assert.Equal(t, "Test Author", response.Book.Author.Name)
}
