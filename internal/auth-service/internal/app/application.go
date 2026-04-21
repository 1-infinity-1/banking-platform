package app

import (
	"context"
	"fmt"
	"log/slog"

	grpcapp "github.com/1-infinity-1/banking-platform/internal/auth-service/internal/app/grpc"
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/config"
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/services/management"
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/storage/role"
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/storage/tx"
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/storage/user"
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

	txManager := tx.NewTxManager(conn)
	roleRepo := role.NewRepository(conn)
	userRepo := user.NewRepository(conn)

	accessManagementSvc := management.NewAccessManagementService(txManager, userRepo, roleRepo)

	grpcApp := grpcapp.NewApp(log, cfg.GRPCconfig.Port, conn, accessManagementSvc)

	return &App{
		GRPCSrv: grpcApp,
	}, nil
}
