package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) UpdateAccountStatus(
	ctx context.Context,
	req *api.UpdateAccountStatusRequest,
	params api.UpdateAccountStatusParams,
) (api.UpdateAccountStatusRes, error) {
	if req.Status == "" {
		return nil, models.NewValidationError("status", "is required", nil)
	}

	account, err := g.accountSvc.UpdateStatus(ctx, models.UpdateAccountStatusParams{
		AccountID: params.AccountID,
		Status:    models.AccountStatus(req.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("g.accountSvc.UpdateStatus: %w", err)
	}

	res := toAPIAccount(*account)
	return &res, nil
}
