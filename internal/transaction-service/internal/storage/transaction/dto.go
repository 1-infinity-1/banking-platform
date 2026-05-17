package transaction

import (
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//nolint:unused // scaffold: used when implementing TODO methods
type transactionDTO struct {
	id             int64
	publicID       uuid.UUID
	fromAccountID  *uuid.UUID
	toAccountID    uuid.UUID
	amount         decimal.Decimal
	currency       string
	status         string
	idempotencyKey string
	createdAt      time.Time
	updatedAt      time.Time
}

//nolint:unused // scaffold: used when implementing TODO methods
func (d *transactionDTO) ToDomain() (*models.Transaction, error) {
	// TODO: implement — map fields + call models.ToTransactionStatus(d.status)
	return nil, fmt.Errorf("ToDomain: %w", models.ErrInternal)
}
