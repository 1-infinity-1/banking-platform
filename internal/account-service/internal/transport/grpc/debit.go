package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *serverAPI) Debit(ctx context.Context, req *accountpb.DebitRequest) (*accountpb.DebitResponse, error) {
	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, models.NewInvalidParamsError("account_id", "must be a valid UUID")
	}
	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		return nil, models.NewInvalidParamsError("amount", "must be a valid decimal string")
	}

	result, err := s.svc.Debit(ctx, models.DebitRequest{
		AccountID:      accountID,
		Amount:         amount,
		IdempotencyKey: req.GetIdempotencyKey(),
	})
	if err != nil {
		return nil, fmt.Errorf("s.svc.Debit: %w", err)
	}
	return toProtoDebitResponse(result), nil
}
