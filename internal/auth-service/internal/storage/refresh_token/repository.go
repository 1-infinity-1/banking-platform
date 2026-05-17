package refreshtoken

import (
	"context"
	"errors"
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

func (r *Repository) CreateTokenTx(
	ctx context.Context,
	tx pgx.Tx,
	sessionID int64,
	tokenHash string,
	expireTime time.Time,
) error {
	query := `
		INSERT INTO refresh_tokens (session_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := tx.Exec(ctx, query, sessionID, tokenHash, expireTime)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}

func (r *Repository) GetByTokenHashTx(
	ctx context.Context,
	tx pgx.Tx,
	tokenHash string,
) (*models.RefreshToken, error) {
	query := `
		SELECT id, session_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var dto refreshTokenDTO
	err := tx.QueryRow(ctx, query, tokenHash).Scan(
		&dto.id,
		&dto.sessionID,
		&dto.tokenHash,
		&dto.expiresAt,
		&dto.revokedAt,
		&dto.createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewNotFoundError("refresh token not found")
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain(), nil
}

func (r *Repository) RevokeTokenTx(ctx context.Context, tx pgx.Tx, id int64) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`

	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}
