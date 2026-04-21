package role

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Conn
}

func NewRepository(db *postgres.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRolesByCodesTx(ctx context.Context, tx pgx.Tx, roleCode []string) ([]models.Role, error) {
	query := `
		SELECT 
		    id,
		    public_id,
			code,
			name
		FROM roles
		WHERE code = ANY($1)
	`

	rows, err := tx.Query(ctx, query, roleCode)
	if err != nil {
		return nil, fmt.Errorf("r.db.Query: %w", err)
	}
	defer rows.Close()

	rolesDTO := make([]Role, 0)
	for rows.Next() {
		var roleDTO Role
		if err = rows.Scan(&roleDTO.id, &roleDTO.publicID, &roleDTO.code, &roleDTO.name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		rolesDTO = append(rolesDTO, roleDTO)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	roles := make([]models.Role, 0, len(rolesDTO))
	for _, roleDTO := range rolesDTO {
		role := roleDTO.ToDomain()
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *Repository) CreateUserRolesTx(ctx context.Context, tx pgx.Tx, userID int64, rolesIDs []int64) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, unnest($2::bigint[])
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := tx.Exec(ctx, query, userID, rolesIDs)
	if err != nil {
		return fmt.Errorf("r.db.Exec: %w", err)
	}

	return nil
}
