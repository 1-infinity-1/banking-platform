package models

import (
	"github.com/google/uuid"
)

type Role struct {
	ID       int64
	PublicID uuid.UUID
	Code     string
	Name     string
}
