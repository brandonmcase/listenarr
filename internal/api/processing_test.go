package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/listenarr/listenarr/internal/models"
)

func TestGetProcessingQueue(t *testing.T) {
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
		Status:        models.DownloadStatusCompleted,
	}
	db.Create(&download)

	task := models.ProcessingTask{
		DownloadID: download.ID,
		Status:     models.ProcessingStatusPending,
		InputPath:  "/tmp/download",
		Progress:   0,
	}
	db.Create(&task)

	router := gin.New()
	router.GET("/api/v1/processing", server.getProcessingQueue)

	t.Run("Get processing queue", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/processing", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 1, response.Pagination.Total)
	})
}

func TestGetProcessingTask(t *testing.T) {
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
		Status:        models.DownloadStatusCompleted,
	}
	db.Create(&download)

	task := models.ProcessingTask{
		DownloadID: download.ID,
		Status:     models.ProcessingStatusPending,
		InputPath:  "/tmp/download",
	}
	db.Create(&task)

	router := gin.New()
	router.GET("/api/v1/processing/:id", server.getProcessingTask)

	t.Run("Get existing processing task", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/processing/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get non-existent processing task", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/processing/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRetryProcessingTask(t *testing.T) {
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
		Status:        models.DownloadStatusCompleted,
	}
	db.Create(&download)

	task := models.ProcessingTask{
		DownloadID: download.ID,
		Status:     models.ProcessingStatusFailed,
		InputPath:  "/tmp/download",
		Error:      "Processing failed",
	}
	db.Create(&task)

	router := gin.New()
	router.POST("/api/v1/processing/:id/retry", server.retryProcessingTask)

	t.Run("Retry failed processing task", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/processing/1/retry", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify task was reset
		var updatedTask models.ProcessingTask
		db.First(&updatedTask, 1)
		assert.Equal(t, models.ProcessingStatusPending, updatedTask.Status)
		assert.Empty(t, updatedTask.Error)
	})

	t.Run("Retry non-failed task", func(t *testing.T) {
		// Create a pending task
		task2 := models.ProcessingTask{
			DownloadID: download.ID,
			Status:     models.ProcessingStatusPending,
			InputPath:  "/tmp/download2",
		}
		db.Create(&task2)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/processing/2/retry", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
