package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) CreateAccount(
	ctx context.Context,
	req *api.CreateAccountRequest,
) (api.CreateAccountRes, error) {
	if req.Currency == "" {
		return nil, models.NewValidationError("currency", "is required", nil)
	}

	params := models.CreateAccountParams{
		UserID:   req.UserID,
		Currency: req.Currency,
	}

	account, err := g.accountSvc.CreateAccount(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("g.accountSvc.CreateAccount: %w", err)
	}

	res := toAPIAccount(*account)
	return &res, nil
}
