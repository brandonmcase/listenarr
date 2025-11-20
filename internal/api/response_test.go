package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSuccessResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		SuccessResponse(c, StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.Empty(t, response.Error)
}

func TestErrorResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		err := errors.New("test error")
		ErrorResponse(c, StatusBadRequest, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "test error", response.Error)
}

func TestErrorResponseWithAPIError(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		err := ErrValidation("validation failed")
		err.WithDetail("field", "title")
		ErrorResponse(c, StatusUnprocessableEntity, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusUnprocessableEntity, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, ErrCodeValidation, response.Code)
	assert.NotNil(t, response.Details)
}

func TestValidationErrorResponse(t *testing.T) {
	t.Run("with ValidationErrors", func(t *testing.T) {
		router := setupTestRouter()
		router.POST("/test", func(c *gin.Context) {
			valErrs := NewValidationErrors()
			valErrs.Add("title", "required")
			valErrs.Add("author", "invalid")
			ValidationErrorResponse(c, valErrs)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, StatusUnprocessableEntity, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, ErrCodeValidation, response.Code)
		assert.NotNil(t, response.Details)
	})

	t.Run("with APIError", func(t *testing.T) {
		router := setupTestRouter()
		router.POST("/test", func(c *gin.Context) {
			err := ErrValidation("validation failed")
			ValidationErrorResponse(c, err)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, StatusUnprocessableEntity, w.Code)
	})

	t.Run("with standard error", func(t *testing.T) {
		router := setupTestRouter()
		router.POST("/test", func(c *gin.Context) {
			err := errors.New("standard error")
			ValidationErrorResponse(c, err)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, StatusUnprocessableEntity, w.Code)
	})
}

func TestNotFoundResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		NotFoundResponse(c, "book")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusNotFound, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, ErrCodeNotFound, response.Code)
	assert.Contains(t, response.Error, "book")
}

func TestConflictResponse(t *testing.T) {
	router := setupTestRouter()
	router.POST("/test", func(c *gin.Context) {
		ConflictResponse(c, "duplicate entry")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusConflict, w.Code)
}

func TestBadRequestResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		BadRequestResponse(c, "invalid request")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusBadRequest, w.Code)
}

func TestInternalErrorResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		InternalErrorResponse(c, "internal error")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusInternalServerError, w.Code)
}

func TestUnauthorizedResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		UnauthorizedResponse(c, "unauthorized")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusUnauthorized, w.Code)
}

func TestPaginatedSuccessResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		data := []string{"item1", "item2", "item3"}
		PaginatedSuccessResponse(c, data, 1, 20, 100)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusOK, w.Code)

	var response PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 20, response.Pagination.Limit)
	assert.Equal(t, 100, response.Pagination.Total)
	assert.Equal(t, 5, response.Pagination.TotalPages)
}

func TestPaginatedSuccessResponse_ZeroTotal(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		data := []string{}
		PaginatedSuccessResponse(c, data, 1, 20, 0)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Pagination.TotalPages) // Should be 1, not 0
}

func TestCreatedResponse(t *testing.T) {
	router := setupTestRouter()
	router.POST("/test", func(c *gin.Context) {
		CreatedResponse(c, gin.H{"id": 1})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusCreated, w.Code)
}

func TestNoContentResponse(t *testing.T) {
	router := setupTestRouter()
	router.DELETE("/test", func(c *gin.Context) {
		NoContentResponse(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}
