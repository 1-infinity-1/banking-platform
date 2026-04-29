package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secretKey string
}

func NewTokenManager(secretKey string) *TokenManager {
	return &TokenManager{
		secretKey: secretKey,
	}
}

func (m *TokenManager) GenerateAccessToken(user models.User, session models.Session, expireAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     user.PublicID,
		"session_id":  session.PublicID.String(),
		"roles":       user.Roles,
		"permissions": user.Permissions,
		"exp":         expireAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *TokenManager) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
