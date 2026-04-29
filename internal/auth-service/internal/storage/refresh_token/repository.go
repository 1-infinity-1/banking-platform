package refreshtoken

import (
	"context"
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Conn
}

func NewRepository(db *postgres.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTokenTx(ctx context.Context, tx pgx.Tx, sessionID int64, tokenHash string, expireTime time.Time) error {
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
