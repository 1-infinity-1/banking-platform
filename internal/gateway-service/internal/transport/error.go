package transport

import (
	"context"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
)

func (g *GatewayHandler) NewError(_ context.Context, _ error) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{}
}
