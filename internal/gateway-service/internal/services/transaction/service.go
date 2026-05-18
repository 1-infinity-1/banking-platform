package transaction

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

type transactionClient interface {
	Transfer(ctx context.Context, params models.TransferParams) (*models.Transaction, error)
	Replenish(ctx context.Context, params models.ReplenishParams) (*models.Transaction, error)
	GetHistory(ctx context.Context, params models.GetHistoryParams) ([]models.Transaction, error)
	GetTransaction(ctx context.Context, params models.GetTransactionParams) (*models.Transaction, error)
}

type Service struct {
	transactionClient transactionClient
}

func New(transactionClient transactionClient) *Service {
	return &Service{transactionClient: transactionClient}
}

func (s *Service) Transfer(ctx context.Context, params models.TransferParams) (*models.Transaction, error) {
	tx, err := s.transactionClient.Transfer(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.transactionClient.Transfer: %w", err)
	}
	return tx, nil
}

func (s *Service) Replenish(ctx context.Context, params models.ReplenishParams) (*models.Transaction, error) {
	tx, err := s.transactionClient.Replenish(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.transactionClient.Replenish: %w", err)
	}
	return tx, nil
}

func (s *Service) GetHistory(ctx context.Context, params models.GetHistoryParams) ([]models.Transaction, error) {
	txs, err := s.transactionClient.GetHistory(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.transactionClient.GetHistory: %w", err)
	}
	return txs, nil
}

func (s *Service) GetTransaction(ctx context.Context, params models.GetTransactionParams) (*models.Transaction, error) {
	tx, err := s.transactionClient.GetTransaction(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.transactionClient.GetTransaction: %w", err)
	}
	return tx, nil
}
