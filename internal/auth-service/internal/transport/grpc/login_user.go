package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverAPI) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req.GetLogin() == "" {
		return nil, models.NewInvalidParamsError("login", "is empty")
	}

	if req.GetPassword() == "" {
		return nil, models.NewInvalidParamsError("password", "is empty")
	}

	if req.GetContext() == nil {
		return nil, models.NewInvalidParamsError("context", "is null")
	}

	if req.GetContext().GetUserAgent() == "" ||
		req.GetContext().GetPlatform() == "" {
		return nil, models.NewInvalidParamsError("context", "not completely filled")
	}

	credentials := models.LoginCredentials{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
		Context: models.InputContext{
			UserAgent: req.GetContext().GetUserAgent(),
			Platform:  req.GetContext().GetPlatform(),
		},
	}

	resp, err := s.authSvc.Login(ctx, credentials)
	if err != nil {
		return nil, fmt.Errorf("s.authSvc.Login: %w", err)
	}

	return &authpb.LoginResponse{
		User: &authpb.User{
			Id:        resp.User.PublicID.String(),
			Status:    toProtoUserStatus(resp.User.Status),
			Login:     resp.User.Login,
			Email:     resp.User.Email,
			Phone:     resp.User.Phone,
			CreatedAt: timestamppb.New(resp.User.CreatedAt),
			UpdatedAt: timestamppb.New(resp.User.UpdatedAt),
		},
		Session: &authpb.Session{
			Id:     resp.Session.PublicID.String(),
			UserId: resp.User.PublicID.String(),
			Status: toProtoSessionStatus(resp.Session.Status),
			Device: &authpb.Device{
				Id:        resp.Device.PublicID.String(),
				Platform:  resp.Device.Platform,
				UserAgent: resp.Device.UserAgent,
			},
			CreatedAt:  timestamppb.New(resp.Session.CreatedAt),
			UpdatedAt:  timestamppb.New(resp.Session.UpdatedAt),
			ExpiresAt:  timestamppb.New(resp.Session.ExpiresAt),
			LastSeenAt: timestamppb.New(resp.Session.LastSeenAt),
		},
		Tokens: &authpb.TokenPair{
			AccessToken:           resp.Tokens.AccessToken,
			RefreshToken:          resp.Tokens.RefreshToken,
			TokenType:             resp.Tokens.TypeToken,
			AccessTokenExpiresAt:  timestamppb.New(resp.Tokens.AccessTokenExpiresAt),
			RefreshTokenExpiresAt: timestamppb.New(resp.Tokens.RefreshTokenExpiresAt),
		},
		AuthContext: &authpb.AuthContext{
			UserId:          resp.User.PublicID.String(),
			SessionId:       resp.Session.PublicID.String(),
			RoleCodes:       resp.User.Roles,
			PermissionCodes: resp.User.Permissions,
		},
	}, nil
}
