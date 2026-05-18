package transaction

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (c *Client) Transfer(ctx context.Context, params models.TransferParams) (*models.Transaction, error) {
	resp, err := c.svc.Transfer(ctx, toProtoTransferRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	tx, err := toTransaction(resp)
	if err != nil {
		return nil, fmt.Errorf("toTransaction: %w", err)
	}
	return tx, nil
}

func (c *Client) Replenish(ctx context.Context, params models.ReplenishParams) (*models.Transaction, error) {
	resp, err := c.svc.Replenish(ctx, toProtoReplenishRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	tx, err := toTransaction(resp)
	if err != nil {
		return nil, fmt.Errorf("toTransaction: %w", err)
	}
	return tx, nil
}

func (c *Client) GetHistory(ctx context.Context, params models.GetHistoryParams) ([]models.Transaction, error) {
	resp, err := c.svc.GetHistory(ctx, toProtoGetHistoryRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	txs := make([]models.Transaction, 0, len(resp.GetTransactions()))
	for _, t := range resp.GetTransactions() {
		tx, mappingErr := toTransaction(t)
		if mappingErr != nil {
			return nil, fmt.Errorf("toTransaction: %w", mappingErr)
		}
		txs = append(txs, *tx)
	}
	return txs, nil
}

func (c *Client) GetTransaction(ctx context.Context, params models.GetTransactionParams) (*models.Transaction, error) {
	resp, err := c.svc.GetTransaction(ctx, toProtoGetTransactionRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	tx, err := toTransaction(resp)
	if err != nil {
		return nil, fmt.Errorf("toTransaction: %w", err)
	}
	return tx, nil
}
