package ledger

import (
	"errors"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
	ledgerpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/ledger"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoGetStatementRequest(p models.GetStatementParams) *ledgerpb.GetStatementRequest {
	return &ledgerpb.GetStatementRequest{
		AccountId: p.AccountID.String(),
		From:      timestamppb.New(p.From),
		To:        timestamppb.New(p.To),
	}
}

func toStatement(s *ledgerpb.Statement) (*models.Statement, error) {
	if s == nil {
		return nil, errors.New("empty statement")
	}
	accountID, err := uuid.Parse(s.GetAccountId())
	if err != nil {
		return nil, fmt.Errorf("parse statement account_id: %w", err)
	}
	entries := make([]models.LedgerEntry, 0, len(s.GetEntries()))
	for _, e := range s.GetEntries() {
		entry, mappingErr := toLedgerEntry(e)
		if mappingErr != nil {
			return nil, fmt.Errorf("toLedgerEntry: %w", mappingErr)
		}
		entries = append(entries, entry)
	}
	return &models.Statement{
		AccountID: accountID,
		Entries:   entries,
	}, nil
}

func toLedgerEntry(e *ledgerpb.LedgerEntry) (models.LedgerEntry, error) {
	if e == nil {
		return models.LedgerEntry{}, errors.New("empty ledger entry")
	}
	id, err := uuid.Parse(e.GetId())
	if err != nil {
		return models.LedgerEntry{}, fmt.Errorf("parse ledger entry id: %w", err)
	}
	txID, err := uuid.Parse(e.GetTransactionId())
	if err != nil {
		return models.LedgerEntry{}, fmt.Errorf("parse transaction_id: %w", err)
	}
	accountID, err := uuid.Parse(e.GetAccountId())
	if err != nil {
		return models.LedgerEntry{}, fmt.Errorf("parse ledger entry account_id: %w", err)
	}
	return models.LedgerEntry{
		ID:            id,
		TransactionID: txID,
		AccountID:     accountID,
		Type:          e.GetType(),
		Amount:        e.GetAmount(),
		Currency:      e.GetCurrency(),
		BalanceAfter:  e.GetBalanceAfter(),
		Description:   e.GetDescription(),
		OccurredAt:    e.GetOccurredAt().AsTime(),
		CreatedAt:     e.GetCreatedAt().AsTime(),
	}, nil
}

func mapGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("unexpected gRPC error: %w", err)
	}
	switch st.Code() { //nolint:exhaustive // only meaningful codes handled; default covers the rest
	case codes.NotFound:
		return models.NewNotFoundError(st.Message(), err)
	case codes.InvalidArgument:
		return models.NewValidationError("", st.Message(), err)
	case codes.Unauthenticated:
		return models.NewUnauthorizedError(st.Message(), err)
	default:
		return fmt.Errorf("ledger service error: %w", err)
	}
}
