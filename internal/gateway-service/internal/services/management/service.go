package management

import (
	"context"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

type authClient interface {
	CreateUser(ctx context.Context, params models.CreateUserParams) (*models.User, error)
}

type AccessManagementService struct {
	authClient authClient
}

func New(authClient authClient) *AccessManagementService {
	return &AccessManagementService{authClient: authClient}
}

func (s *AccessManagementService) CreateUser(
	ctx context.Context,
	params models.CreateUserParams,
) (*models.User, error) {
	return s.authClient.CreateUser(ctx, params)
}
