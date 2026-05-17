package models

import "time"

const (
	TokenTypeBearer = "Bearer"
)

type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	TypeToken             string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type RefreshToken struct {
	ID        int64
	SessionID int64
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
