package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/1-infinity-1/banking-platform/pkg/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
)

type Repository struct {
	db *postgres.Conn
}

func NewRepository(db *postgres.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUserTx(ctx context.Context, tx pgx.Tx, user models.CreateUser, status models.Status) (*models.User, error) {
	query := `
		INSERT INTO users (
			login,
			email,
		    phone,
			password_hash,
		    status
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
		    id, 
		    public_id, 
		    login, 
		    email, 
		    phone, 
		    status, 
		    created_at, 
		    updated_at
	`

	var userDTO CreateUserDTO
	err := tx.QueryRow(
		ctx,
		query,
		user.Login,
		user.Email,
		user.Phone,
		user.Password,
		status,
	).Scan(
		&userDTO.id,
		&userDTO.publicID,
		&userDTO.login,
		&userDTO.email,
		&userDTO.phone,
		&userDTO.status,
		&userDTO.createdAt,
		&userDTO.updatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, models.NewInvalidParamsError("unique", "this data is already in use")
		}
		return nil, fmt.Errorf("r.db.Exec %w", err)
	}

	userModel, err := userDTO.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("userDTO.ToDomain: %w", err)
	}

	return userModel, nil
}
