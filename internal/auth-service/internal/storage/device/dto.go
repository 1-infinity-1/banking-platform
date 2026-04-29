package device

import (
	"time"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/google/uuid"
)

type DeviceDTO struct {
	id        int64
	publicID  uuid.UUID
	userAgent string
	platform  string
	createdAt time.Time
	updatedAt time.Time
}

func (d *DeviceDTO) ToDomain() *models.Device {
	return &models.Device{
		ID:        d.id,
		PublicID:  d.publicID,
		UserAgent: d.userAgent,
		Platform:  d.platform,
	}
}
