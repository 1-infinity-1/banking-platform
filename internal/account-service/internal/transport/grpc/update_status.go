package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/google/uuid"
)

func (s *serverAPI) UpdateStatus(
	ctx context.Context,
	req *accountpb.UpdateStatusRequest,
) (*accountpb.UpdateStatusResponse, error) {
	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, models.NewInvalidParamsError("account_id", "must be a valid UUID")
	}

	updated, err := s.svc.UpdateStatus(ctx, models.UpdateStatusRequest{
		AccountID: accountID,
		Status:    fromProtoAccountStatus(req.GetStatus()),
	})
	if err != nil {
		return nil, fmt.Errorf("s.svc.UpdateStatus: %w", err)
	}
	return toProtoUpdateStatusResponse(updated), nil
}
