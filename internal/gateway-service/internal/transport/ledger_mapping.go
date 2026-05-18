package transport

import (
	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func toAPILedgerEntry(e models.LedgerEntry) api.LedgerEntry {
	entry := api.LedgerEntry{
		ID:            e.ID,
		TransactionID: e.TransactionID,
		AccountID:     e.AccountID,
		Type:          e.Type,
		Amount:        e.Amount,
		Currency:      e.Currency,
		BalanceAfter:  e.BalanceAfter,
		OccurredAt:    e.OccurredAt,
		CreatedAt:     e.CreatedAt,
	}
	if e.Description != "" {
		entry.Description = api.NewOptString(e.Description)
	}
	return entry
}

func toAPIStatement(s models.Statement) api.Statement {
	entries := make([]api.LedgerEntry, 0, len(s.Entries))
	for _, e := range s.Entries {
		entries = append(entries, toAPILedgerEntry(e))
	}
	return api.Statement{
		AccountID: s.AccountID,
		Entries:   entries,
	}
}
