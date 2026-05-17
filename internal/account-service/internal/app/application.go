package app

import (
	"context"
	"fmt"
	"log/slog"

	grpcapp "github.com/1-infinity-1/banking-platform/internal/account-service/internal/app/grpc"
	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/config"
	accountsvc "github.com/1-infinity-1/banking-platform/internal/account-service/internal/services/account"
	accountrepo "github.com/1-infinity-1/banking-platform/internal/account-service/internal/storage/account"
	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/storage/tx"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func NewApp(ctx context.Context, log *slog.Logger, cfg config.Config) (*App, error) {
	conn, err := postgres.NewDB(ctx, postgres.Cfg{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
	})
	if err != nil {
		return nil, fmt.Errorf("postgres.NewDB: %w", err)
	}

	txManager := tx.NewManager(conn)
	repo := accountrepo.NewRepository(conn)
	svc := accountsvc.NewService(txManager, repo)

	grpcApp := grpcapp.NewApp(log, cfg.GRPCconfig.Port, conn, svc)

	return &App{
		GRPCSrv: grpcApp,
	}, nil
}
