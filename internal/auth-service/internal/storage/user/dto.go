package user

import (
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/google/uuid"
)

type CreateUserDTO struct {
	id        int64
	publicID  uuid.UUID
	login     string
	email     *string
	phone     *string
	status    string
	createdAt time.Time
	updatedAt time.Time
}

func (u *CreateUserDTO) ToDomain() (*models.User, error) {
	user := &models.User{
		ID:        u.id,
		PublicID:  u.publicID,
		Login:     u.login,
		Email:     u.email,
		Phone:     u.phone,
		CreatedAt: u.createdAt,
		UpdatedAt: u.updatedAt,
	}

	status, err := models.ToUserStatus(u.status)
	if err != nil {
		return nil, fmt.Errorf("models.ToStatus: %w", err)
	}
	user.Status = status

	return user, nil
}

type UserDTO struct {
	id           int64
	publicID     uuid.UUID
	login        string
	email        *string
	phone        *string
	passwordHash string
	status       string
	roles        []string
	permissions  []string
	createdAt    time.Time
	updatedAt    time.Time
}

func (u *UserDTO) ToDomain() (*models.User, error) {
	user := &models.User{
		ID:           u.id,
		PublicID:     u.publicID,
		Login:        u.login,
		Email:        u.email,
		Phone:        u.phone,
		PasswordHash: u.passwordHash,
		Roles:        u.roles,
		Permissions:  u.permissions,
		CreatedAt:    u.createdAt,
		UpdatedAt:    u.updatedAt,
	}

	status, err := models.ToUserStatus(u.status)
	if err != nil {
		return nil, fmt.Errorf("models.ToStatus: %w", err)
	}
	user.Status = status

	return user, nil
}
