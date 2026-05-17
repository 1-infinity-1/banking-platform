package entry

import (
	"fmt"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type entryDTO struct {
	id            int64
	publicID      uuid.UUID
	transactionID uuid.UUID
	accountID     uuid.UUID
	entryType     string
	amount        string
	currency      string
	balanceAfter  string
	description   *string
	occurredAt    time.Time
	createdAt     time.Time
}

func (d *entryDTO) ToDomain() (models.LedgerEntry, error) {
	amount, err := decimal.NewFromString(d.amount)
	if err != nil {
		return models.LedgerEntry{}, fmt.Errorf("decimal.NewFromString amount: %w", err)
	}

	balanceAfter, err := decimal.NewFromString(d.balanceAfter)
	if err != nil {
		return models.LedgerEntry{}, fmt.Errorf("decimal.NewFromString balanceAfter: %w", err)
	}

	desc := ""
	if d.description != nil {
		desc = *d.description
	}

	return models.LedgerEntry{
		ID:            d.id,
		PublicID:      d.publicID,
		TransactionID: d.transactionID,
		AccountID:     d.accountID,
		Type:          models.EntryType(d.entryType),
		Amount:        amount,
		Currency:      d.currency,
		BalanceAfter:  balanceAfter,
		Description:   desc,
		OccurredAt:    d.occurredAt,
		CreatedAt:     d.createdAt,
	}, nil
}
