package account

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
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

func (s *Service) CreateAccount(ctx context.Context, req models.CreateAccountRequest) (*models.Account, error) {
	if req.UserID == uuid.Nil {
		return nil, models.NewInvalidParamsError("user_id", "must be non-zero UUID")
	}
	if req.Currency == "" {
		return nil, models.NewInvalidParamsError("currency", "must be non-empty")
	}

	var account *models.Account
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		account, inner = s.accountRepo.CreateAccountTx(ctx, tx, req)
		if inner != nil {
			return fmt.Errorf("s.accountRepo.CreateAccountTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return account, nil
}

func (s *Service) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*models.Account, error) {
	if userID == uuid.Nil {
		return nil, models.NewInvalidParamsError("user_id", "must be non-zero UUID")
	}

	accounts, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("s.accountRepo.GetByUserID: %w", err)
	}
	return accounts, nil
}

func (s *Service) GetAccount(ctx context.Context, accountID uuid.UUID) (*models.Account, error) {
	if accountID == uuid.Nil {
		return nil, models.NewInvalidParamsError("account_id", "must be non-zero UUID")
	}

	var account *models.Account
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		account, inner = s.accountRepo.GetByIDTx(ctx, tx, accountID)
		if inner != nil {
			return fmt.Errorf("s.accountRepo.GetByIDTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return account, nil
}

func (s *Service) GetBalance(ctx context.Context, accountID uuid.UUID) (*models.Balance, error) {
	if accountID == uuid.Nil {
		return nil, models.NewInvalidParamsError("account_id", "must be non-zero UUID")
	}

	var account *models.Account
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		account, inner = s.accountRepo.GetByIDTx(ctx, tx, accountID)
		if inner != nil {
			return fmt.Errorf("s.accountRepo.GetByIDTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return &models.Balance{
		AccountID: account.PublicID.String(),
		Amount:    account.Balance,
		Currency:  account.Currency,
	}, nil
}

func (s *Service) UpdateStatus(ctx context.Context, req models.UpdateStatusRequest) (*models.Account, error) {
	if req.AccountID == uuid.Nil {
		return nil, models.NewInvalidParamsError("account_id", "must be non-zero UUID")
	}
	switch req.Status {
	case models.AccountStatusActive, models.AccountStatusBlocked, models.AccountStatusClosed:
	case models.AccountStatusUnspecified:
		return nil, models.NewInvalidParamsError("status", "must be active|blocked|closed")
	default:
		return nil, models.NewInvalidParamsError("status", "must be active|blocked|closed")
	}

	var updated *models.Account
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		updated, inner = s.accountRepo.UpdateStatusTx(ctx, tx, req)
		if inner != nil {
			return fmt.Errorf("s.accountRepo.UpdateStatusTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return updated, nil
}

func (s *Service) Debit(ctx context.Context, req models.DebitRequest) (*models.DebitResult, error) {
	if err := validateMoneyOp(req.AccountID, req.Amount, req.IdempotencyKey); err != nil {
		return nil, err
	}

	var result *models.DebitResult
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		result, inner = s.accountRepo.DebitTx(ctx, tx, req)
		if inner != nil {
			return fmt.Errorf("s.accountRepo.DebitTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return result, nil
}

func (s *Service) Credit(ctx context.Context, req models.CreditRequest) (*models.CreditResult, error) {
	if err := validateMoneyOp(req.AccountID, req.Amount, req.IdempotencyKey); err != nil {
		return nil, err
	}

	var result *models.CreditResult
	err := s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		var inner error
		result, inner = s.accountRepo.CreditTx(ctx, tx, req)
		if inner != nil {
			return fmt.Errorf("s.accountRepo.CreditTx: %w", inner)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}
	return result, nil
}

func validateMoneyOp(accountID uuid.UUID, amount decimal.Decimal, idempotencyKey string) error {
	if accountID == uuid.Nil {
		return models.NewInvalidParamsError("account_id", "must be non-zero UUID")
	}
	if !amount.IsPositive() {
		return models.NewInvalidParamsError("amount", "must be positive")
	}
	if idempotencyKey == "" {
		return models.NewInvalidParamsError("idempotency_key", "must be non-empty")
	}
	return nil
}
