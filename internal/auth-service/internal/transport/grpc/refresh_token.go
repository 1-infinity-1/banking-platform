package grpc

import (
	"context"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
)

func (s *serverAPI) RefreshToken(
	ctx context.Context,
	req *authpb.RefreshTokenRequest,
) (*authpb.RefreshTokenResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, models.NewInvalidParamsError("refresh_token", "is empty")
	}

	result, err := s.authSvc.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &authpb.RefreshTokenResponse{
		Tokens: toProtoTokenPair(result.Tokens),
		AuthContext: &authpb.AuthContext{
			UserId:          result.User.PublicID.String(),
			SessionId:       result.Session.PublicID.String(),
			RoleCodes:       result.User.Roles,
			PermissionCodes: result.User.Permissions,
		},
	}, nil
}
