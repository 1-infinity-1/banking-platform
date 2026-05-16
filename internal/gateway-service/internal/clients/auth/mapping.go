package auth

import (
	"errors"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toProtoCreateUserRequest(p models.CreateUserParams) *authpb.CreateUserRequest {
	return &authpb.CreateUserRequest{
		Login:     p.Login,
		Email:     p.Email,
		Phone:     p.Phone,
		Password:  p.Password,
		RoleCodes: p.RoleCodes,
	}
}

func toUser(u *authpb.User) (*models.User, error) {
	if u == nil {
		return nil, errors.New("empty user")
	}
	id, err := uuid.Parse(u.GetId())
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	user := &models.User{
		ID:        id,
		Login:     u.GetLogin(),
		Status:    toUserStatus(u.GetStatus()),
		CreatedAt: u.GetCreatedAt().AsTime(),
		UpdatedAt: u.GetUpdatedAt().AsTime(),
	}
	if v := u.GetEmail(); v != "" {
		user.Email = &v
	}
	if v := u.GetPhone(); v != "" {
		user.Phone = &v
	}
	return user, nil
}

func toUserStatus(s authpb.UserStatus) models.UserStatus {
	switch s {
	case authpb.UserStatus_USER_STATUS_ACTIVE:
		return models.UserStatusActive
	case authpb.UserStatus_USER_STATUS_BLOCKED:
		return models.UserStatusBlocked
	case authpb.UserStatus_USER_STATUS_LOCKED:
		return models.UserStatusLocked
	case authpb.UserStatus_USER_STATUS_DISABLED:
		return models.UserStatusDisabled
	case authpb.UserStatus_USER_STATUS_UNSPECIFIED:
		return models.UserStatusUnspecified
	default:
		return models.UserStatusUnspecified
	}
}

func mapGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("unexpected gRPC error: %w", err)
	}
	switch st.Code() { //nolint:exhaustive // only meaningful codes handled; default covers the rest
	case codes.AlreadyExists:
		return models.NewConflictError(st.Message(), err)
	case codes.InvalidArgument:
		return models.NewValidationError("", st.Message(), err)
	case codes.NotFound:
		return models.NewNotFoundError(st.Message(), err)
	default:
		return fmt.Errorf("auth service error: %w", err)
	}
}
