package grpc

import (
	"context"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService interface {
	Login(ctx context.Context, credentials models.LoginCredentials) (*models.LoginResult, error)
}

type AccessManagementService interface {
	CreateUser(ctx context.Context, userCreate models.CreateUser, password string) (*models.User, error)
}

type serverAPI struct {
	authpb.UnimplementedAuthServiceServer
	authpb.UnimplementedAccessManagementServiceServer

	accessManagementSvc AccessManagementService
	authSvc             AuthService
}

func NewServerAPI(gRPC *grpc.Server, accessManagementSvc AccessManagementService, authSvc AuthService) {
	srv := &serverAPI{
		accessManagementSvc: accessManagementSvc,
		authSvc:             authSvc,
	}

	authpb.RegisterAuthServiceServer(gRPC, srv)
	authpb.RegisterAccessManagementServiceServer(gRPC, srv)
}

func (s *serverAPI) Logout(_ context.Context, _ *authpb.LogoutRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (s *serverAPI) RefreshToken(_ context.Context, _ *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	panic("implement me")
}
