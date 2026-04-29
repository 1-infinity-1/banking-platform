package grpc

import (
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	authpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/auth"
)

func toProtoSessionStatus(s models.SessionStatus) authpb.SessionStatus {
	switch s {
	case models.SessionStatusActive:
		return authpb.SessionStatus_SESSION_STATUS_ACTIVE
	case models.SessionStatusRevoked:
		return authpb.SessionStatus_SESSION_STATUS_REVOKED
	case models.SessionStatusExpired:
		return authpb.SessionStatus_SESSION_STATUS_EXPIRED
	default:
		return authpb.SessionStatus_SESSION_STATUS_UNSPECIFIED
	}
}

func toProtoUserStatus(s models.UserStatus) authpb.UserStatus {
	switch s {
	case models.UserStatusActive:
		return authpb.UserStatus_USER_STATUS_ACTIVE
	case models.UserStatusBlocked:
		return authpb.UserStatus_USER_STATUS_BLOCKED
	case models.UserStatusDisabled:
		return authpb.UserStatus_USER_STATUS_DISABLED
	case models.UserStatusLocked:
		return authpb.UserStatus_USER_STATUS_LOCKED
	default:
		return authpb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}
