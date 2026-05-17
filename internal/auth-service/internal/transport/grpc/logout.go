package grpc

import (
	"context"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverAPI) Logout(ctx context.Context, req *authpb.LogoutRequest) (*emptypb.Empty, error) {
	if req.GetRefreshToken() == "" {
		return nil, models.NewInvalidParamsError("refresh_token", "is empty")
	}

	if err := s.authSvc.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
