package management

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
)

type txManager interface {
	BeginFunc(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type roleRepo interface {
	GetRolesByCodesTx(ctx context.Context, tx pgx.Tx, roleCode []string) ([]models.Role, error)
	CreateUserRolesTx(ctx context.Context, tx pgx.Tx, userID int64, rolesIDs []int64) error
}

type userRepo interface {
	CreateUserTx(ctx context.Context, tx pgx.Tx, user models.CreateUser, passwordHashed string, status models.UserStatus) (*models.User, error)
}

type AccessManagementService struct {
	txManager txManager
	userRepo  userRepo
	roleRepo  roleRepo
}

func NewAccessManagementService(txManager txManager, userRepo userRepo, roleRepo roleRepo) *AccessManagementService {
	return &AccessManagementService{
		txManager: txManager,
		userRepo:  userRepo,
		roleRepo:  roleRepo,
	}
}

func (u *AccessManagementService) CreateUser(ctx context.Context, userCreate models.CreateUser, password string) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	passwordHashed := string(hashed)

	statusUser := models.UserStatusActive

	var user *models.User
	err = u.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		user, err = u.userRepo.CreateUserTx(ctx, tx, userCreate, passwordHashed, statusUser)
		if err != nil {
			return fmt.Errorf("u.repo.Create: %w", err)
		}

		roles, err := u.roleRepo.GetRolesByCodesTx(ctx, tx, userCreate.Role)
		if err != nil {
			return fmt.Errorf("u.repo.GetRoles: %w", err)
		}

		if len(roles) == 0 {
			return models.NewNotFoundError("no roles found")
		}

		var roleIDs = make([]int64, 0, len(roles))
		var roleCodes = make([]string, 0, len(roles))
		for _, role := range roles {
			roleIDs = append(roleIDs, role.ID)
			roleCodes = append(roleCodes, role.Code)
		}

		user.Roles = roleCodes

		err = u.roleRepo.CreateUserRolesTx(ctx, tx, user.ID, roleIDs)
		if err != nil {
			return fmt.Errorf("u.repo.CreateUserRoles: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("u.txManager.BeginFunc: %w", err)
	}

	return user, nil
}
