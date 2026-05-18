package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) GetTransactionHistory(
	ctx context.Context,
	params api.GetTransactionHistoryParams,
) (api.GetTransactionHistoryRes, error) {
	txs, err := g.transactionSvc.GetHistory(ctx, models.GetHistoryParams{
		AccountID: params.AccountID,
		Limit:     params.Limit.Or(0),
		Offset:    params.Offset.Or(0),
	})
	if err != nil {
		return nil, fmt.Errorf("g.transactionSvc.GetHistory: %w", err)
	}

	apiTxs := make([]api.Transaction, 0, len(txs))
	for _, t := range txs {
		apiTxs = append(apiTxs, toAPITransaction(t))
	}

	return &api.GetTransactionHistoryResponse{Transactions: apiTxs}, nil
}
