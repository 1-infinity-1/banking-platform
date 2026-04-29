package device

import (
	"context"
	"fmt"

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

func (r *Repository) CreateDeviceTx(ctx context.Context, tx pgx.Tx, userAgent, platform string) (*models.Device, error) {
	query := `
		INSERT INTO devices (user_agent, platform)
		VALUES ($1, $2)
		RETURNING id, public_id, user_agent, platform, created_at, updated_at
	`

	var dto DeviceDTO
	err := tx.QueryRow(ctx, query, userAgent, platform).Scan(
		&dto.id,
		&dto.publicID,
		&dto.userAgent,
		&dto.platform,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain(), nil
}
