package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) Login(ctx context.Context, req *api.LoginRequest) (api.LoginRes, error) {
	if req.Login == "" {
		return nil, models.NewValidationError("login", "is required", nil)
	}
	if req.Password == "" {
		return nil, models.NewValidationError("password", "is required", nil)
	}

	params := models.LoginParams{
		Login:    req.Login,
		Password: req.Password,
	}
	if rc, ok := req.Context.Get(); ok {
		params.UserAgent = rc.UserAgent.Or("")
		params.Platform = rc.Platform.Or("")
	}

	result, err := g.authSvc.Login(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("g.authSvc.Login: %w", err)
	}

	return &api.LoginResponse{
		User:        toAPIUser(result.User),
		Session:     toAPISession(result.Session),
		Tokens:      toAPITokenPair(result.Tokens),
		AuthContext: toAPIAuthContext(result.AuthCtx),
	}, nil
}
