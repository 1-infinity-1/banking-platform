package account

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

type accountClient interface {
	CreateAccount(ctx context.Context, params models.CreateAccountParams) (*models.Account, error)
	GetUserAccounts(ctx context.Context, params models.GetUserAccountsParams) ([]models.Account, error)
	GetAccount(ctx context.Context, params models.GetAccountParams) (*models.Account, error)
	GetBalance(ctx context.Context, params models.GetBalanceParams) (*models.Balance, error)
	UpdateStatus(ctx context.Context, params models.UpdateAccountStatusParams) (*models.Account, error)
}

type Service struct {
	accountClient accountClient
}

func New(accountClient accountClient) *Service {
	return &Service{accountClient: accountClient}
}

func (s *Service) CreateAccount(ctx context.Context, params models.CreateAccountParams) (*models.Account, error) {
	account, err := s.accountClient.CreateAccount(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.accountClient.CreateAccount: %w", err)
	}
	return account, nil
}

func (s *Service) GetUserAccounts(ctx context.Context, params models.GetUserAccountsParams) ([]models.Account, error) {
	accounts, err := s.accountClient.GetUserAccounts(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.accountClient.GetUserAccounts: %w", err)
	}
	return accounts, nil
}

func (s *Service) GetAccount(ctx context.Context, params models.GetAccountParams) (*models.Account, error) {
	account, err := s.accountClient.GetAccount(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.accountClient.GetAccount: %w", err)
	}
	return account, nil
}

func (s *Service) GetBalance(ctx context.Context, params models.GetBalanceParams) (*models.Balance, error) {
	balance, err := s.accountClient.GetBalance(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.accountClient.GetBalance: %w", err)
	}
	return balance, nil
}

func (s *Service) UpdateStatus(ctx context.Context, params models.UpdateAccountStatusParams) (*models.Account, error) {
	account, err := s.accountClient.UpdateStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.accountClient.UpdateStatus: %w", err)
	}
	return account, nil
}
