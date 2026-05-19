package account

import (
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type accountDTO struct {
	id        int64
	publicID  uuid.UUID
	userID    uuid.UUID
	currency  string
	balance   decimal.Decimal
	status    string
	createdAt time.Time
	updatedAt time.Time
}

func (d *accountDTO) ToDomain() (*models.Account, error) {
	status, err := models.ToAccountStatus(d.status)
	if err != nil {
		return nil, fmt.Errorf("models.ToAccountStatus: %w", err)
	}
	return &models.Account{
		ID:        d.id,
		PublicID:  d.publicID,
		UserID:    d.userID,
		Currency:  d.currency,
		Balance:   d.balance,
		Status:    status,
		CreatedAt: d.createdAt,
		UpdatedAt: d.updatedAt,
	}, nil
}
