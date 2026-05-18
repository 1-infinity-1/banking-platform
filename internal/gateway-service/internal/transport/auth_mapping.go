package transport

import (
	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func toAPIUser(u *models.User) api.User {
	apiUser := api.User{
		ID:        u.ID,
		Login:     u.Login,
		Status:    api.UserStatus(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.Email != nil {
		apiUser.Email = api.NewOptNilString(*u.Email)
	}
	if u.Phone != nil {
		apiUser.Phone = api.NewOptNilString(*u.Phone)
	}
	return apiUser
}

func toAPISession(s models.Session) api.Session {
	return api.Session{
		ID:         s.ID,
		UserID:     s.UserID,
		Status:     api.SessionStatus(s.Status),
		Device:     toAPIDevice(s.Device),
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		ExpiresAt:  s.ExpiresAt,
		LastSeenAt: s.LastSeenAt,
	}
}

func toAPIDevice(d models.Device) api.Device {
	return api.Device{
		ID:        d.ID,
		Platform:  d.Platform,
		UserAgent: d.UserAgent,
	}
}

func toAPITokenPair(t models.TokenPair) api.TokenPair {
	return api.TokenPair{
		AccessToken:           t.AccessToken,
		RefreshToken:          t.RefreshToken,
		TokenType:             t.TokenType,
		AccessTokenExpiresAt:  t.AccessExpiresAt,
		RefreshTokenExpiresAt: t.RefreshExpiresAt,
	}
}

func toAPIAuthContext(a models.AuthContext) api.AuthContext {
	return api.AuthContext{
		UserID:          a.UserID,
		SessionID:       a.SessionID,
		RoleCodes:       a.RoleCodes,
		PermissionCodes: a.PermissionCodes,
	}
}
