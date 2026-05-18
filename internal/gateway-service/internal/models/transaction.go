package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusUnspecified TransactionStatus = "unspecified"
	TransactionStatusPending     TransactionStatus = "pending"
	TransactionStatusCompleted   TransactionStatus = "completed"
	TransactionStatusFailed      TransactionStatus = "failed"
	TransactionStatusCancelled   TransactionStatus = "cancelled"
)

type Transaction struct {
	ID             uuid.UUID
	FromAccountID  *uuid.UUID
	ToAccountID    uuid.UUID
	Amount         string
	Currency       string
	Status         TransactionStatus
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TransferParams struct {
	FromAccountID  uuid.UUID
	ToAccountID    uuid.UUID
	Amount         string
	Currency       string
	IdempotencyKey string
}

type ReplenishParams struct {
	ToAccountID    uuid.UUID
	Amount         string
	Currency       string
	IdempotencyKey string
}

type GetHistoryParams struct {
	AccountID uuid.UUID
	Limit     int32
	Offset    int32
}

type GetTransactionParams struct {
	TransactionID uuid.UUID
}
