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

	status, err := models.ToStatus(u.status)
	if err != nil {
		return nil, fmt.Errorf("models.ToStatus: %w", err)
	}

	user.Status = status

	return user, nil
}
