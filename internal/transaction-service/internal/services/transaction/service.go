package transaction

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type txManager interface {
	BeginFunc(ctx context.Context, fn func(pgx.Tx) error) error
}

type transactionRepo interface {
	CreateTx(ctx context.Context, tx pgx.Tx, req models.TransferRequest) (*models.Transaction, error)
	CreateReplenishTx(ctx context.Context, tx pgx.Tx, req models.ReplenishRequest) (*models.Transaction, error)
	GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.Transaction, error)
	GetHistory(ctx context.Context, req models.GetHistoryRequest) ([]*models.Transaction, error)
	UpdateStatusTx(
		ctx context.Context,
		tx pgx.Tx,
		id uuid.UUID,
		status models.TransactionStatus,
	) (*models.Transaction, error)
}

// accountClient wraps the gRPC client to account-service behind a consumer-side interface.
type accountClient interface {
	Debit(ctx context.Context, req models.DebitRequest) (*models.DebitResult, error)
	Credit(ctx context.Context, req models.CreditRequest) (*models.CreditResult, error)
}

// eventPublisher abstracts the Kafka producer.
type eventPublisher interface {
	PublishTransactionCompleted(ctx context.Context, event models.TransactionEvent) error
}

type Service struct {
	txManager      txManager
	txRepo         transactionRepo
	accountClient  accountClient
	eventPublisher eventPublisher
}

func NewService(
	txManager txManager,
	txRepo transactionRepo,
	accountClient accountClient,
	eventPublisher eventPublisher,
) *Service {
	return &Service{
		txManager:      txManager,
		txRepo:         txRepo,
		accountClient:  accountClient,
		eventPublisher: eventPublisher,
	}
}

func (s *Service) Transfer(ctx context.Context, req models.TransferRequest) (*models.Transaction, error) {
	if err := validateTransfer(req); err != nil {
		return nil, err
	}

	pending, err := s.createPendingTransfer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("s.createPendingTransfer: %w", err)
	}

	switch pending.Status {
	case models.TransactionStatusCompleted:
		if err = s.eventPublisher.PublishTransactionCompleted(ctx, toEvent(pending)); err != nil {
			return nil, fmt.Errorf("s.eventPublisher.PublishTransactionCompleted: %w", err)
		}
		return pending, nil
	case models.TransactionStatusFailed:
		return nil, models.NewBusinessError("transaction already failed for given idempotency_key")
	case models.TransactionStatusPending,
		models.TransactionStatusCancelled,
		models.TransactionStatusUnspecified:
		// fall through to saga execution
	}

	_, err = s.accountClient.Debit(ctx, models.DebitRequest{
		AccountID:      req.FromAccountID,
		Amount:         req.Amount,
		IdempotencyKey: pending.PublicID.String() + ":debit",
	})
	if err != nil {
		s.markFailed(ctx, pending.PublicID)
		return nil, fmt.Errorf("s.accountClient.Debit: %w", err)
	}

	_, err = s.accountClient.Credit(ctx, models.CreditRequest{
		AccountID:      req.ToAccountID,
		Amount:         req.Amount,
		IdempotencyKey: pending.PublicID.String() + ":credit",
	})
	if err != nil {
		s.markFailed(ctx, pending.PublicID)
		return nil, fmt.Errorf("s.accountClient.Credit: %w", err)
	}

	completed, err := s.markCompleted(ctx, pending.PublicID)
	if err != nil {
		return nil, fmt.Errorf("s.markCompleted: %w", err)
	}

	if err = s.eventPublisher.PublishTransactionCompleted(ctx, toEvent(completed)); err != nil {
		return nil, fmt.Errorf("s.eventPublisher.PublishTransactionCompleted: %w", err)
	}

	return completed, nil
}

