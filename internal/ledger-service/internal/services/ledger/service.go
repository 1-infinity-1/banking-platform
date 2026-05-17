package ledger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type entryRepo interface {
	CreateEntryTx(ctx context.Context, tx pgx.Tx, entry models.LedgerEntry) error
	GetStatementByAccountID(ctx context.Context, accountID uuid.UUID, from, to time.Time) ([]models.LedgerEntry, error)
}

type txManager interface {
	BeginFunc(ctx context.Context, fn func(pgx.Tx) error) error
}

type Service struct {
	entryRepo entryRepo
	txManager txManager
}

func NewService(txManager txManager, entryRepo entryRepo) *Service {
	return &Service{
		entryRepo: entryRepo,
		txManager: txManager,
	}
}

// RecordEntry persists a completed transaction event as an immutable ledger entry.
// Duplicate transaction_id (ConflictError) is silently ignored — idempotent at-least-once delivery.
func (s *Service) RecordEntry(ctx context.Context, event models.TransactionCompletedEvent) error {
	transactionID, err := uuid.Parse(event.TransactionID)
	if err != nil {
		return models.NewInvalidParamsError("transaction_id", fmt.Sprintf("invalid UUID: %s", err))
	}

	accountID, err := uuid.Parse(event.AccountID)
	if err != nil {
		return models.NewInvalidParamsError("account_id", fmt.Sprintf("invalid UUID: %s", err))
	}

	amount, err := decimal.NewFromString(event.Amount)
	if err != nil {
		return models.NewInvalidParamsError("amount", fmt.Sprintf("invalid decimal: %s", err))
	}

	balanceAfter, err := decimal.NewFromString(event.BalanceAfter)
	if err != nil {
		return models.NewInvalidParamsError("balance_after", fmt.Sprintf("invalid decimal: %s", err))
	}

	entry := models.LedgerEntry{
		TransactionID: transactionID,
		AccountID:     accountID,
		Type:          models.EntryType(event.Type),
		Amount:        amount,
		Currency:      event.Currency,
		BalanceAfter:  balanceAfter,
		Description:   event.Description,
		OccurredAt:    event.OccurredAt,
	}

	err = s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		return s.entryRepo.CreateEntryTx(ctx, tx, entry)
	})
	if err != nil {
		var conflictErr *models.ConflictError
		if errors.As(err, &conflictErr) {
			return nil
		}
		return fmt.Errorf("RecordEntry: %w", err)
	}

	return nil
}

// GetStatement returns the account statement for the requested period.
func (s *Service) GetStatement(ctx context.Context, accountID uuid.UUID, from, to time.Time) (models.Statement, error) {
	if !from.Before(to) {
		return models.Statement{}, models.NewInvalidParamsError("from/to", "from must be before to")
	}

	entries, err := s.entryRepo.GetStatementByAccountID(ctx, accountID, from, to)
	if err != nil {
		return models.Statement{}, fmt.Errorf("GetStatement: %w", err)
	}

	return models.Statement{
		AccountID: accountID,
		Entries:   entries,
	}, nil
}
