package transaction

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Conn
}

func NewRepository(db *postgres.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTx(ctx context.Context, tx pgx.Tx, req models.TransferRequest) (*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("CreateTx: %w", models.ErrInternal)
}

func (r *Repository) CreateReplenishTx(ctx context.Context, tx pgx.Tx, req models.ReplenishRequest) (*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("CreateReplenishTx: %w", models.ErrInternal)
}

func (r *Repository) GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetByIDTx: %w", models.ErrInternal)
}

func (r *Repository) GetHistory(ctx context.Context, req models.GetHistoryRequest) ([]*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetHistory: %w", models.ErrInternal)
}

func (r *Repository) UpdateStatusTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, status models.TransactionStatus) (*models.Transaction, error) {
	// TODO: implement (idempotency_key unique constraint ensures at-most-once)
	return nil, fmt.Errorf("UpdateStatusTx: %w", models.ErrInternal)
}
