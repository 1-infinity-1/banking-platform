package transport

import (
	"context"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

type managementService interface {
	CreateUser(ctx context.Context, params models.CreateUserParams) (*models.User, error)
}

type GatewayHandler struct {
	api.UnimplementedHandler

	managementSvc managementService
}

func NewGatewayHandler(managementSvc managementService) *GatewayHandler {
	return &GatewayHandler{managementSvc: managementSvc}
}
