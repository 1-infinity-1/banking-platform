package transport

import (
	"context"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

type managementService interface {
	CreateUser(ctx context.Context, params models.CreateUserParams) (*models.User, error)
}

type authService interface {
	Login(ctx context.Context, params models.LoginParams) (*models.LoginResult, error)
	Logout(ctx context.Context, params models.LogoutParams) error
	RefreshToken(ctx context.Context, params models.RefreshTokenParams) (*models.RefreshTokenResult, error)
}

type accountService interface {
	CreateAccount(ctx context.Context, params models.CreateAccountParams) (*models.Account, error)
	GetUserAccounts(ctx context.Context, params models.GetUserAccountsParams) ([]models.Account, error)
	GetAccount(ctx context.Context, params models.GetAccountParams) (*models.Account, error)
	GetBalance(ctx context.Context, params models.GetBalanceParams) (*models.Balance, error)
	UpdateStatus(ctx context.Context, params models.UpdateAccountStatusParams) (*models.Account, error)
}

type transactionService interface {
	Transfer(ctx context.Context, params models.TransferParams) (*models.Transaction, error)
	Replenish(ctx context.Context, params models.ReplenishParams) (*models.Transaction, error)
	GetHistory(ctx context.Context, params models.GetHistoryParams) ([]models.Transaction, error)
	GetTransaction(ctx context.Context, params models.GetTransactionParams) (*models.Transaction, error)
}

type ledgerService interface {
	GetStatement(ctx context.Context, params models.GetStatementParams) (*models.Statement, error)
}

type GatewayHandler struct {
	api.UnimplementedHandler

	managementSvc  managementService
	authSvc        authService
	accountSvc     accountService
	transactionSvc transactionService
	ledgerSvc      ledgerService
}

func NewGatewayHandler(
	managementSvc managementService,
	authSvc authService,
	accountSvc accountService,
	transactionSvc transactionService,
	ledgerSvc ledgerService,
) *GatewayHandler {
	return &GatewayHandler{
		managementSvc:  managementSvc,
		authSvc:        authSvc,
		accountSvc:     accountSvc,
		transactionSvc: transactionSvc,
		ledgerSvc:      ledgerSvc,
	}
}
