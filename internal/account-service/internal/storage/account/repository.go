package account

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
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

func (r *Repository) CreateAccountTx(
	_ context.Context,
	_ pgx.Tx,
	_ models.CreateAccountRequest,
) (*models.Account, error) {
	// TODO: implement
	return nil, fmt.Errorf("CreateAccountTx: %w", models.ErrInternal)
}

func (r *Repository) GetByIDTx(_ context.Context, _ pgx.Tx, _ uuid.UUID) (*models.Account, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetByIDTx: %w", models.ErrInternal)
}

func (r *Repository) GetByUserID(_ context.Context, _ uuid.UUID) ([]*models.Account, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetByUserID: %w", models.ErrInternal)
}

func (r *Repository) UpdateStatusTx(
	_ context.Context,
	_ pgx.Tx,
	_ models.UpdateStatusRequest,
) (*models.Account, error) {
	// TODO: implement
	return nil, fmt.Errorf("UpdateStatusTx: %w", models.ErrInternal)
}

func (r *Repository) DebitTx(_ context.Context, _ pgx.Tx, _ models.DebitRequest) (*models.DebitResult, error) {
	// TODO: implement (use idempotency_key unique constraint)
	return nil, fmt.Errorf("DebitTx: %w", models.ErrInternal)
}

func (r *Repository) CreditTx(_ context.Context, _ pgx.Tx, _ models.CreditRequest) (*models.CreditResult, error) {
	// TODO: implement (use idempotency_key unique constraint)
	return nil, fmt.Errorf("CreditTx: %w", models.ErrInternal)
}
