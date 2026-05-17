package refreshtoken

import (
	"time"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
)

type refreshTokenDTO struct {
	id        int64
	sessionID int64
	tokenHash string
	expiresAt time.Time
	revokedAt *time.Time
	createdAt time.Time
}

func (d *refreshTokenDTO) ToDomain() *models.RefreshToken {
	return &models.RefreshToken{
		ID:        d.id,
		SessionID: d.sessionID,
		TokenHash: d.tokenHash,
		ExpiresAt: d.expiresAt,
		RevokedAt: d.revokedAt,
		CreatedAt: d.createdAt,
	}
}
