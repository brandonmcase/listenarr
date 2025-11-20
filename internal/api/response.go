package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Success bool                   `json:"success"`
	Data    interface{}            `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool                   `json:"success"`
	Data       interface{}            `json:"data,omitempty"`
	Pagination PaginationInfo         `json:"pagination,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Code       string                 `json:"code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// HTTP status code constants
const (
	StatusOK                  = http.StatusOK                  // 200
	StatusCreated             = http.StatusCreated             // 201
	StatusNoContent           = http.StatusNoContent           // 204
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusNotFound            = http.StatusNotFound            // 404
	StatusConflict            = http.StatusConflict            // 409
	StatusUnprocessableEntity = http.StatusUnprocessableEntity // 422
	StatusInternalServerError = http.StatusInternalServerError // 500
)

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Error:   "",
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, err error) {
	response := Response{
		Success: false,
		Data:    nil,
		Error:   err.Error(),
	}

	// If error is an APIError, include code and details
	if apiErr, ok := AsAPIError(err); ok {
		response.Code = apiErr.Code
		if len(apiErr.Details) > 0 {
			response.Details = apiErr.Details
		}
	}

	c.JSON(statusCode, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	response := Response{
		Success: false,
		Data:    nil,
		Code:    ErrCodeValidation,
	}

	// Handle ValidationErrors type
	if valErrs, ok := AsValidationErrors(err); ok {
		response.Error = "Validation failed"
		response.Details = map[string]interface{}{
			"errors": valErrs.Errors,
		}
		c.JSON(StatusUnprocessableEntity, response)
		return
	}

	// Handle APIError
	if apiErr, ok := AsAPIError(err); ok {
		response.Error = apiErr.Message
		response.Code = apiErr.Code
		if len(apiErr.Details) > 0 {
			response.Details = apiErr.Details
		}
		c.JSON(StatusUnprocessableEntity, response)
		return
	}

	// Fallback to generic error
	response.Error = err.Error()
	c.JSON(StatusUnprocessableEntity, response)
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, resource string) {
	err := ErrNotFound(resource)
	ErrorResponse(c, StatusNotFound, err)
}

// ConflictResponse sends a conflict response
func ConflictResponse(c *gin.Context, message string) {
	err := ErrConflict(message)
	ErrorResponse(c, StatusConflict, err)
}

// BadRequestResponse sends a bad request response
func BadRequestResponse(c *gin.Context, message string) {
	err := ErrBadRequest(message)
	ErrorResponse(c, StatusBadRequest, err)
}

// InternalErrorResponse sends an internal server error response
func InternalErrorResponse(c *gin.Context, message string) {
	err := ErrInternal(message)
	ErrorResponse(c, StatusInternalServerError, err)
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	err := ErrUnauthorized(message)
	ErrorResponse(c, StatusUnauthorized, err)
}

// PaginatedSuccessResponse sends a successful paginated response
func PaginatedSuccessResponse(c *gin.Context, data interface{}, page, limit, total int) {
	totalPages := (total + limit - 1) / limit // Ceiling division
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Error: "",
	})
}

// CreatedResponse sends a created response (201)
func CreatedResponse(c *gin.Context, data interface{}) {
	SuccessResponse(c, StatusCreated, data)
}

// NoContentResponse sends a no content response (204)
func NoContentResponse(c *gin.Context) {
	c.Status(StatusNoContent)
}
