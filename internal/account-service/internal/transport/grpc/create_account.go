package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/google/uuid"
)

func (s *serverAPI) CreateAccount(
	ctx context.Context,
	req *accountpb.CreateAccountRequest,
) (*accountpb.Account, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, models.NewInvalidParamsError("user_id", "must be a valid UUID")
	}

	acc, err := s.svc.CreateAccount(ctx, models.CreateAccountRequest{
		UserID:   userID,
		Currency: req.GetCurrency(),
	})
	if err != nil {
		return nil, fmt.Errorf("s.svc.CreateAccount: %w", err)
	}
	return toProtoAccount(acc), nil
}
