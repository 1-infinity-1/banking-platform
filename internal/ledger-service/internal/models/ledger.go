package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type EntryType string

const (
	EntryTypeCredit EntryType = "credit"
	EntryTypeDebit  EntryType = "debit"
)

type LedgerEntry struct {
	ID            int64
	PublicID      uuid.UUID
	TransactionID uuid.UUID
	AccountID     uuid.UUID
	Type          EntryType
	Amount        decimal.Decimal
	Currency      string
	BalanceAfter  decimal.Decimal
	Description   string
	OccurredAt    time.Time
	CreatedAt     time.Time
}

type Statement struct {
	AccountID uuid.UUID
	Entries   []LedgerEntry
}

// TransactionCompletedEvent is the JSON structure of a Kafka message
// produced by transaction-service on topic transactions.completed.
type TransactionCompletedEvent struct {
	TransactionID string    `json:"transaction_id"`
	AccountID     string    `json:"account_id"`
	Type          string    `json:"type"`
	Amount        string    `json:"amount"`
	Currency      string    `json:"currency"`
	BalanceAfter  string    `json:"balance_after"`
	Description   string    `json:"description"`
	OccurredAt    time.Time `json:"occurred_at"`
}
