package app

import (
	"log/slog"

	grpcapp "github.com/1-infinity-1/banking-platform/internal/auth-service/internal/app/grpc"
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/config"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func NewApp(log *slog.Logger, cfg config.Config) *App {
	// TODO: implementation repo

	// TODO: implementation service

	grpcApp := grpcapp.NewApp(log, cfg.GRPCconfig.Port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
