package models

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusUnspecified UserStatus = "unspecified"
	UserStatusActive      UserStatus = "active"
	UserStatusBlocked     UserStatus = "blocked"
	UserStatusLocked      UserStatus = "locked"
	UserStatusDisabled    UserStatus = "disabled"
)

type CreateUserParams struct {
	Login     string
	Email     *string
	Phone     *string
	Password  string
	RoleCodes []string
}

type User struct {
	ID        uuid.UUID
	Login     string
	Email     *string
	Phone     *string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
