package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountStatusUnspecified AccountStatus = "unspecified"
	AccountStatusActive      AccountStatus = "active"
	AccountStatusBlocked     AccountStatus = "blocked"
	AccountStatusClosed      AccountStatus = "closed"
)

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Currency  string
	Balance   string
	Status    AccountStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Balance struct {
	AccountID uuid.UUID
	Amount    string
	Currency  string
}

type CreateAccountParams struct {
	UserID   uuid.UUID
	Currency string
}

type GetUserAccountsParams struct {
	UserID uuid.UUID
}

type GetAccountParams struct {
	AccountID uuid.UUID
}

type GetBalanceParams struct {
	AccountID uuid.UUID
}

type UpdateAccountStatusParams struct {
	AccountID uuid.UUID
	Status    AccountStatus
}
