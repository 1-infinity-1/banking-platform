package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountStatus string

const (
	AccountStatusUnspecified AccountStatus = "unspecified"
	AccountStatusActive      AccountStatus = "active"
	AccountStatusBlocked     AccountStatus = "blocked"
	AccountStatusClosed      AccountStatus = "closed"
)

func ToAccountStatus(s string) (AccountStatus, error) {
	switch AccountStatus(s) {
	case AccountStatusActive:
		return AccountStatusActive, nil
	case AccountStatusBlocked:
		return AccountStatusBlocked, nil
	case AccountStatusClosed:
		return AccountStatusClosed, nil
	case AccountStatusUnspecified:
		return AccountStatusUnspecified, nil
	}
	return "", fmt.Errorf("unknown account status: %s", s)
}

type Account struct {
	ID        int64
	PublicID  uuid.UUID
	UserID    uuid.UUID
	Currency  string
	Balance   decimal.Decimal
	Status    AccountStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Balance struct {
	AccountID string
	Amount    decimal.Decimal
	Currency  string
}

type CreateAccountRequest struct {
	UserID   uuid.UUID
	Currency string
}

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

type UpdateStatusRequest struct {
	AccountID uuid.UUID
	Status    AccountStatus
}
