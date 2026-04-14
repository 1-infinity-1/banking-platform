package srv

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/transport"
)

type App struct {
	log        *slog.Logger
	ogenServer *api.Server
	port       string
}

func NewApp(hnd *transport.GatewayHandler, log *slog.Logger, port string) (*App, error) {
	srv, err := api.NewServer(hnd)
	if err != nil {
		return nil, fmt.Errorf("api.NewServer: %w", err)
	}

	return &App{
		log:        log,
		ogenServer: srv,
		port:       port,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	const op = "httpApp.Run"

	a.log.With("op", op)

	srv := &http.Server{
		Addr:    ":" + a.port,
		Handler: a.ogenServer,
	}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	a.log.Info("api server is running", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
