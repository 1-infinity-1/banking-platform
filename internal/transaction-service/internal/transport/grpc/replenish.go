package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *serverAPI) Replenish(
	ctx context.Context,
	req *transactionpb.ReplenishRequest,
) (*transactionpb.Transaction, error) {
	to, err := uuid.Parse(req.GetToAccountId())
	if err != nil {
		return nil, models.NewInvalidParamsError("to_account_id", "must be a valid UUID")
	}

	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		return nil, models.NewInvalidParamsError("amount", "must be a valid decimal string")
	}

	transaction, err := s.svc.Replenish(ctx, models.ReplenishRequest{
		ToAccountID:    to,
		Amount:         amount,
		Currency:       req.GetCurrency(),
		IdempotencyKey: req.GetIdempotencyKey(),
	})
	if err != nil {
		return nil, fmt.Errorf("s.svc.Replenish: %w", err)
	}

	return toProtoTransaction(transaction), nil
}
