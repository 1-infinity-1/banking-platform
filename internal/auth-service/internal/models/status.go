package models

import "fmt"

type UserStatus string

const (
	UserStatusUnspecified UserStatus = ""

	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusBlocked  UserStatus = "BLOCKED"
	UserStatusLocked   UserStatus = "LOCKED"
	UserStatusDisabled UserStatus = "DISABLED"
)

func ToUserStatus(st string) (UserStatus, error) {
	switch st {
	case string(UserStatusActive):
		return UserStatusActive, nil
	case string(UserStatusBlocked):
		return UserStatusBlocked, nil
	case string(UserStatusLocked):
		return UserStatusLocked, nil
	case string(UserStatusDisabled):
		return UserStatusDisabled, nil
	default:
		return UserStatusUnspecified, fmt.Errorf("invalid status: %s", st)
	}
}

type SessionStatus string

const (
	SessionStatusUnspecified SessionStatus = ""

	SessionStatusActive  SessionStatus = "ACTIVE"
	SessionStatusRevoked SessionStatus = "REVOKED"
	SessionStatusExpired SessionStatus = "EXPIRED"
)

func ToSessionStatus(st string) (SessionStatus, error) {
	switch st {
	case string(SessionStatusActive):
		return SessionStatusActive, nil
	case string(SessionStatusRevoked):
		return SessionStatusRevoked, nil
	case string(SessionStatusExpired):
		return SessionStatusExpired, nil
	default:
		return SessionStatusUnspecified, fmt.Errorf("invalid status: %s", st)
	}
}
