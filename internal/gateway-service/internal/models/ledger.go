package models

import (
	"time"

	"github.com/google/uuid"
)

type LedgerEntry struct {
	ID            uuid.UUID
	TransactionID uuid.UUID
	AccountID     uuid.UUID
	Type          string
	Amount        string
	Currency      string
	BalanceAfter  string
	Description   string
	OccurredAt    time.Time
	CreatedAt     time.Time
}

type Statement struct {
	AccountID uuid.UUID
	Entries   []LedgerEntry
}

type GetStatementParams struct {
	AccountID uuid.UUID
	From      time.Time
	To        time.Time
}
