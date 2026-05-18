package middleware

import (
	"context"
	"errors"
	"fmt"
	"slices"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authClaimsKey struct{}

// AuthClaims holds the validated JWT payload injected into the request context.
type AuthClaims struct {
	UserID          uuid.UUID
	SessionID       uuid.UUID
	RoleCodes       []string
	PermissionCodes []string
}

// FromContext returns the AuthClaims stored by JWTSecurityHandler.
func FromContext(ctx context.Context) (AuthClaims, bool) {
	v, ok := ctx.Value(authClaimsKey{}).(AuthClaims)
	return v, ok
}

func newContextWithClaims(ctx context.Context, claims AuthClaims) context.Context {
	return context.WithValue(ctx, authClaimsKey{}, claims)
}

// JWTSecurityHandler implements api.SecurityHandler using HS256-signed JWTs.
// operationRoles maps an operation name to the set of roles of which at least
// one must be present in the token. Operations absent from the map allow any
// authenticated user through.
type JWTSecurityHandler struct {
	secret         []byte
	operationRoles map[api.OperationName][]string
}

// NewJWTSecurityHandler creates a JWTSecurityHandler.
// operationRoles may be nil (no role restrictions beyond authentication).
func NewJWTSecurityHandler(secret string, operationRoles map[api.OperationName][]string) *JWTSecurityHandler {
	return &JWTSecurityHandler{
		secret:         []byte(secret),
		operationRoles: operationRoles,
	}
}

// HandleBearerAuth is called by ogen for every protected operation.
// It validates the JWT, checks required roles, and injects claims into the context.
func (h *JWTSecurityHandler) HandleBearerAuth(
	ctx context.Context,
	op api.OperationName,
	t api.BearerAuth,
) (context.Context, error) {
	mapClaims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(t.Token, mapClaims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.secret, nil
	})
	if err != nil || !token.Valid {
		return ctx, fmt.Errorf("invalid token: %w", err)
	}

	claims, err := extractClaims(mapClaims)
	if err != nil {
		return ctx, fmt.Errorf("extract claims: %w", err)
	}

	ctx = newContextWithClaims(ctx, claims)

	if required, ok := h.operationRoles[op]; ok {
		if !hasAnyRole(claims.RoleCodes, required) {
			return ctx, models.NewForbiddenError("insufficient role", nil)
		}
	}

	return ctx, nil
}

// hasAnyRole reports whether roles contains at least one value from required.
func hasAnyRole(roles, required []string) bool {
	for _, r := range required {
		if slices.Contains(roles, r) {
			return true
		}
	}
	return false
}

func extractClaims(m jwt.MapClaims) (AuthClaims, error) {
	userIDStr, ok := m["user_id"].(string)
	if !ok || userIDStr == "" {
		return AuthClaims{}, errors.New("missing user_id claim")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return AuthClaims{}, fmt.Errorf("parse user_id: %w", err)
	}

	sessionIDStr, ok := m["session_id"].(string)
	if !ok || sessionIDStr == "" {
		return AuthClaims{}, errors.New("missing session_id claim")
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return AuthClaims{}, fmt.Errorf("parse session_id: %w", err)
	}

	roles := toStringSlice(m["roles"])
	permissions := toStringSlice(m["permissions"])

	return AuthClaims{
		UserID:          userID,
		SessionID:       sessionID,
		RoleCodes:       roles,
		PermissionCodes: permissions,
	}, nil
}

func toStringSlice(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, strOk := item.(string); strOk {
			result = append(result, s)
		}
	}
	return result
}
