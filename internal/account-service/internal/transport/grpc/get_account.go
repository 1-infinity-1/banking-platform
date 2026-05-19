package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/google/uuid"
)

func (s *serverAPI) GetAccount(ctx context.Context, req *accountpb.GetAccountRequest) (*accountpb.Account, error) {
	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, models.NewInvalidParamsError("account_id", "must be a valid UUID")
	}

	acc, err := s.svc.GetAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("s.svc.GetAccount: %w", err)
	}
	return toProtoAccount(acc), nil
}
