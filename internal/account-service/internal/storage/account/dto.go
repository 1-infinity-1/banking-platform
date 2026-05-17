package account

import (
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//nolint:unused // scaffold: used when implementing TODO methods
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

//nolint:unused // scaffold: used when implementing TODO methods
func (d *accountDTO) ToDomain() (*models.Account, error) {
	// TODO: implement — map fields + call models.ToAccountStatus(d.status)
	return nil, fmt.Errorf("ToDomain: %w", models.ErrInternal)
}
