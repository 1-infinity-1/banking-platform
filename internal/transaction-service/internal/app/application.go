package app

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	grpcapp "github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/app/grpc"
	accountclient "github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/clients/account"
	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/config"
	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/kafka"
	transactionsvc "github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/services/transaction"
	transactionrepo "github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/storage/transaction"
	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/storage/tx"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	brokers := strings.Split(cfg.Kafka.Brokers, ",")
	kafkaProducer := kafka.NewProducer(brokers, cfg.Kafka.Topic)

	accountGRPCConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.AccountService.Host, cfg.AccountService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient (account-service): %w", err)
	}

	accClient := accountclient.NewClient(accountGRPCConn)

	txManager := tx.NewManager(conn)
	txRepo := transactionrepo.NewRepository(conn)
	txSvc := transactionsvc.NewService(txManager, txRepo, accClient, kafkaProducer)

	grpcApp := grpcapp.NewApp(log, cfg.GRPCconfig.Port, conn, txSvc)

	return &App{
		GRPCSrv: grpcApp,
	}, nil
}
