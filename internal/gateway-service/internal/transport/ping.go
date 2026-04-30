package transport

import (
	"context"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
)

const (
	pong = "pong"
)

func (g *GatewayHandler) Ping(_ context.Context) (*api.PingResponse, error) {
	return &api.PingResponse{
		Message: pong,
	}, nil
}
