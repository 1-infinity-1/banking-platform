package transport

import (
	"context"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) CreateUser(ctx context.Context, req *api.CreateUserRequest) (api.CreateUserRes, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	params := models.CreateUserParams{
		Login:     req.Login,
		Password:  req.Password,
		RoleCodes: req.RoleCodes,
	}
	if v, ok := req.Email.Get(); ok {
		params.Email = &v
	}
	if v, ok := req.Phone.Get(); ok {
		params.Phone = &v
	}

	user, err := g.managementSvc.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return toCreateUserResponse(user), nil
}

func validateCreateUserRequest(req *api.CreateUserRequest) error {
	if req.Login == "" {
		return models.NewValidationError("login", "is required", nil)
	}
	if req.Password == "" {
		return models.NewValidationError("password", "is required", nil)
	}
	if len(req.RoleCodes) == 0 {
		return models.NewValidationError("role_codes", "is required", nil)
	}
	_, hasEmail := req.Email.Get()
	_, hasPhone := req.Phone.Get()
	if !hasEmail && !hasPhone {
		return models.NewValidationError("", "email or phone is required", nil)
	}
	return nil
}

func toCreateUserResponse(u *models.User) *api.CreateUserResponse {
	apiUser := api.User{
		ID:        u.ID,
		Login:     u.Login,
		Status:    api.UserStatus(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.Email != nil {
		apiUser.Email = api.NewOptNilString(*u.Email)
	}
	if u.Phone != nil {
		apiUser.Phone = api.NewOptNilString(*u.Phone)
	}
	return &api.CreateUserResponse{User: apiUser}
}
