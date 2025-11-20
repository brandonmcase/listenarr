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

	"github.com/listenarr/listenarr/internal/models"
)

func TestGetDownloads(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{Title: "Test Book", AuthorID: author.ID}
	db.Create(&book)

	libraryItem := models.LibraryItem{
		BookID:    book.ID,
		Status:    models.LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	db.Create(&libraryItem)

	release := models.Release{BookID: book.ID, Format: "m4b"}
	db.Create(&release)

	download := models.Download{
		LibraryItemID: libraryItem.ID,
		ReleaseID:     release.ID,
		Status:        models.DownloadStatusQueued,
		Progress:      0,
	}
	db.Create(&download)

	router := gin.New()
	router.GET("/api/v1/downloads", server.getDownloads)

	t.Run("Get downloads", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/downloads", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 1, response.Pagination.Total)
	})

	t.Run("Get downloads with status filter", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/downloads?status=queued", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetDownload(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{Title: "Test Book", AuthorID: author.ID}
	db.Create(&book)

	libraryItem := models.LibraryItem{
		BookID:    book.ID,
		Status:    models.LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	db.Create(&libraryItem)

	release := models.Release{BookID: book.ID}
	db.Create(&release)

	download := models.Download{
		LibraryItemID: libraryItem.ID,
		ReleaseID:     release.ID,
		Status:        models.DownloadStatusQueued,
	}
	db.Create(&download)

	router := gin.New()
	router.GET("/api/v1/downloads/:id", server.getDownload)

	t.Run("Get existing download", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/downloads/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get non-existent download", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/downloads/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestStartDownload(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{Title: "Test Book", AuthorID: author.ID}
	db.Create(&book)

	libraryItem := models.LibraryItem{
		BookID:    book.ID,
		Status:    models.LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	db.Create(&libraryItem)

	release := models.Release{BookID: book.ID, Format: "m4b"}
	db.Create(&release)

	router := gin.New()
	router.POST("/api/v1/downloads", server.startDownload)

	t.Run("Start download", func(t *testing.T) {
		reqBody := StartDownloadRequest{
			LibraryItemID: libraryItem.ID,
			ReleaseID:     release.ID,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/downloads", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify download was created
		var count int64
		db.Model(&models.Download{}).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Start download with invalid library item", func(t *testing.T) {
		reqBody := StartDownloadRequest{
			LibraryItemID: 999,
			ReleaseID:     release.ID,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/downloads", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCancelDownload(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{Title: "Test Book", AuthorID: author.ID}
	db.Create(&book)

	libraryItem := models.LibraryItem{
		BookID:    book.ID,
		Status:    models.LibraryItemStatusDownloading,
		AddedDate: time.Now(),
	}
	db.Create(&libraryItem)

	release := models.Release{BookID: book.ID}
	db.Create(&release)

	download := models.Download{
		LibraryItemID: libraryItem.ID,
		ReleaseID:     release.ID,
		Status:        models.DownloadStatusDownloading,
	}
	db.Create(&download)

	router := gin.New()
	router.DELETE("/api/v1/downloads/:id", server.cancelDownload)

	t.Run("Cancel active download", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/downloads/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify download status updated
		var updatedDownload models.Download
		db.First(&updatedDownload, 1)
		assert.Equal(t, models.DownloadStatusFailed, updatedDownload.Status)
	})
}
