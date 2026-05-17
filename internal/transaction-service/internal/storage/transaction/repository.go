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

func (r *Repository) CreateTx(_ context.Context, _ pgx.Tx, _ models.TransferRequest) (*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("CreateTx: %w", models.ErrInternal)
}

func (r *Repository) CreateReplenishTx(
	_ context.Context,
	_ pgx.Tx,
	_ models.ReplenishRequest,
) (*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("CreateReplenishTx: %w", models.ErrInternal)
}

func (r *Repository) GetByIDTx(_ context.Context, _ pgx.Tx, _ uuid.UUID) (*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetByIDTx: %w", models.ErrInternal)
}

func (r *Repository) GetHistory(_ context.Context, _ models.GetHistoryRequest) ([]*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetHistory: %w", models.ErrInternal)
}

func (r *Repository) UpdateStatusTx(
	_ context.Context,
	_ pgx.Tx,
	_ uuid.UUID,
	_ models.TransactionStatus,
) (*models.Transaction, error) {
	// TODO: implement (idempotency_key unique constraint ensures at-most-once)
	return nil, fmt.Errorf("UpdateStatusTx: %w", models.ErrInternal)
}
