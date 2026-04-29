package models

import (
	"errors"
	"fmt"
)

var (
	ErrInternal = errors.New("internal error")
)

type NotFoundError struct {
	message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{message: message}
}

func (e *NotFoundError) Error() string {
	return e.message
}

type InvalidParamsError struct {
	field  string
	reason string
}

func NewInvalidParamsError(field, reason string) *InvalidParamsError {
	return &InvalidParamsError{field: field, reason: reason}
}

func (e *InvalidParamsError) Error() string {
	return fmt.Sprintf("Field `%s` invalid: %s", e.field, e.reason)
}

type BusinessError struct {
	message string
}

func NewBusinessError(message string) *BusinessError {
	return &BusinessError{message: message}
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("business error: %s", e.message)
}
