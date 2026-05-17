package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	transactiongrpc "github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/transport/grpc"
	internal_interceptor "github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/transport/grpc/interceptor"
	"github.com/1-infinity-1/banking-platform/pkg/grpc/interceptor"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	conn       *postgres.Conn
}

func NewApp(
	log *slog.Logger,
	port int,
	conn *postgres.Conn,
	svc transactiongrpc.TransactionService,
) *App {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.TraceUnaryServerInterceptor(),
			interceptor.LoggingUnaryServerInterceptor(log),
			internal_interceptor.UnaryErrorInterceptor(log),
			interceptor.RecoveryUnaryServerInterceptor(log),
		),
	)

	transactiongrpc.NewServerAPI(gRPCServer, svc)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
		conn:       conn,
	}
}

func (a *App) Run(ctx context.Context) error {
	const op = "grpcapp.Run"

	a.log.With("op", op)

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

	err = a.gRPCServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
