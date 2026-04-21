package role

import (
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/google/uuid"
)

type Role struct {
	id       int64
	publicID uuid.UUID
	code     string
	name     string
}

func (r *Role) ToDomain() models.Role {
	return models.Role{
		ID:       r.id,
		PublicID: r.publicID,
		Code:     r.code,
		Name:     r.name,
	}
}
