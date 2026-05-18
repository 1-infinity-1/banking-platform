package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "active"
	SessionStatusRevoked SessionStatus = "revoked"
	SessionStatusExpired SessionStatus = "expired"
)

type Device struct {
	ID        uuid.UUID
	Platform  string
	UserAgent string
}

type Session struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Status     SessionStatus
	Device     Device
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ExpiresAt  time.Time
	LastSeenAt time.Time
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	TokenType        string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

type AuthContext struct {
	UserID          uuid.UUID
	SessionID       uuid.UUID
	RoleCodes       []string
	PermissionCodes []string
}

type LoginParams struct {
	Login     string
	Password  string
	UserAgent string
	Platform  string
}

type LoginResult struct {
	User    *User
	Session Session
	Tokens  TokenPair
	AuthCtx AuthContext
}

type LogoutParams struct {
	RefreshToken string
	UserAgent    string
	Platform     string
}

type RefreshTokenParams struct {
	RefreshToken string
	UserAgent    string
	Platform     string
}

type RefreshTokenResult struct {
	Tokens  TokenPair
	AuthCtx AuthContext
}
