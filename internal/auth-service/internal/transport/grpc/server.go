package grpc

import (
	"context"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
	"google.golang.org/grpc"
)

type AuthService interface {
	Login(ctx context.Context, credentials models.LoginCredentials) (*models.LoginResult, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResult, error)
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
