package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/google/uuid"
)

func (s *serverAPI) GetBalance(ctx context.Context, req *accountpb.GetBalanceRequest) (*accountpb.Balance, error) {
	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, models.NewInvalidParamsError("account_id", "must be a valid UUID")
	}

	balance, err := s.svc.GetBalance(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("s.svc.GetBalance: %w", err)
	}
	return toProtoBalance(balance), nil
}
