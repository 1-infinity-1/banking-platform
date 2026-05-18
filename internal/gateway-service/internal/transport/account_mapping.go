package transport

import (
	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func toAPIAccount(a models.Account) api.Account {
	return api.Account{
		ID:        a.ID,
		UserID:    a.UserID,
		Currency:  a.Currency,
		Balance:   a.Balance,
		Status:    api.AccountStatus(a.Status),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func toAPIBalance(b models.Balance) api.BalanceResponse {
	return api.BalanceResponse{
		AccountID: b.AccountID,
		Amount:    b.Amount,
		Currency:  b.Currency,
	}
}
