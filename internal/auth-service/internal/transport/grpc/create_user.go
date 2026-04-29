package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
	"go.openly.dev/pointy"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
		Login: req.GetLogin(),
		Role:  req.GetRoleCodes(),
	}

	if req.GetEmail() != "" {
		user.Email = pointy.String(req.GetEmail())
	}

	if req.GetPhone() != "" {
		user.Phone = pointy.String(req.GetPhone())
	}

	userResp, err := s.accessManagementSvc.CreateUser(ctx, user, req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("s.accessManagementSvc.CreateUser: %w", err)
	}

	return &authpb.CreateUserResponse{
		User: &authpb.User{
			Id:        userResp.PublicID.String(),
			Status:    toProtoUserStatus(userResp.Status),
			Login:     userResp.Login,
			Email:     userResp.Email,
			Phone:     userResp.Phone,
			CreatedAt: timestamppb.New(userResp.CreatedAt),
			UpdatedAt: timestamppb.New(userResp.UpdatedAt),
		},
	}, nil
}
