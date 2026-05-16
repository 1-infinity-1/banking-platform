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
