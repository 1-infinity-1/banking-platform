package grpc

import (
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/models"
	ledgerpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/ledger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func domainEntryToProto(e models.LedgerEntry) *ledgerpb.LedgerEntry {
	return &ledgerpb.LedgerEntry{
		Id:            e.PublicID.String(),
		TransactionId: e.TransactionID.String(),
		AccountId:     e.AccountID.String(),
		Type:          string(e.Type),
		Amount:        e.Amount.String(),
		Currency:      e.Currency,
		BalanceAfter:  e.BalanceAfter.String(),
		Description:   e.Description,
		OccurredAt:    timestamppb.New(e.OccurredAt),
		CreatedAt:     timestamppb.New(e.CreatedAt),
	}
}
