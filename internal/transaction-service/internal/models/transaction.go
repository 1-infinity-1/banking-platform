package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionStatus string

const (
	TransactionStatusUnspecified TransactionStatus = "unspecified"
	TransactionStatusPending     TransactionStatus = "pending"
	TransactionStatusCompleted   TransactionStatus = "completed"
	TransactionStatusFailed      TransactionStatus = "failed"
	TransactionStatusCancelled   TransactionStatus = "cancelled"
)

func ToTransactionStatus(s string) (TransactionStatus, error) {
	switch TransactionStatus(s) {
	case TransactionStatusPending:
		return TransactionStatusPending, nil
	case TransactionStatusCompleted:
		return TransactionStatusCompleted, nil
	case TransactionStatusFailed:
		return TransactionStatusFailed, nil
	case TransactionStatusCancelled:
		return TransactionStatusCancelled, nil
	case TransactionStatusUnspecified:
		return TransactionStatusUnspecified, nil
	}
	return "", fmt.Errorf("unknown transaction status: %s", s)
}

type Transaction struct {
	ID             int64
	PublicID       uuid.UUID
	FromAccountID  *uuid.UUID // nil for replenishments
	ToAccountID    uuid.UUID
	Amount         decimal.Decimal
	Currency       string
	Status         TransactionStatus
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TransactionEvent is the Kafka payload published on transaction.completed.
type TransactionEvent struct {
	TransactionID string            `json:"transaction_id"`
	FromAccountID *string           `json:"from_account_id,omitempty"`
	ToAccountID   string            `json:"to_account_id"`
	Amount        string            `json:"amount"`
	Currency      string            `json:"currency"`
	Status        TransactionStatus `json:"status"`
	OccurredAt    time.Time         `json:"occurred_at"`
}

type TransferRequest struct {
	FromAccountID  uuid.UUID
	ToAccountID    uuid.UUID
	Amount         decimal.Decimal
	Currency       string
	IdempotencyKey string
}

type ReplenishRequest struct {
	ToAccountID    uuid.UUID
	Amount         decimal.Decimal
	Currency       string
	IdempotencyKey string
}

type GetHistoryRequest struct {
	AccountID uuid.UUID
	Limit     int32
	Offset    int32
}

// DebitRequest and CreditRequest are used by both the service layer and the account gRPC client.
type DebitRequest struct {
	AccountID      uuid.UUID
	Amount         decimal.Decimal
	IdempotencyKey string
}

type DebitResult struct {
	AccountID    string
	BalanceAfter decimal.Decimal
}

type CreditRequest struct {
	AccountID      uuid.UUID
	Amount         decimal.Decimal
	IdempotencyKey string
}

type CreditResult struct {
	AccountID    string
	BalanceAfter decimal.Decimal
}
