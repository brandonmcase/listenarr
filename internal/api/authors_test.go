package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/listenarr/listenarr/internal/models"
)

func TestGetAuthors(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	// Create test data
	author1 := models.Author{Name: "Author One"}
	author2 := models.Author{Name: "Author Two"}
	db.Create(&author1)
	db.Create(&author2)

	router := gin.New()
	router.GET("/api/v1/authors", server.getAuthors)

	t.Run("Get authors", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/authors", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
		assert.Equal(t, 2, response.Pagination.Total)
	})

	t.Run("Get authors with pagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/authors?page=1&limit=1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 1, response.Pagination.Limit)
		assert.Equal(t, 2, response.Pagination.Total)
	})

	t.Run("Get authors with search", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/authors?search=One", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})
}

func TestGetAuthor(t *testing.T) {
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

	router := gin.New()
	router.GET("/api/v1/authors/:id", server.getAuthor)

	t.Run("Get existing author", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/authors/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	t.Run("Get non-existent author", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/authors/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Get author with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/authors/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateAuthor(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	router := gin.New()
	router.POST("/api/v1/authors", server.createAuthor)

	t.Run("Create author", func(t *testing.T) {
		reqBody := CreateAuthorRequest{
			Name:      "New Author",
			Biography: "Author biography",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authors", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		// Verify author was created
		var count int64
		db.Model(&models.Author{}).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Create author with missing required fields", func(t *testing.T) {
		reqBody := map[string]interface{}{}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/authors", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Create duplicate author", func(t *testing.T) {
		// First create
		reqBody := CreateAuthorRequest{Name: "Duplicate Author"}
		body, _ := json.Marshal(reqBody)

		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("POST", "/api/v1/authors", bytes.NewBuffer(body))
		req1.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		// Try to create again
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/api/v1/authors", bytes.NewBuffer(body))
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusConflict, w2.Code)
	})
}

func TestUpdateAuthor(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	router := gin.New()
	router.PUT("/api/v1/authors/:id", server.updateAuthor)

	t.Run("Update author", func(t *testing.T) {
		// Create test author
		author := models.Author{Name: "Original Author"}
		db.Create(&author)

		newName := "Updated Author"
		reqBody := UpdateAuthorRequest{Name: &newName}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/authors/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify update
		var updatedAuthor models.Author
		db.First(&updatedAuthor, 1)
		assert.Equal(t, "Updated Author", updatedAuthor.Name)
	})

	t.Run("Update non-existent author", func(t *testing.T) {
		newName := "Updated Author"
		reqBody := UpdateAuthorRequest{Name: &newName}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/authors/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Update author with duplicate name", func(t *testing.T) {
		// Create two authors
		author1 := models.Author{Name: "First Author"}
		author2 := models.Author{Name: "Second Author"}
		db.Create(&author1)
		db.Create(&author2)

		// Reload to ensure IDs are set
		db.First(&author1, author1.ID)
		db.First(&author2, author2.ID)

		// Try to update author2 with author1's name
		duplicateName := author1.Name
		reqBody := UpdateAuthorRequest{Name: &duplicateName}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/authors/"+strconv.Itoa(int(author2.ID)), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestDeleteAuthor(t *testing.T) {
	db := setupTestDB(t)
	server := setupLibraryTestServer(db)

	router := gin.New()
	router.DELETE("/api/v1/authors/:id", server.deleteAuthor)

	t.Run("Delete author without books", func(t *testing.T) {
		author := models.Author{Name: "Author to Delete"}
		db.Create(&author)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/authors/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify soft delete
		var count int64
		db.Model(&models.Author{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Delete author with books", func(t *testing.T) {
		author := models.Author{Name: "Author with Books"}
		db.Create(&author)

		book := models.Book{
			Title:    "Book",
			AuthorID: author.ID,
		}
		db.Create(&book)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/authors/2", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Delete non-existent author", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/authors/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestToAuthorResponse(t *testing.T) {
	author := &models.Author{
		ID:          1,
		Name:        "Test Author",
		Biography:   "Biography",
		ImageURL:    "http://example.com/image.jpg",
		GoodreadsID: "12345",
	}

	response := toAuthorResponseDetailed(author)
	assert.NotNil(t, response)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Test Author", response.Name)
	assert.Equal(t, "Biography", response.Biography)
}