func (s *Service) Replenish(ctx context.Context, req models.ReplenishRequest) (*models.Transaction, error) {
	if err := validateReplenish(req); err != nil {
		return nil, err
	}

	pending, err := s.createPendingReplenish(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("s.createPendingReplenish: %w", err)
	}

	switch pending.Status {
	case models.TransactionStatusCompleted:
		if err = s.eventPublisher.PublishTransactionCompleted(ctx, toEvent(pending)); err != nil {
			return nil, fmt.Errorf("s.eventPublisher.PublishTransactionCompleted: %w", err)
		}
		return pending, nil
	case models.TransactionStatusFailed:
		return nil, models.NewBusinessError("transaction already failed for given idempotency_key")
	case models.TransactionStatusPending,
		models.TransactionStatusCancelled,
		models.TransactionStatusUnspecified:
		// fall through to saga execution
	}

	_, err = s.accountClient.Credit(ctx, models.CreditRequest{
		AccountID:      req.ToAccountID,
		Amount:         req.Amount,
		IdempotencyKey: pending.PublicID.String() + ":credit",
	})
	if err != nil {
		s.markFailed(ctx, pending.PublicID)
		return nil, fmt.Errorf("s.accountClient.Credit: %w", err)
	}

	completed, err := s.markCompleted(ctx, pending.PublicID)
	if err != nil {
		return nil, fmt.Errorf("s.markCompleted: %w", err)
	}

	if err = s.eventPublisher.PublishTransactionCompleted(ctx, toEvent(completed)); err != nil {
		return nil, fmt.Errorf("s.eventPublisher.PublishTransactionCompleted: %w", err)
	}

	return completed, nil
}

func (s *Service) GetHistory(ctx context.Context, req models.GetHistoryRequest) ([]*models.Transaction, error) {
	transactions, err := s.txRepo.GetHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("s.txRepo.GetHistory: %w", err)
	}
	return transactions, nil
}

func (s *Service) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction *models.Transaction
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		transaction, inner = s.txRepo.GetByIDTx(ctx, tx, id)
		if inner != nil {
			return fmt.Errorf("s.txRepo.GetByIDTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return transaction, nil
}

func validateReplenish(req models.ReplenishRequest) error {
	if req.ToAccountID == uuid.Nil {
		return models.NewInvalidParamsError("to_account_id", "must be non-zero UUID")
	}
	if !req.Amount.IsPositive() {
		return models.NewInvalidParamsError("amount", "must be positive")
	}
	if req.Currency == "" {
		return models.NewInvalidParamsError("currency", "must be non-empty")
	}
	if req.IdempotencyKey == "" {
		return models.NewInvalidParamsError("idempotency_key", "must be non-empty")
	}
	return nil
}

func (s *Service) createPendingReplenish(
	ctx context.Context,
	req models.ReplenishRequest,
) (*models.Transaction, error) {
	var pending *models.Transaction
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		pending, inner = s.txRepo.CreateReplenishTx(ctx, tx, req)
		if inner != nil {
			return fmt.Errorf("s.txRepo.CreateReplenishTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return pending, nil
}

func validateTransfer(req models.TransferRequest) error {
	if req.FromAccountID == req.ToAccountID {
		return models.NewInvalidParamsError("to_account_id", "must differ from from_account_id")
	}
	if !req.Amount.IsPositive() {
		return models.NewInvalidParamsError("amount", "must be positive")
	}
	if req.Currency == "" {
		return models.NewInvalidParamsError("currency", "must be non-empty")
	}
	if req.IdempotencyKey == "" {
		return models.NewInvalidParamsError("idempotency_key", "must be non-empty")
	}
	return nil
}

func (s *Service) createPendingTransfer(ctx context.Context, req models.TransferRequest) (*models.Transaction, error) {
	var pending *models.Transaction
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		pending, inner = s.txRepo.CreateTx(ctx, tx, req)
		if inner != nil {
			return fmt.Errorf("s.txRepo.CreateTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return pending, nil
}

func (s *Service) markFailed(ctx context.Context, id uuid.UUID) {
	_ = s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := s.txRepo.UpdateStatusTx(ctx, tx, id, models.TransactionStatusFailed)
		return err
	})
}

func (s *Service) markCompleted(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var completed *models.Transaction
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		completed, inner = s.txRepo.UpdateStatusTx(ctx, tx, id, models.TransactionStatusCompleted)
		if inner != nil {
			return fmt.Errorf("s.txRepo.UpdateStatusTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return completed, nil
}

func toEvent(tx *models.Transaction) models.TransactionEvent {
	var from *string
	if tx.FromAccountID != nil {
		s := tx.FromAccountID.String()
		from = &s
	}
	return models.TransactionEvent{
		TransactionID: tx.PublicID.String(),
		FromAccountID: from,
		ToAccountID:   tx.ToAccountID.String(),
		Amount:        tx.Amount.String(),
		Currency:      tx.Currency,
		Status:        tx.Status,
		OccurredAt:    tx.UpdatedAt,
	}
}
