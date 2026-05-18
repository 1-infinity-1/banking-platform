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

func toProtoLoginRequest(p models.LoginParams) *authpb.LoginRequest {
	return &authpb.LoginRequest{
		Login:    p.Login,
		Password: p.Password,
		Context: &authpb.RequestContext{
			UserAgent: p.UserAgent,
			Platform:  p.Platform,
		},
	}
}

func toProtoLogoutRequest(p models.LogoutParams) *authpb.LogoutRequest {
	return &authpb.LogoutRequest{
		RefreshToken: p.RefreshToken,
		Context: &authpb.RequestContext{
			UserAgent: p.UserAgent,
			Platform:  p.Platform,
		},
	}
}

func toProtoRefreshTokenRequest(p models.RefreshTokenParams) *authpb.RefreshTokenRequest {
	return &authpb.RefreshTokenRequest{
		RefreshToken: p.RefreshToken,
		Context: &authpb.RequestContext{
			UserAgent: p.UserAgent,
			Platform:  p.Platform,
		},
	}
}

func toLoginResult(resp *authpb.LoginResponse) (*models.LoginResult, error) {
	if resp == nil {
		return nil, errors.New("empty login response")
	}
	user, err := toUser(resp.GetUser())
	if err != nil {
		return nil, fmt.Errorf("toUser: %w", err)
	}
	session, err := toSession(resp.GetSession())
	if err != nil {
		return nil, fmt.Errorf("toSession: %w", err)
	}
	tokens, err := toTokenPair(resp.GetTokens())
	if err != nil {
		return nil, fmt.Errorf("toTokenPair: %w", err)
	}
	authCtx, err := toAuthContext(resp.GetAuthContext())
	if err != nil {
		return nil, fmt.Errorf("toAuthContext: %w", err)
	}
	return &models.LoginResult{
		User:    user,
		Session: session,
		Tokens:  tokens,
		AuthCtx: authCtx,
	}, nil
}

func toRefreshTokenResult(resp *authpb.RefreshTokenResponse) (*models.RefreshTokenResult, error) {
	if resp == nil {
		return nil, errors.New("empty refresh token response")
	}
	tokens, err := toTokenPair(resp.GetTokens())
	if err != nil {
		return nil, fmt.Errorf("toTokenPair: %w", err)
	}
	authCtx, err := toAuthContext(resp.GetAuthContext())
	if err != nil {
		return nil, fmt.Errorf("toAuthContext: %w", err)
	}
	return &models.RefreshTokenResult{
		Tokens:  tokens,
		AuthCtx: authCtx,
	}, nil
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

func toSession(s *authpb.Session) (models.Session, error) {
	if s == nil {
		return models.Session{}, errors.New("empty session")
	}
	id, err := uuid.Parse(s.GetId())
	if err != nil {
		return models.Session{}, fmt.Errorf("parse session id: %w", err)
	}
	userID, err := uuid.Parse(s.GetUserId())
	if err != nil {
		return models.Session{}, fmt.Errorf("parse session user_id: %w", err)
	}
	device, err := toDevice(s.GetDevice())
	if err != nil {
		return models.Session{}, fmt.Errorf("toDevice: %w", err)
	}
	return models.Session{
		ID:         id,
		UserID:     userID,
		Status:     toSessionStatus(s.GetStatus()),
		Device:     device,
		CreatedAt:  s.GetCreatedAt().AsTime(),
		UpdatedAt:  s.GetUpdatedAt().AsTime(),
		ExpiresAt:  s.GetExpiresAt().AsTime(),
		LastSeenAt: s.GetLastSeenAt().AsTime(),
	}, nil
}

func toDevice(d *authpb.Device) (models.Device, error) {
	if d == nil {
		return models.Device{}, errors.New("empty device")
	}
	id, err := uuid.Parse(d.GetId())
	if err != nil {
		return models.Device{}, fmt.Errorf("parse device id: %w", err)
	}
	return models.Device{
		ID:        id,
		Platform:  d.GetPlatform(),
		UserAgent: d.GetUserAgent(),
	}, nil
}

func toTokenPair(t *authpb.TokenPair) (models.TokenPair, error) {
	if t == nil {
		return models.TokenPair{}, errors.New("empty token pair")
	}
	return models.TokenPair{
		AccessToken:      t.GetAccessToken(),
		RefreshToken:     t.GetRefreshToken(),
		TokenType:        t.GetTokenType(),
		AccessExpiresAt:  t.GetAccessTokenExpiresAt().AsTime(),
		RefreshExpiresAt: t.GetRefreshTokenExpiresAt().AsTime(),
	}, nil
}

func toAuthContext(a *authpb.AuthContext) (models.AuthContext, error) {
	if a == nil {
		return models.AuthContext{}, errors.New("empty auth context")
	}
	userID, err := uuid.Parse(a.GetUserId())
	if err != nil {
		return models.AuthContext{}, fmt.Errorf("parse auth context user_id: %w", err)
	}
	sessionID, err := uuid.Parse(a.GetSessionId())
	if err != nil {
		return models.AuthContext{}, fmt.Errorf("parse auth context session_id: %w", err)
	}
	return models.AuthContext{
		UserID:          userID,
		SessionID:       sessionID,
		RoleCodes:       a.GetRoleCodes(),
		PermissionCodes: a.GetPermissionCodes(),
	}, nil
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

func toSessionStatus(s authpb.SessionStatus) models.SessionStatus {
	switch s { //nolint:exhaustive // unspecified handled by default
	case authpb.SessionStatus_SESSION_STATUS_ACTIVE:
		return models.SessionStatusActive
	case authpb.SessionStatus_SESSION_STATUS_REVOKED:
		return models.SessionStatusRevoked
	case authpb.SessionStatus_SESSION_STATUS_EXPIRED:
		return models.SessionStatusExpired
	default:
		return models.SessionStatusActive
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
	case codes.Unauthenticated:
		return models.NewUnauthorizedError(st.Message(), err)
	case codes.PermissionDenied:
		return models.NewBusinessError(st.Message(), err)
	default:
		return fmt.Errorf("auth service error: %w", err)
	}
}
