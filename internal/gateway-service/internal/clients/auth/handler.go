package auth

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (c *Client) CreateUser(ctx context.Context, params models.CreateUserParams) (*models.User, error) {
	resp, err := c.accessManagement.CreateUser(ctx, toProtoCreateUserRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	user, err := toUser(resp.GetUser())
	if err != nil {
		return nil, fmt.Errorf("toUser: %w", err)
	}
	return user, nil
}

func (c *Client) Login(ctx context.Context, params models.LoginParams) (*models.LoginResult, error) {
	resp, err := c.authSvc.Login(ctx, toProtoLoginRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	result, err := toLoginResult(resp)
	if err != nil {
		return nil, fmt.Errorf("toLoginResult: %w", err)
	}
	return result, nil
}

func (c *Client) Logout(ctx context.Context, params models.LogoutParams) error {
	_, err := c.authSvc.Logout(ctx, toProtoLogoutRequest(params))
	if err != nil {
		return mapGRPCError(err)
	}
	return nil
}

func (c *Client) RefreshToken(
	ctx context.Context,
	params models.RefreshTokenParams,
) (*models.RefreshTokenResult, error) {
	resp, err := c.authSvc.RefreshToken(ctx, toProtoRefreshTokenRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	result, err := toRefreshTokenResult(resp)
	if err != nil {
		return nil, fmt.Errorf("toRefreshTokenResult: %w", err)
	}
	return result, nil
}
