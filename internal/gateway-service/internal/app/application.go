package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/app/srv"
	authclient "github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/clients/auth"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/config"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/services/management"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/transport"
)

type App struct {
	httpSrv    *srv.App
	authClient *authclient.Client
}

func NewApp(log *slog.Logger, cfg config.Config) (*App, error) {
	authClient, err := authclient.NewClient(cfg.AuthGRPC.Host, cfg.AuthGRPC.Port)
	if err != nil {
		return nil, fmt.Errorf("authclient.NewClient: %w", err)
	}

	managementSvc := management.New(authClient)

	hnd := transport.NewGatewayHandler(managementSvc)

	httpServer, err := srv.NewApp(hnd, log, cfg.HTTPConfig.Port)
	if err != nil {
		_ = authClient.Close()
		return nil, fmt.Errorf("srv.NewApp: %w", err)
	}

	return &App{
		httpSrv:    httpServer,
		authClient: authClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		_ = a.authClient.Close()
	}()

	return a.httpSrv.Run(ctx)
}
