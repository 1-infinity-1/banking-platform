package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) Replenish(ctx context.Context, req *api.ReplenishRequest) (api.ReplenishRes, error) {
	if req.Amount == "" {
		return nil, models.NewValidationError("amount", "is required", nil)
	}
	if req.Currency == "" {
		return nil, models.NewValidationError("currency", "is required", nil)
	}
	if req.IdempotencyKey == "" {
		return nil, models.NewValidationError("idempotency_key", "is required", nil)
	}

	tx, err := g.transactionSvc.Replenish(ctx, models.ReplenishParams{
		ToAccountID:    req.ToAccountID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, fmt.Errorf("g.transactionSvc.Replenish: %w", err)
	}

	res := toAPITransaction(*tx)
	return &res, nil
}
