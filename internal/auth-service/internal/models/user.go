package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           int64
	PublicID     uuid.UUID
	Login        string
	Email        *string
	Phone        *string
	PasswordHash string
	Status       UserStatus
	Roles        []string
	Permissions  []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUser struct {
	Login string
	Email *string
	Phone *string
	Role  []string
}
