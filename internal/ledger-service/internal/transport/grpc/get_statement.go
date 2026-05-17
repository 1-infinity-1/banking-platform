package grpc

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/models"
	ledgerpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/ledger"
	"github.com/google/uuid"
)

func (s *serverAPI) GetStatement(
	ctx context.Context,
	req *ledgerpb.GetStatementRequest,
) (*ledgerpb.Statement, error) {
	if req.GetAccountId() == "" {
		return nil, models.NewInvalidParamsError("account_id", "is empty")
	}

	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, models.NewInvalidParamsError("account_id", fmt.Sprintf("invalid UUID: %s", err))
	}

	if req.GetFrom() == nil {
		return nil, models.NewInvalidParamsError("from", "is required")
	}

	if req.GetTo() == nil {
		return nil, models.NewInvalidParamsError("to", "is required")
	}

	from := req.GetFrom().AsTime()
	to := req.GetTo().AsTime()

	stmt, err := s.ledgerSvc.GetStatement(ctx, accountID, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetStatement: %w", err)
	}

	return domainStatementToProto(stmt), nil
}

func domainStatementToProto(stmt models.Statement) *ledgerpb.Statement {
	entries := make([]*ledgerpb.LedgerEntry, 0, len(stmt.Entries))
	for _, e := range stmt.Entries {
		entries = append(entries, domainEntryToProto(e))
	}

	return &ledgerpb.Statement{
		AccountId: stmt.AccountID.String(),
		Entries:   entries,
	}
}
