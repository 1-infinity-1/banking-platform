package models

import (
	"errors"
	"fmt"
)

var ErrInternal = errors.New("internal error")

type ForbiddenError struct {
	message string
	cause   error
}

func NewForbiddenError(message string, cause error) *ForbiddenError {
	return &ForbiddenError{message: message, cause: cause}
}

func (e *ForbiddenError) Error() string { return e.message }
func (e *ForbiddenError) Unwrap() error { return e.cause }

type UnauthorizedError struct {
	message string
	cause   error
}

func NewUnauthorizedError(message string, cause error) *UnauthorizedError {
	return &UnauthorizedError{message: message, cause: cause}
}

func (e *UnauthorizedError) Error() string { return e.message }
func (e *UnauthorizedError) Unwrap() error { return e.cause }

type BusinessError struct {
	message string
	cause   error
}

func NewBusinessError(message string, cause error) *BusinessError {
	return &BusinessError{message: message, cause: cause}
}

func (e *BusinessError) Error() string { return e.message }
func (e *BusinessError) Unwrap() error { return e.cause }

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
