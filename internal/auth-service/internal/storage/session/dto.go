package session

import (
	"time"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/google/uuid"
)

type SessionDTO struct {
	id         int64
	publicID   uuid.UUID
	userID     int64
	deviceID   int64
	status     string
	createdAt  time.Time
	updatedAt  time.Time
	expiresAt  time.Time
	lastSeenAt time.Time
}

func (s *SessionDTO) ToDomain() (*models.Session, error) {
	status, err := models.ToSessionStatus(s.status)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:         s.id,
		PublicID:   s.publicID,
		UserID:     s.userID,
		DeviceID:   s.deviceID,
		Status:     status,
		CreatedAt:  s.createdAt,
		UpdatedAt:  s.updatedAt,
		ExpiresAt:  s.expiresAt,
		LastSeenAt: s.lastSeenAt,
	}, nil
}
