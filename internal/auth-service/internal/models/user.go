package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64
	PublicID  uuid.UUID
	Login     string
	Email     *string
	Phone     *string
	Password  string
	Status    Status
	Role      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUser struct {
	Login    string
	Email    *string
	Phone    *string
	Password string
	Role     []string
}
