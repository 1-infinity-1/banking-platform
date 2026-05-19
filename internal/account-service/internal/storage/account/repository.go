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
	"github.com/shopspring/decimal"
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

func (r *Repository) DebitTx(
	ctx context.Context,
	tx pgx.Tx,
	req models.DebitRequest,
) (*models.DebitResult, error) {
	var opID int64
	err := tx.QueryRow(ctx, `
		INSERT INTO account_operations (account_id, op_type, amount, idempotency_key)
		VALUES ($1, 'debit', $2, $3)
		RETURNING id
	`, req.AccountID, req.Amount, req.IdempotencyKey).Scan(&opID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			accountID, balanceAfter, lookupErr := r.getOperationByKeyTx(ctx, tx, req.IdempotencyKey)
			if lookupErr != nil {
				return nil, fmt.Errorf("r.getOperationByKeyTx: %w", lookupErr)
			}
			return &models.DebitResult{
				AccountID:    accountID,
				BalanceAfter: balanceAfter,
			}, nil
		}
		return nil, fmt.Errorf("tx.QueryRow(insert account_operation): %w", err)
	}

	var newBalance decimal.Decimal
	err = tx.QueryRow(ctx, `
		UPDATE accounts
		SET balance = balance - $1, updated_at = NOW()
		WHERE public_id = $2 AND balance >= $1 AND status = 'active'
		RETURNING balance
	`, req.Amount, req.AccountID).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, r.classifyDebitFailure(ctx, tx, req.AccountID, req.Amount)
		}
		return nil, fmt.Errorf("tx.QueryRow(update accounts debit): %w", err)
	}

	if _, err = tx.Exec(ctx, `
		UPDATE account_operations SET balance_after = $1 WHERE id = $2
	`, newBalance, opID); err != nil {
		return nil, fmt.Errorf("tx.Exec(update account_operation balance_after): %w", err)
	}

	return &models.DebitResult{
		AccountID:    req.AccountID.String(),
		BalanceAfter: newBalance,
	}, nil
}

func (r *Repository) CreditTx(
	ctx context.Context,
	tx pgx.Tx,
	req models.CreditRequest,
) (*models.CreditResult, error) {
	var opID int64
	err := tx.QueryRow(ctx, `
		INSERT INTO account_operations (account_id, op_type, amount, idempotency_key)
		VALUES ($1, 'credit', $2, $3)
		RETURNING id
	`, req.AccountID, req.Amount, req.IdempotencyKey).Scan(&opID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			accountID, balanceAfter, lookupErr := r.getOperationByKeyTx(ctx, tx, req.IdempotencyKey)
			if lookupErr != nil {
				return nil, fmt.Errorf("r.getOperationByKeyTx: %w", lookupErr)
			}
			return &models.CreditResult{
				AccountID:    accountID,
				BalanceAfter: balanceAfter,
			}, nil
		}
		return nil, fmt.Errorf("tx.QueryRow(insert account_operation): %w", err)
	}

	var newBalance decimal.Decimal
	err = tx.QueryRow(ctx, `
		UPDATE accounts
		SET balance = balance + $1, updated_at = NOW()
		WHERE public_id = $2 AND status = 'active'
		RETURNING balance
	`, req.Amount, req.AccountID).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, r.classifyCreditFailure(ctx, tx, req.AccountID)
		}
		return nil, fmt.Errorf("tx.QueryRow(update accounts credit): %w", err)
	}

	if _, err = tx.Exec(ctx, `
		UPDATE account_operations SET balance_after = $1 WHERE id = $2
	`, newBalance, opID); err != nil {
		return nil, fmt.Errorf("tx.Exec(update account_operation balance_after): %w", err)
	}

	return &models.CreditResult{
		AccountID:    req.AccountID.String(),
		BalanceAfter: newBalance,
	}, nil
}

func (r *Repository) getOperationByKeyTx(
	ctx context.Context,
	tx pgx.Tx,
	key string,
) (string, decimal.Decimal, error) {
	var accID uuid.UUID
	var balanceAfter decimal.Decimal
	err := tx.QueryRow(ctx, `
		SELECT account_id, balance_after FROM account_operations WHERE idempotency_key = $1
	`, key).Scan(&accID, &balanceAfter)
	if err != nil {
		return "", decimal.Zero, fmt.Errorf("tx.QueryRow(select account_operation): %w", err)
	}
	return accID.String(), balanceAfter, nil
}

func (r *Repository) classifyDebitFailure(
	ctx context.Context,
	tx pgx.Tx,
	accountID uuid.UUID,
	amount decimal.Decimal,
) error {
	var balance decimal.Decimal
	var status string
	err := tx.QueryRow(ctx, `SELECT balance, status FROM accounts WHERE public_id = $1`, accountID).
		Scan(&balance, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.NewNotFoundError("account not found")
		}
		return fmt.Errorf("tx.QueryRow(classify debit): %w", err)
	}
	if status != string(models.AccountStatusActive) {
		return models.NewBusinessError(fmt.Sprintf("account is not active: %s", status))
	}
	if balance.LessThan(amount) {
		return models.NewBusinessError("insufficient funds")
	}
	return models.NewBusinessError("debit precondition failed")
}

func (r *Repository) classifyCreditFailure(ctx context.Context, tx pgx.Tx, accountID uuid.UUID) error {
	var status string
	err := tx.QueryRow(ctx, `SELECT status FROM accounts WHERE public_id = $1`, accountID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.NewNotFoundError("account not found")
		}
		return fmt.Errorf("tx.QueryRow(classify credit): %w", err)
	}
	if status != string(models.AccountStatusActive) {
		return models.NewBusinessError(fmt.Sprintf("account is not active: %s", status))
	}
	return models.NewBusinessError("credit precondition failed")
}
