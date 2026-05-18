package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
	"github.com/google/uuid"
)

func (s *serverAPI) GetHistory(
	ctx context.Context,
	req *transactionpb.GetHistoryRequest,
) (*transactionpb.TransactionsList, error) {
	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, models.NewInvalidParamsError("account_id", "must be a valid UUID")
	}

	transactions, err := s.svc.GetHistory(ctx, models.GetHistoryRequest{
		AccountID: accountID,
		Limit:     req.GetLimit(),
		Offset:    req.GetOffset(),
	})
	if err != nil {
		return nil, fmt.Errorf("s.svc.GetHistory: %w", err)
	}

	return toProtoTransactionsList(transactions), nil
}
