package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.openly.dev/pointy"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccessManagementService interface {
	CreateUser(ctx context.Context, userCreate models.CreateUser) (*models.User, error)
}

type serverAPI struct {
	authpb.UnimplementedAuthServiceServer
	authpb.UnimplementedAccessManagementServiceServer

	accessManagementSvc AccessManagementService
}

func NewServerAPI(gRPC *grpc.Server, accessManagementSvc AccessManagementService) {
	srv := &serverAPI{
		accessManagementSvc: accessManagementSvc,
	}

	authpb.RegisterAuthServiceServer(gRPC, srv)
	authpb.RegisterAccessManagementServiceServer(gRPC, srv)
}

func (s *serverAPI) Login(_ context.Context, _ *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Logout(_ context.Context, _ *authpb.LogoutRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (s *serverAPI) RefreshToken(_ context.Context, _ *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	panic("implement me")
}

func (s *serverAPI) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	if req.GetLogin() == "" {
		return nil, models.NewInvalidParamsError("login", "is empty")
	}

	if req.GetPassword() == "" {
		return nil, models.NewInvalidParamsError("password", "is empty")
	}

	if req.GetEmail() == "" && req.GetPhone() == "" {
		return nil, models.NewInvalidParamsError("email and phone", "at least one field is required")
	}

	if len(req.GetRoleCodes()) == 0 {
		return nil, models.NewInvalidParamsError("role_codes", "is empty")
	}

	user := models.CreateUser{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
		Role:     req.GetRoleCodes(),
	}

	if req.GetEmail() != "" {
		user.Email = pointy.String(req.GetEmail())
	}

	if req.GetPhone() != "" {
		user.Phone = pointy.String(req.GetPhone())
	}

	userResp, err := s.accessManagementSvc.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("s.accessManagementSvc.CreateUser: %w", err)
	}

	return &authpb.CreateUserResponse{
		User: &authpb.User{
			Id:        userResp.PublicID.String(),
			Status:    toProtoStatus(userResp.Status),
			Login:     userResp.Login,
			Email:     userResp.Email,
			Phone:     userResp.Phone,
			CreatedAt: timestamppb.New(userResp.CreatedAt),
			UpdatedAt: timestamppb.New(userResp.UpdatedAt),
		},
	}, nil
}

func toProtoStatus(s models.Status) authpb.UserStatus {
	switch s {
	case models.StatusActive:
		return authpb.UserStatus_USER_STATUS_ACTIVE
	case models.StatusBlocked:
		return authpb.UserStatus_USER_STATUS_BLOCKED
	case models.StatusDisabled:
		return authpb.UserStatus_USER_STATUS_DISABLED
	case models.StatusLocked:
		return authpb.UserStatus_USER_STATUS_LOCKED
	default:
		return authpb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}
