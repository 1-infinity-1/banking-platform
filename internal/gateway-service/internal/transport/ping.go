package transport

import (
	"context"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
)

const (
	pong = "pong"
)

func (g *GatewayHandler) Ping(ctx context.Context) (*api.PingResponse, error) {
	return &api.PingResponse{
		Message: pong,
	}, nil
}
