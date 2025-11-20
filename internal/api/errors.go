package api

import (
	"errors"
	"fmt"
)

// Error codes for machine-readable error identification
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeBadRequest    = "BAD_REQUEST"
	ErrCodeUnprocessable = "UNPROCESSABLE_ENTITY"
)

// APIError represents an API error with code and message
type APIError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new APIError
func NewAPIError(code, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithDetail adds a detail to the error
func (e *APIError) WithDetail(key string, value interface{}) *APIError {
	e.Details[key] = value
	return e
}

// Predefined error constructors

// ErrValidation creates a validation error
func ErrValidation(message string) *APIError {
	return NewAPIError(ErrCodeValidation, message)
}

// ErrNotFound creates a not found error
func ErrNotFound(resource string) *APIError {
	return NewAPIError(ErrCodeNotFound, fmt.Sprintf("%s not found", resource))
}

// ErrConflict creates a conflict error
func ErrConflict(message string) *APIError {
	return NewAPIError(ErrCodeConflict, message)
}

// ErrUnauthorized creates an unauthorized error
func ErrUnauthorized(message string) *APIError {
	if message == "" {
		message = "Unauthorized"
	}
	return NewAPIError(ErrCodeUnauthorized, message)
}

// ErrInternal creates an internal server error
func ErrInternal(message string) *APIError {
	return NewAPIError(ErrCodeInternal, message)
}

// ErrBadRequest creates a bad request error
func ErrBadRequest(message string) *APIError {
	return NewAPIError(ErrCodeBadRequest, message)
}

// ErrUnprocessable creates an unprocessable entity error
func ErrUnprocessable(message string) *APIError {
	return NewAPIError(ErrCodeUnprocessable, message)
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError
}

// Error implements the error interface
func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %d error(s)", len(e.Errors))
}

// Add adds a validation error
func (e *ValidationErrors) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// NewValidationErrors creates a new ValidationErrors
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}

// IsAPIError checks if an error is an APIError
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// AsAPIError converts an error to APIError if possible
func AsAPIError(err error) (*APIError, bool) {
	apiErr, ok := err.(*APIError)
	return apiErr, ok
}

// IsValidationErrors checks if an error is ValidationErrors
func IsValidationErrors(err error) bool {
	_, ok := err.(*ValidationErrors)
	return ok
}

// AsValidationErrors converts an error to ValidationErrors if possible
func AsValidationErrors(err error) (*ValidationErrors, bool) {
	valErr, ok := err.(*ValidationErrors)
	return valErr, ok
}

// WrapError wraps a standard error as an APIError
func WrapError(err error, code, message string) *APIError {
	if err == nil {
		return nil
	}
	apiErr := NewAPIError(code, message)
	apiErr.WithDetail("original_error", err.Error())
	return apiErr
}

// Standard error variables for common cases
var (
	ErrInvalidID       = errors.New("invalid ID format")
	ErrMissingRequired = errors.New("missing required field")
	ErrInvalidFormat   = errors.New("invalid format")
	ErrDatabaseError   = errors.New("database error")
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateEntry  = errors.New("duplicate entry")
)
