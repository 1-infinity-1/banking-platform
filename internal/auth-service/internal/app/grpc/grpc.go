package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/1-infinity-1/banking-platform/internal/auth-service/internal/transport/grpc"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.NewServerAPI(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
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
	}()

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
