package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) Logout(ctx context.Context, req *api.LogoutRequest) (api.LogoutRes, error) {
	if req.RefreshToken == "" {
		return nil, models.NewValidationError("refresh_token", "is required", nil)
	}

	params := models.LogoutParams{RefreshToken: req.RefreshToken}
	if rc, ok := req.Context.Get(); ok {
		params.UserAgent = rc.UserAgent.Or("")
		params.Platform = rc.Platform.Or("")
	}

	if err := g.authSvc.Logout(ctx, params); err != nil {
		return nil, fmt.Errorf("g.authSvc.Logout: %w", err)
	}

	return &api.LogoutNoContent{}, nil
}
