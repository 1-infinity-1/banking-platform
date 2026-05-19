package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	accountColumns = `
		id,
		public_id,
		user_id,
		currency,
		balance,
		status,
		created_at,
		updated_at
	`

	pgUniqueViolation = "23505"
)

type Repository struct {
	db *postgres.Conn
}

func NewRepository(db *postgres.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAccountTx(
	ctx context.Context,
	tx pgx.Tx,
	req models.CreateAccountRequest,
) (*models.Account, error) {
	query := `
		INSERT INTO accounts (user_id, currency)
		VALUES ($1, $2)
		RETURNING ` + accountColumns

	var dto accountDTO
	err := tx.QueryRow(ctx, query, req.UserID, req.Currency).Scan(
		&dto.id,
		&dto.publicID,
		&dto.userID,
		&dto.currency,
		&dto.balance,
		&dto.status,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, models.NewConflictError("account already exists")
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain()
}

func (r *Repository) GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.Account, error) {
	query := `SELECT ` + accountColumns + ` FROM accounts WHERE public_id = $1`

	var dto accountDTO
	err := tx.QueryRow(ctx, query, id).Scan(
		&dto.id,
		&dto.publicID,
		&dto.userID,
		&dto.currency,
		&dto.balance,
		&dto.status,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewNotFoundError("account not found")
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain()
}

func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Account, error) {
	query := `SELECT ` + accountColumns + ` FROM accounts WHERE user_id = $1 ORDER BY created_at`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("r.db.Query: %w", err)
	}
	defer rows.Close()

	accounts := make([]*models.Account, 0)
	for rows.Next() {
		var dto accountDTO
		if err = rows.Scan(
			&dto.id,
			&dto.publicID,
			&dto.userID,
			&dto.currency,
			&dto.balance,
			&dto.status,
			&dto.createdAt,
			&dto.updatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		account, convErr := dto.ToDomain()
		if convErr != nil {
			return nil, fmt.Errorf("dto.ToDomain: %w", convErr)
		}
		accounts = append(accounts, account)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return accounts, nil
}

func (r *Repository) UpdateStatusTx(
	ctx context.Context,
	tx pgx.Tx,
	req models.UpdateStatusRequest,
) (*models.Account, error) {
	query := `
		UPDATE accounts
		SET status = $1, updated_at = NOW()
		WHERE public_id = $2
		RETURNING ` + accountColumns

	var dto accountDTO
	err := tx.QueryRow(ctx, query, string(req.Status), req.AccountID).Scan(
		&dto.id,
		&dto.publicID,
		&dto.userID,
		&dto.currency,
		&dto.balance,
		&dto.status,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewNotFoundError("account not found")
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain()
}

func (r *Repository) DebitTx(_ context.Context, _ pgx.Tx, _ models.DebitRequest) (*models.DebitResult, error) {
	// TODO: implement (use idempotency_key unique constraint)
	return nil, fmt.Errorf("DebitTx: %w", models.ErrInternal)
}

func (r *Repository) CreditTx(_ context.Context, _ pgx.Tx, _ models.CreditRequest) (*models.CreditResult, error) {
	// TODO: implement (use idempotency_key unique constraint)
	return nil, fmt.Errorf("CreditTx: %w", models.ErrInternal)
}
