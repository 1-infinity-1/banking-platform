package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
	"github.com/google/uuid"
)

func (s *serverAPI) GetTransaction(
	ctx context.Context,
	req *transactionpb.GetTransactionRequest,
) (*transactionpb.Transaction, error) {
	id, err := uuid.Parse(req.GetTransactionId())
	if err != nil {
		return nil, models.NewInvalidParamsError("transaction_id", "must be a valid UUID")
	}

	transaction, err := s.svc.GetTransaction(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.svc.GetTransaction: %w", err)
	}

	return toProtoTransaction(transaction), nil
}
