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

func (s *Service) Transfer(_ context.Context, _ models.TransferRequest) (*models.Transaction, error) {
	// TODO: implement (Saga: debit from → credit to → save → publish event)
	return nil, fmt.Errorf("Transfer: %w", models.ErrInternal)
}

func (s *Service) Replenish(_ context.Context, _ models.ReplenishRequest) (*models.Transaction, error) {
	// TODO: implement (credit to → save → publish event)
	return nil, fmt.Errorf("Replenish: %w", models.ErrInternal)
}

func (s *Service) GetHistory(_ context.Context, _ models.GetHistoryRequest) ([]*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetHistory: %w", models.ErrInternal)
}

func (s *Service) GetTransaction(_ context.Context, _ uuid.UUID) (*models.Transaction, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetTransaction: %w", models.ErrInternal)
}
