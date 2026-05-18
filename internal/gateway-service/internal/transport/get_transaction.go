package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) GetTransaction(
	ctx context.Context,
	params api.GetTransactionParams,
) (api.GetTransactionRes, error) {
	tx, err := g.transactionSvc.GetTransaction(ctx, models.GetTransactionParams{TransactionID: params.TransactionID})
	if err != nil {
		return nil, fmt.Errorf("g.transactionSvc.GetTransaction: %w", err)
	}

	res := toAPITransaction(*tx)
	return &res, nil
}
