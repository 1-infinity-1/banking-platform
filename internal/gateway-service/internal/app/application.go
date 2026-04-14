package app

import (
	"fmt"
	"log/slog"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/app/srv"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/config"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/transport"
)

type App struct {
	HTTPSrv *srv.App
}

func NewApp(log *slog.Logger, cfg config.Config) (*App, error) {
	// TODO: implementation repo

	// TODO: implementation service

	hnd := transport.NewGatewayHandler()

	httpServer, err := srv.NewApp(hnd, log, cfg.HTTPConfig.Port)
	if err != nil {
		return nil, fmt.Errorf("srv.NewApp: %w", err)
	}

	return &App{
		HTTPSrv: httpServer,
	}, nil
}
