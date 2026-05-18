package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	transactionColumns = `
		id,
		public_id,
		from_account_id,
		to_account_id,
		amount,
		currency,
		status,
		idempotency_key,
		created_at,
		updated_at
	`

	defaultHistoryLimit = 50
	maxHistoryLimit     = 200

	pgUniqueViolation = "23505"
)

type Repository struct {
	db *postgres.Conn
}

func NewRepository(db *postgres.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTx(ctx context.Context, tx pgx.Tx, req models.TransferRequest) (*models.Transaction, error) {
	query := `
		INSERT INTO transactions
			(from_account_id, to_account_id, amount, currency, status, idempotency_key)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING ` + transactionColumns

	var dto transactionDTO
	err := tx.QueryRow(
		ctx,
		query,
		req.FromAccountID,
		req.ToAccountID,
		req.Amount,
		req.Currency,
		string(models.TransactionStatusPending),
		req.IdempotencyKey,
	).Scan(
		&dto.id,
		&dto.publicID,
		&dto.fromAccountID,
		&dto.toAccountID,
		&dto.amount,
		&dto.currency,
		&dto.status,
		&dto.idempotencyKey,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			existing, lookupErr := r.getByIdempotencyKeyTx(ctx, tx, req.IdempotencyKey)
			if lookupErr != nil {
				return nil, fmt.Errorf("r.getByIdempotencyKeyTx: %w", lookupErr)
			}
			return existing, nil
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain()
}

func (r *Repository) CreateReplenishTx(
	ctx context.Context,
	tx pgx.Tx,
	req models.ReplenishRequest,
) (*models.Transaction, error) {
	query := `
		INSERT INTO transactions
			(from_account_id, to_account_id, amount, currency, status, idempotency_key)
		VALUES (NULL, $1, $2, $3, $4, $5)
		RETURNING ` + transactionColumns

	var dto transactionDTO
	err := tx.QueryRow(
		ctx,
		query,
		req.ToAccountID,
		req.Amount,
		req.Currency,
		string(models.TransactionStatusPending),
		req.IdempotencyKey,
	).Scan(
		&dto.id,
		&dto.publicID,
		&dto.fromAccountID,
		&dto.toAccountID,
		&dto.amount,
		&dto.currency,
		&dto.status,
		&dto.idempotencyKey,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			existing, lookupErr := r.getByIdempotencyKeyTx(ctx, tx, req.IdempotencyKey)
			if lookupErr != nil {
				return nil, fmt.Errorf("r.getByIdempotencyKeyTx: %w", lookupErr)
			}
			return existing, nil
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain()
}

func (r *Repository) getByIdempotencyKeyTx(ctx context.Context, tx pgx.Tx, key string) (*models.Transaction, error) {
	query := `SELECT ` + transactionColumns + ` FROM transactions WHERE idempotency_key = $1`

	var dto transactionDTO
	err := tx.QueryRow(ctx, query, key).Scan(
		&dto.id,
		&dto.publicID,
		&dto.fromAccountID,
		&dto.toAccountID,
		&dto.amount,
		&dto.currency,
		&dto.status,
		&dto.idempotencyKey,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewNotFoundError("transaction not found by idempotency_key")
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain()
}

func (r *Repository) GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.Transaction, error) {
	query := `SELECT ` + transactionColumns + ` FROM transactions WHERE public_id = $1`

	var dto transactionDTO
	err := tx.QueryRow(ctx, query, id).Scan(
		&dto.id,
		&dto.publicID,
		&dto.fromAccountID,
		&dto.toAccountID,
		&dto.amount,
		&dto.currency,
		&dto.status,
		&dto.idempotencyKey,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewNotFoundError("transaction not found")
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain()
}

func (r *Repository) GetHistory(ctx context.Context, req models.GetHistoryRequest) ([]*models.Transaction, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = defaultHistoryLimit
	}
	if limit > maxHistoryLimit {
		limit = maxHistoryLimit
	}
	offset := max(req.Offset, 0)

	query := `
		SELECT ` + transactionColumns + `
		FROM transactions
		WHERE to_account_id = $1 OR from_account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, req.AccountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("r.db.Query: %w", err)
	}
	defer rows.Close()

	transactions := make([]*models.Transaction, 0, limit)
	for rows.Next() {
		var dto transactionDTO
		if err = rows.Scan(
			&dto.id,
			&dto.publicID,
			&dto.fromAccountID,
			&dto.toAccountID,
			&dto.amount,
			&dto.currency,
			&dto.status,
			&dto.idempotencyKey,
			&dto.createdAt,
			&dto.updatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		var domain *models.Transaction
		domain, err = dto.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("dto.ToDomain: %w", err)
		}

		transactions = append(transactions, domain)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return transactions, nil
}

func (r *Repository) UpdateStatusTx(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
	status models.TransactionStatus,
) (*models.Transaction, error) {
	query := `
		UPDATE transactions
		SET status = $1, updated_at = NOW()
		WHERE public_id = $2
		RETURNING ` + transactionColumns

	var dto transactionDTO
	err := tx.QueryRow(ctx, query, string(status), id).Scan(
		&dto.id,
		&dto.publicID,
		&dto.fromAccountID,
		&dto.toAccountID,
		&dto.amount,
		&dto.currency,
		&dto.status,
		&dto.idempotencyKey,
		&dto.createdAt,
		&dto.updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewNotFoundError("transaction not found")
		}
		return nil, fmt.Errorf("tx.QueryRow: %w", err)
	}

	return dto.ToDomain()
}
