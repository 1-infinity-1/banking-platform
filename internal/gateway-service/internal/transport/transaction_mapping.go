package transport

import (
	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func toAPITransaction(t models.Transaction) api.Transaction {
	tx := api.Transaction{
		ID:             t.ID,
		ToAccountID:    t.ToAccountID,
		Amount:         t.Amount,
		Currency:       t.Currency,
		Status:         api.TransactionStatus(t.Status),
		IdempotencyKey: t.IdempotencyKey,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
	if t.FromAccountID != nil {
		tx.FromAccountID = api.NewOptNilUUID(*t.FromAccountID)
	} else {
		tx.FromAccountID = api.OptNilUUID{Set: true, Null: true}
	}
	return tx
}
