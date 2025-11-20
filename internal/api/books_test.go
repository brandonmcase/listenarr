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

func TestGetBooks(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book1 := models.Book{Title: "Book One", AuthorID: author.ID}
	book2 := models.Book{Title: "Book Two", AuthorID: author.ID}
	db.Create(&book1)
	db.Create(&book2)

	router := gin.New()
	router.GET("/api/v1/books", server.getBooks)

	t.Run("Get books", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/books", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 2, response.Pagination.Total)
	})

	t.Run("Get books with author filter", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/books?author_id=1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get books with search", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/books?search=One", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetBook(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{Title: "Test Book", AuthorID: author.ID}
	db.Create(&book)

	router := gin.New()
	router.GET("/api/v1/books/:id", server.getBook)

	t.Run("Get existing book", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/books/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Get non-existent book", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/books/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCreateBook(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create author
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	router := gin.New()
	router.POST("/api/v1/books", server.createBook)

	t.Run("Create book", func(t *testing.T) {
		isbn := "1234567890"
		reqBody := CreateBookRequest{
			Title:    "New Book",
			AuthorID: author.ID,
			ISBN:     &isbn,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify book was created
		var count int64
		db.Model(&models.Book{}).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Create book with invalid author", func(t *testing.T) {
		reqBody := CreateBookRequest{
			Title:    "New Book",
			AuthorID: 999,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Create duplicate book", func(t *testing.T) {
		reqBody := CreateBookRequest{
			Title:    "Duplicate Book",
			AuthorID: author.ID,
		}
		body, _ := json.Marshal(reqBody)

		// First create
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(body))
		req1.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		// Try to create again
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(body))
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusConflict, w2.Code)
	})
}

func TestUpdateBook(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{Title: "Original Book", AuthorID: author.ID}
	db.Create(&book)

	router := gin.New()
	router.PUT("/api/v1/books/:id", server.updateBook)

	t.Run("Update book", func(t *testing.T) {
		newTitle := "Updated Book"
		reqBody := UpdateBookRequest{Title: &newTitle}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/books/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify update
		var updatedBook models.Book
		db.First(&updatedBook, 1)
		assert.Equal(t, "Updated Book", updatedBook.Title)
	})

	t.Run("Update non-existent book", func(t *testing.T) {
		newTitle := "Updated Book"
		reqBody := UpdateBookRequest{Title: &newTitle}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/books/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteBook(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author := models.Author{Name: "Test Author"}
	db.Create(&author)

	book := models.Book{Title: "Book to Delete", AuthorID: author.ID}
	db.Create(&book)

	router := gin.New()
	router.DELETE("/api/v1/books/:id", server.deleteBook)

	t.Run("Delete book without library items", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/books/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify soft delete
		var count int64
		db.Model(&models.Book{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Delete book with library items", func(t *testing.T) {
		// Create book with library item
		book2 := models.Book{Title: "Book with Library Item", AuthorID: author.ID}
		db.Create(&book2)

		libraryItem := models.LibraryItem{
			BookID:    book2.ID,
			Status:    models.LibraryItemStatusWanted,
			AddedDate: time.Now(),
		}
		db.Create(&libraryItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/books/2", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
