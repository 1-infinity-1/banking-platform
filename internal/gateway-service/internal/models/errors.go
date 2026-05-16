package models

import (
	"errors"
	"fmt"
)

var ErrInternal = errors.New("internal error")

type ConflictError struct {
	message string
	cause   error
}

func NewConflictError(message string, cause error) *ConflictError {
	return &ConflictError{message: message, cause: cause}
}

func (e *ConflictError) Error() string { return e.message }
func (e *ConflictError) Unwrap() error { return e.cause }

type ValidationError struct {
	field   string
	message string
	cause   error
}

func NewValidationError(field, message string, cause error) *ValidationError {
	return &ValidationError{field: field, message: message, cause: cause}
}

func (e *ValidationError) Error() string {
	if e.field == "" {
		return e.message
	}
	return fmt.Sprintf("field %s: %s", e.field, e.message)
}

func (e *ValidationError) Unwrap() error { return e.cause }

type NotFoundError struct {
	message string
	cause   error
}

func NewNotFoundError(message string, cause error) *NotFoundError {
	return &NotFoundError{message: message, cause: cause}
}

func (e *NotFoundError) Error() string { return e.message }
func (e *NotFoundError) Unwrap() error { return e.cause }
