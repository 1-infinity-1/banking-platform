package entry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/models"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *postgres.Conn
}

func NewRepository(db *postgres.Conn) *Repository {
	return &Repository{db: db}
}

// CreateEntryTx inserts an immutable ledger entry inside the given transaction.
// Returns ConflictError when transaction_id already exists (idempotent re-delivery).
func (r *Repository) CreateEntryTx(ctx context.Context, tx pgx.Tx, entry models.LedgerEntry) error {
	query := `
		INSERT INTO ledger_entries (
			transaction_id,
			account_id,
			type,
			amount,
			currency,
			balance_after,
			description,
			occurred_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	var desc *string
	if entry.Description != "" {
		desc = &entry.Description
	}

	_, err := tx.Exec(ctx, query,
		entry.TransactionID,
		entry.AccountID,
		string(entry.Type),
		entry.Amount.String(),
		entry.Currency,
		entry.BalanceAfter.String(),
		desc,
		entry.OccurredAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.NewConflictError("ledger entry already exists for transaction_id")
		}
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}

// GetStatementByAccountID returns all ledger entries for an account in the given period,
// ordered by occurred_at ascending. Returns an empty slice when no entries match.
func (r *Repository) GetStatementByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
	from, to time.Time,
) ([]models.LedgerEntry, error) {
	query := `
		SELECT
			id,
			public_id,
			transaction_id,
			account_id,
			type,
			amount::text,
			currency,
			balance_after::text,
			description,
			occurred_at,
			created_at
		FROM ledger_entries
		WHERE account_id = $1
		  AND occurred_at >= $2
		  AND occurred_at <= $3
		ORDER BY occurred_at ASC
	`

	rows, err := r.db.Query(ctx, query, accountID, from, to)
	if err != nil {
		return nil, fmt.Errorf("r.db.Query: %w", err)
	}
	defer rows.Close()

	var entries []models.LedgerEntry
	for rows.Next() {
		var dto entryDTO
		if err = rows.Scan(
			&dto.id,
			&dto.publicID,
			&dto.transactionID,
			&dto.accountID,
			&dto.entryType,
			&dto.amount,
			&dto.currency,
			&dto.balanceAfter,
			&dto.description,
			&dto.occurredAt,
			&dto.createdAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		domainEntry, domainErr := dto.ToDomain()
		if domainErr != nil {
			return nil, fmt.Errorf("dto.ToDomain: %w", domainErr)
		}

		entries = append(entries, domainEntry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return entries, nil
}
