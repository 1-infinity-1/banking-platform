package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         int64
	PublicID   uuid.UUID
	UserID     int64
	DeviceID   int64
	Status     SessionStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ExpiresAt  time.Time
	LastSeenAt time.Time
}
