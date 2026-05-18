package transaction

import (
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type transactionDTO struct {
	id             int64
	publicID       uuid.UUID
	fromAccountID  pgtype.UUID
	toAccountID    uuid.UUID
	amount         decimal.Decimal
	currency       string
	status         string
	idempotencyKey string
	createdAt      time.Time
	updatedAt      time.Time
}

func (d *transactionDTO) ToDomain() (*models.Transaction, error) {
	status, err := models.ToTransactionStatus(d.status)
	if err != nil {
		return nil, fmt.Errorf("models.ToTransactionStatus: %w", err)
	}

	var fromAccountID *uuid.UUID
	if d.fromAccountID.Valid {
		from := uuid.UUID(d.fromAccountID.Bytes)
		fromAccountID = &from
	}

	return &models.Transaction{
		ID:             d.id,
		PublicID:       d.publicID,
		FromAccountID:  fromAccountID,
		ToAccountID:    d.toAccountID,
		Amount:         d.amount,
		Currency:       d.currency,
		Status:         status,
		IdempotencyKey: d.idempotencyKey,
		CreatedAt:      d.createdAt,
		UpdatedAt:      d.updatedAt,
	}, nil
}
