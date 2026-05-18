package ledger

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

type ledgerClient interface {
	GetStatement(ctx context.Context, params models.GetStatementParams) (*models.Statement, error)
}

type Service struct {
	ledgerClient ledgerClient
}

func New(ledgerClient ledgerClient) *Service {
	return &Service{ledgerClient: ledgerClient}
}

func (s *Service) GetStatement(ctx context.Context, params models.GetStatementParams) (*models.Statement, error) {
	statement, err := s.ledgerClient.GetStatement(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.ledgerClient.GetStatement: %w", err)
	}
	return statement, nil
}
