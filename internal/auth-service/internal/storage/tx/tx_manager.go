package tx

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
)

type Manager struct {
	db *postgres.Conn
}

func NewManager(db *postgres.Conn) *Manager {
	return &Manager{db: db}
}

func (t *Manager) BeginFunc(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("t.db.Begin: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = fn(tx); err != nil {
		return fmt.Errorf("fn: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
