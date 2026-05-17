package account

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type txManager interface {
	BeginFunc(ctx context.Context, fn func(pgx.Tx) error) error
}

type accountRepo interface {
	CreateAccountTx(ctx context.Context, tx pgx.Tx, req models.CreateAccountRequest) (*models.Account, error)
	GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Account, error)
	UpdateStatusTx(ctx context.Context, tx pgx.Tx, req models.UpdateStatusRequest) (*models.Account, error)
	DebitTx(ctx context.Context, tx pgx.Tx, req models.DebitRequest) (*models.DebitResult, error)
	CreditTx(ctx context.Context, tx pgx.Tx, req models.CreditRequest) (*models.CreditResult, error)
}

type Service struct {
	txManager   txManager
	accountRepo accountRepo
}

func NewService(txManager txManager, accountRepo accountRepo) *Service {
	return &Service{
		txManager:   txManager,
		accountRepo: accountRepo,
	}
}

func (s *Service) CreateAccount(_ context.Context, _ models.CreateAccountRequest) (*models.Account, error) {
	// TODO: implement
	return nil, fmt.Errorf("CreateAccount: %w", models.ErrInternal)
}

func (s *Service) GetUserAccounts(_ context.Context, _ uuid.UUID) ([]*models.Account, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetUserAccounts: %w", models.ErrInternal)
}

func (s *Service) GetAccount(_ context.Context, _ uuid.UUID) (*models.Account, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetAccount: %w", models.ErrInternal)
}

func (s *Service) GetBalance(_ context.Context, _ uuid.UUID) (*models.Balance, error) {
	// TODO: implement
	return nil, fmt.Errorf("GetBalance: %w", models.ErrInternal)
}

func (s *Service) UpdateStatus(_ context.Context, _ models.UpdateStatusRequest) (*models.Account, error) {
	// TODO: implement
	return nil, fmt.Errorf("UpdateStatus: %w", models.ErrInternal)
}

func (s *Service) Debit(_ context.Context, _ models.DebitRequest) (*models.DebitResult, error) {
	// TODO: implement
	return nil, fmt.Errorf("Debit: %w", models.ErrInternal)
}

func (s *Service) Credit(_ context.Context, _ models.CreditRequest) (*models.CreditResult, error) {
	// TODO: implement
	return nil, fmt.Errorf("Credit: %w", models.ErrInternal)
}
