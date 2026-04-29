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
