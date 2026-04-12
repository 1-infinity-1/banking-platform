package grpc

import (
	"context"

	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverAPI struct {
	authpb.UnimplementedAuthServiceServer
	authpb.UnimplementedAccessManagementServiceServer

	//TODO: имплиментируем logger, service и т.д.
}

func NewServerAPI(gRPC *grpc.Server) {
	srv := &serverAPI{}

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

func (s *serverAPI) CreateUser(_ context.Context, _ *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	panic("implement me")
}
