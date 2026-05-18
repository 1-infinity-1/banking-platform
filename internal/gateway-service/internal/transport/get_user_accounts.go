package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) GetUserAccounts(
	ctx context.Context,
	params api.GetUserAccountsParams,
) (api.GetUserAccountsRes, error) {
	accounts, err := g.accountSvc.GetUserAccounts(ctx, models.GetUserAccountsParams{UserID: params.UserID})
	if err != nil {
		return nil, fmt.Errorf("g.accountSvc.GetUserAccounts: %w", err)
	}

	apiAccounts := make([]api.Account, 0, len(accounts))
	for _, a := range accounts {
		apiAccounts = append(apiAccounts, toAPIAccount(a))
	}

	return &api.GetUserAccountsResponse{Accounts: apiAccounts}, nil
}
