package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/google/uuid"
)

func (s *serverAPI) GetUserAccounts(
	ctx context.Context,
	req *accountpb.GetUserAccountsRequest,
) (*accountpb.AccountsList, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, models.NewInvalidParamsError("user_id", "must be a valid UUID")
	}

	accounts, err := s.svc.GetUserAccounts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("s.svc.GetUserAccounts: %w", err)
	}
	return toProtoAccountsList(accounts), nil
}
