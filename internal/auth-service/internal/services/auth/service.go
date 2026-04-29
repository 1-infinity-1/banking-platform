package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type txManager interface {
	BeginFunc(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type userRepo interface {
	GetByLogin(ctx context.Context, login string) (*models.User, error)
}

type deviceRepo interface {
	CreateDeviceTx(ctx context.Context, tx pgx.Tx, userAgent, platform string) (*models.Device, error)
}

type sessionRepo interface {
	CreateSessionTx(ctx context.Context, tx pgx.Tx, userID, deviceID int64, status models.SessionStatus, expireTime time.Time) (*models.Session, error)
}

type refreshTokenRepo interface {
	CreateTokenTx(ctx context.Context, tx pgx.Tx, sessionID int64, tokenHash string, expireTime time.Time) error
}

type tokenManager interface {
	GenerateAccessToken(user models.User, session models.Session, expireAt time.Time) (string, error)
	GenerateRefreshToken() (string, error)
}

type Config struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type AuthService struct {
	txManager        txManager
	userRepo         userRepo
	deviceRepo       deviceRepo
	sessionRepo      sessionRepo
	tokenManager     tokenManager
	refreshTokenRepo refreshTokenRepo
	cfg              Config
}

func NewAuthService(txManager txManager, userRepo userRepo, deviceRepo deviceRepo, sessionRepo sessionRepo, tokenManager tokenManager, refreshTokenRepo refreshTokenRepo, cfg Config) *AuthService {
	return &AuthService{
		txManager:        txManager,
		userRepo:         userRepo,
		deviceRepo:       deviceRepo,
		sessionRepo:      sessionRepo,
		tokenManager:     tokenManager,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
	}
}

func (s *AuthService) Login(ctx context.Context, credentials models.LoginCredentials) (*models.LoginResult, error) {
	user, err := s.userRepo.GetByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.GetByLogin: %w", err)
	}

	isPasswordsMatched := checkPassword(user.PasswordHash, credentials.Password)
	if !isPasswordsMatched {
		return nil, models.NewInvalidParamsError("password", "incorrect password")
	}

	if user.Status != models.UserStatusActive {
		return nil, models.NewBusinessError("account is not available for login")
	}

	var loginResult *models.LoginResult
	err = s.txManager.BeginFunc(ctx, func(tx pgx.Tx) error {
		device, err := s.deviceRepo.CreateDeviceTx(ctx, tx, credentials.Context.UserAgent, credentials.Context.Platform)
		if err != nil {
			return fmt.Errorf("s.deviceRepo.CreateDeviceTx: %w", err)
		}

		expireAt := time.Now().Add(s.cfg.RefreshTokenTTL)
		session, err := s.sessionRepo.CreateSessionTx(ctx, tx, user.ID, device.ID, models.SessionStatusActive, expireAt)
		if err != nil {
			return fmt.Errorf("s.sessionRepo.CreateSessionTx: %w", err)
		}

		expireAtAccessToken := time.Now().Add(s.cfg.AccessTokenTTL)
		accessToken, err := s.tokenManager.GenerateAccessToken(*user, *session, expireAtAccessToken)
		if err != nil {
			return fmt.Errorf("s.tokenManager.GenerateAccessToken: %w", err)
		}

		refreshToken, err := s.tokenManager.GenerateRefreshToken()
		if err != nil {
			return fmt.Errorf("s.tokenManager.GenerateRefreshToken: %w", err)
		}

		hash := sha256.Sum256([]byte(refreshToken))
		tokenHash := hex.EncodeToString(hash[:])
		err = s.refreshTokenRepo.CreateTokenTx(ctx, tx, session.ID, tokenHash, expireAt)
		if err != nil {
			return fmt.Errorf("s.refreshTokenRepo.CreateTokenTx: %w", err)
		}

		loginResult = mappingToLoginResult(user, session, device, accessToken, refreshToken, expireAtAccessToken, expireAt)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.txManager.BeginFunc: %w", err)
	}

	return loginResult, nil
}

func mappingToLoginResult(
	user *models.User,
	session *models.Session,
	device *models.Device,
	accessToken string,
	refreshToken string,
	expireAtAccessToken time.Time,
	expireAtRefreshToken time.Time,
) *models.LoginResult {
	return &models.LoginResult{
		User:    *user,
		Session: *session,
		Device:  *device,
		Tokens: models.TokenPair{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			TypeToken:             models.TokenTypeBearer,
			AccessTokenExpiresAt:  expireAtAccessToken,
			RefreshTokenExpiresAt: expireAtRefreshToken,
		},
	}
}
