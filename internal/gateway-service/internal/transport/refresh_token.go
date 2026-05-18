package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) RefreshToken(ctx context.Context, req *api.RefreshTokenRequest) (api.RefreshTokenRes, error) {
	if req.RefreshToken == "" {
		return nil, models.NewValidationError("refresh_token", "is required", nil)
	}

	params := models.RefreshTokenParams{RefreshToken: req.RefreshToken}
	if rc, ok := req.Context.Get(); ok {
		params.UserAgent = rc.UserAgent.Or("")
		params.Platform = rc.Platform.Or("")
	}

	result, err := g.authSvc.RefreshToken(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("g.authSvc.RefreshToken: %w", err)
	}

	return &api.RefreshTokenResponse{
		Tokens:      toAPITokenPair(result.Tokens),
		AuthContext: toAPIAuthContext(result.AuthCtx),
	}, nil
}
