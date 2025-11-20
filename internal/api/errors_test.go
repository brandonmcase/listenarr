package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError(t *testing.T) {
	t.Run("NewAPIError", func(t *testing.T) {
		err := NewAPIError(ErrCodeValidation, "test error")
		assert.NotNil(t, err)
		assert.Equal(t, ErrCodeValidation, err.Code)
		assert.Equal(t, "test error", err.Message)
		assert.NotNil(t, err.Details)
	})

	t.Run("Error", func(t *testing.T) {
		err := NewAPIError(ErrCodeValidation, "test error")
		assert.Equal(t, "test error", err.Error())
	})

	t.Run("WithDetail", func(t *testing.T) {
		err := NewAPIError(ErrCodeValidation, "test error")
		err.WithDetail("field", "value")
		assert.Equal(t, "value", err.Details["field"])
	})
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name        string
		constructor func(string) *APIError
		code        string
		message     string
	}{
		{"ErrValidation", ErrValidation, ErrCodeValidation, "validation error"},
		{"ErrNotFound", func(msg string) *APIError { return ErrNotFound("resource") }, ErrCodeNotFound, "resource not found"},
		{"ErrConflict", ErrConflict, ErrCodeConflict, "conflict error"},
		{"ErrUnauthorized", ErrUnauthorized, ErrCodeUnauthorized, "unauthorized"},
		{"ErrInternal", ErrInternal, ErrCodeInternal, "internal error"},
		{"ErrBadRequest", ErrBadRequest, ErrCodeBadRequest, "bad request"},
		{"ErrUnprocessable", ErrUnprocessable, ErrCodeUnprocessable, "unprocessable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor(tt.message)
			assert.Equal(t, tt.code, err.Code)
			assert.Equal(t, tt.message, err.Message)
		})
	}

	t.Run("ErrNotFound with resource", func(t *testing.T) {
		err := ErrNotFound("book")
		assert.Equal(t, "book not found", err.Message)
	})

	t.Run("ErrUnauthorized with empty message", func(t *testing.T) {
		err := ErrUnauthorized("")
		assert.Equal(t, "Unauthorized", err.Message)
	})
}

func TestValidationError(t *testing.T) {
	t.Run("ValidationError", func(t *testing.T) {
		err := ValidationError{
			Field:   "title",
			Message: "required",
		}
		assert.Equal(t, "title: required", err.Error())
	})
}

func TestValidationErrors(t *testing.T) {
	t.Run("NewValidationErrors", func(t *testing.T) {
		errs := NewValidationErrors()
		assert.NotNil(t, errs)
		assert.False(t, errs.HasErrors())
	})

	t.Run("Add", func(t *testing.T) {
		errs := NewValidationErrors()
		errs.Add("title", "required")
		errs.Add("author", "invalid")

		assert.True(t, errs.HasErrors())
		assert.Len(t, errs.Errors, 2)
		assert.Equal(t, "title", errs.Errors[0].Field)
		assert.Equal(t, "required", errs.Errors[0].Message)
	})

	t.Run("Error", func(t *testing.T) {
		errs := NewValidationErrors()
		assert.Equal(t, "validation failed", errs.Error())

		errs.Add("field", "error")
		assert.Contains(t, errs.Error(), "1 error(s)")
	})
}

func TestErrorTypeChecks(t *testing.T) {
	t.Run("IsAPIError", func(t *testing.T) {
		apiErr := ErrValidation("test")
		assert.True(t, IsAPIError(apiErr))

		stdErr := errors.New("standard error")
		assert.False(t, IsAPIError(stdErr))
	})

	t.Run("AsAPIError", func(t *testing.T) {
		apiErr := ErrValidation("test")
		converted, ok := AsAPIError(apiErr)
		assert.True(t, ok)
		assert.Equal(t, apiErr, converted)

		stdErr := errors.New("standard error")
		_, ok = AsAPIError(stdErr)
		assert.False(t, ok)
	})

	t.Run("IsValidationErrors", func(t *testing.T) {
		valErrs := NewValidationErrors()
		valErrs.Add("field", "error")
		assert.True(t, IsValidationErrors(valErrs))

		stdErr := errors.New("standard error")
		assert.False(t, IsValidationErrors(stdErr))
	})

	t.Run("AsValidationErrors", func(t *testing.T) {
		valErrs := NewValidationErrors()
		valErrs.Add("field", "error")
		converted, ok := AsValidationErrors(valErrs)
		assert.True(t, ok)
		assert.Equal(t, valErrs, converted)

		stdErr := errors.New("standard error")
		_, ok = AsValidationErrors(stdErr)
		assert.False(t, ok)
	})
}

func TestWrapError(t *testing.T) {
	t.Run("WrapError", func(t *testing.T) {
		originalErr := errors.New("original error")
		wrapped := WrapError(originalErr, ErrCodeInternal, "wrapped error")

		assert.NotNil(t, wrapped)
		assert.Equal(t, ErrCodeInternal, wrapped.Code)
		assert.Equal(t, "wrapped error", wrapped.Message)
		assert.Equal(t, "original error", wrapped.Details["original_error"])
	})

	t.Run("WrapError with nil", func(t *testing.T) {
		wrapped := WrapError(nil, ErrCodeInternal, "wrapped error")
		assert.Nil(t, wrapped)
	})
}
