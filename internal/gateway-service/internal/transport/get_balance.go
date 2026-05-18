package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) GetBalance(ctx context.Context, params api.GetBalanceParams) (api.GetBalanceRes, error) {
	balance, err := g.accountSvc.GetBalance(ctx, models.GetBalanceParams{AccountID: params.AccountID})
	if err != nil {
		return nil, fmt.Errorf("g.accountSvc.GetBalance: %w", err)
	}

	res := toAPIBalance(*balance)
	return &res, nil
}
