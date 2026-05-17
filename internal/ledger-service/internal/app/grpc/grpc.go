package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/config"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/services/ledger"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/storage/entry"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/storage/tx"
	ledgergrpc "github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/transport/grpc"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/transport/grpc/interceptor"
	pkginterceptor "github.com/1-infinity-1/banking-platform/pkg/grpc/interceptor"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	conn       *postgres.Conn
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

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			pkginterceptor.TraceUnaryServerInterceptor(),
			pkginterceptor.LoggingUnaryServerInterceptor(log),
			interceptor.UnaryErrorInterceptor(log),
			pkginterceptor.RecoveryUnaryServerInterceptor(log),
		),
	)

	ledgergrpc.NewServerAPI(gRPCServer, ledgerSvc)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfg.GRPCConfig.Port,
		conn:       conn,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	const op = "grpcapp.Run"

	lc := net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.InfoContext(ctx, "grpc server is running", slog.String("address", lis.Addr().String()))

	go func() {
		<-ctx.Done()
		a.gRPCServer.GracefulStop()
		a.conn.Close()
	}()

	if err = a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
