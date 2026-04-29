package models

import (
	"github.com/google/uuid"
)

type Device struct {
	ID        int64
	PublicID  uuid.UUID
	Platform  string
	UserAgent string
}
