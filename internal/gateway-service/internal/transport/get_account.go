package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) GetAccount(ctx context.Context, params api.GetAccountParams) (api.GetAccountRes, error) {
	account, err := g.accountSvc.GetAccount(ctx, models.GetAccountParams{AccountID: params.AccountID})
	if err != nil {
		return nil, fmt.Errorf("g.accountSvc.GetAccount: %w", err)
	}

	res := toAPIAccount(*account)
	return &res, nil
}
