package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/1-infinity-1/banking-platform/internal/auth-service/internal/transport/grpc"
	internal_interceptor "github.com/1-infinity-1/banking-platform/internal/auth-service/internal/transport/grpc/interceptor"
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

func NewApp(log *slog.Logger, port int, conn *postgres.Conn, accessManagementSvc authgrpc.AccessManagementService, authSvc authgrpc.AuthService) *App {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.TraceUnaryServerInterceptor(),
			interceptor.LoggingUnaryServerInterceptor(log),
			internal_interceptor.UnaryErrorInterceptor(log),
			interceptor.RecoveryUnaryServerInterceptor(log),
		),
	)

	authgrpc.NewServerAPI(gRPCServer, accessManagementSvc, authSvc)

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server is running", slog.String("address", lis.Addr().String()))

	go func() {
		<-ctx.Done()
		a.gRPCServer.GracefulStop()
		a.conn.Close()
	}()

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
