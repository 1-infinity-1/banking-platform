package auth

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

type authClient interface {
	Login(ctx context.Context, params models.LoginParams) (*models.LoginResult, error)
	Logout(ctx context.Context, params models.LogoutParams) error
	RefreshToken(ctx context.Context, params models.RefreshTokenParams) (*models.RefreshTokenResult, error)
}

type Service struct {
	authClient authClient
}

func New(authClient authClient) *Service {
	return &Service{authClient: authClient}
}

func (s *Service) Login(ctx context.Context, params models.LoginParams) (*models.LoginResult, error) {
	result, err := s.authClient.Login(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.authClient.Login: %w", err)
	}
	return result, nil
}

func (s *Service) Logout(ctx context.Context, params models.LogoutParams) error {
	if err := s.authClient.Logout(ctx, params); err != nil {
		return fmt.Errorf("s.authClient.Logout: %w", err)
	}
	return nil
}

func (s *Service) RefreshToken(
	ctx context.Context,
	params models.RefreshTokenParams,
) (*models.RefreshTokenResult, error) {
	result, err := s.authClient.RefreshToken(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.authClient.RefreshToken: %w", err)
	}
	return result, nil
}
