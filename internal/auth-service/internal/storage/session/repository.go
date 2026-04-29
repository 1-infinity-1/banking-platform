package session

import (
	"context"
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Conn
}

func NewRepository(db *postgres.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSessionTx(ctx context.Context, tx pgx.Tx, userID, deviceID int64, status models.SessionStatus, expireTime time.Time) (*models.Session, error) {
	query := `
		INSERT INTO sessions (user_id, device_id, status, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, public_id, user_id, device_id, status, created_at, updated_at, expires_at, COALESCE(last_seen_at, created_at)
	`

	var dto SessionDTO
	err := tx.QueryRow(ctx, query, userID, deviceID, status, expireTime).Scan(
		&dto.id,
		&dto.publicID,
		&dto.userID,
		&dto.deviceID,
		&dto.status,
		&dto.createdAt,
		&dto.updatedAt,
		&dto.expiresAt,
		&dto.lastSeenAt,
	)
	if err != nil {
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	session, err := dto.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("dto.ToDomain: %w", err)
	}

	return session, nil
}
