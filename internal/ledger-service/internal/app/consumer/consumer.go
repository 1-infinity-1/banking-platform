package consumerapp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/config"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/services/ledger"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/storage/entry"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/storage/tx"
	kafkatransport "github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/transport/kafka"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
)

type App struct {
	consumer *kafkatransport.Consumer
	conn     *postgres.Conn
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
	entryRepo := entry.NewRepository(conn)
	ledgerSvc := ledger.NewService(txManager, entryRepo)
	consumer := kafkatransport.NewConsumer(cfg, ledgerSvc, log)

	return &App{
		consumer: consumer,
		conn:     conn,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.conn.Close()
	return a.consumer.Run(ctx)
}
