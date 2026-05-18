package app

import (
	"context"
	"fmt"
	"log/slog"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/app/srv"
	accountclient "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/clients/account"
	authclient "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/clients/auth"
	ledgerclient "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/clients/ledger"
	transactionclient "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/clients/transaction"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/config"
	accountsvc "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/services/account"
	authsvc "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/services/auth"
	ledgersvc "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/services/ledger"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/services/management"
	transactionsvc "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/services/transaction"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/transport"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/transport/middleware"
)

type App struct {
	httpSrv           *srv.App
	authClient        *authclient.Client
	accountClient     *accountclient.Client
	transactionClient *transactionclient.Client
	ledgerClient      *ledgerclient.Client
}

func NewApp(log *slog.Logger, cfg config.Config) (*App, error) {
	authClient, err := authclient.NewClient(cfg.AuthGRPC.Host, cfg.AuthGRPC.Port)
	if err != nil {
		return nil, fmt.Errorf("authclient.NewClient: %w", err)
	}

	accountClient, err := accountclient.NewClient(cfg.AccountGRPC.Host, cfg.AccountGRPC.Port)
	if err != nil {
		_ = authClient.Close()
		return nil, fmt.Errorf("accountclient.NewClient: %w", err)
	}

	transactionClient, err := transactionclient.NewClient(cfg.TransactionGRPC.Host, cfg.TransactionGRPC.Port)
	if err != nil {
		_ = authClient.Close()
		_ = accountClient.Close()
		return nil, fmt.Errorf("transactionclient.NewClient: %w", err)
	}

	ledgerClient, err := ledgerclient.NewClient(cfg.LedgerGRPC.Host, cfg.LedgerGRPC.Port)
	if err != nil {
		_ = authClient.Close()
		_ = accountClient.Close()
		_ = transactionClient.Close()
		return nil, fmt.Errorf("ledgerclient.NewClient: %w", err)
	}

	managementSvc := management.New(authClient)
	authService := authsvc.New(authClient)
	accountService := accountsvc.New(accountClient)
	transactionService := transactionsvc.New(transactionClient)
	ledgerService := ledgersvc.New(ledgerClient)

	hnd := transport.NewGatewayHandler(managementSvc, authService, accountService, transactionService, ledgerService)

	// TODO: roles need to be mapped to API handlers
	operationRoles := map[api.OperationName][]string{
		api.UpdateAccountStatusOperation: {"admin"},
	}
	sec := middleware.NewJWTSecurityHandler(cfg.JWT.Secret, operationRoles)

	httpServer, err := srv.NewApp(hnd, sec, log, cfg.HTTPConfig.Port)
	if err != nil {
		_ = authClient.Close()
		_ = accountClient.Close()
		_ = transactionClient.Close()
		_ = ledgerClient.Close()
		return nil, fmt.Errorf("srv.NewApp: %w", err)
	}

	return &App{
		httpSrv:           httpServer,
		authClient:        authClient,
		accountClient:     accountClient,
		transactionClient: transactionClient,
		ledgerClient:      ledgerClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		_ = a.authClient.Close()
		_ = a.accountClient.Close()
		_ = a.transactionClient.Close()
		_ = a.ledgerClient.Close()
	}()

	return a.httpSrv.Run(ctx)
}
