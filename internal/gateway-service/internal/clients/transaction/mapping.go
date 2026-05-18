package transaction

import (
	"errors"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toProtoTransferRequest(p models.TransferParams) *transactionpb.TransferRequest {
	return &transactionpb.TransferRequest{
		FromAccountId:  p.FromAccountID.String(),
		ToAccountId:    p.ToAccountID.String(),
		Amount:         p.Amount,
		Currency:       p.Currency,
		IdempotencyKey: p.IdempotencyKey,
	}
}

func toProtoReplenishRequest(p models.ReplenishParams) *transactionpb.ReplenishRequest {
	return &transactionpb.ReplenishRequest{
		ToAccountId:    p.ToAccountID.String(),
		Amount:         p.Amount,
		Currency:       p.Currency,
		IdempotencyKey: p.IdempotencyKey,
	}
}

func toProtoGetHistoryRequest(p models.GetHistoryParams) *transactionpb.GetHistoryRequest {
	return &transactionpb.GetHistoryRequest{
		AccountId: p.AccountID.String(),
		Limit:     p.Limit,
		Offset:    p.Offset,
	}
}

func toProtoGetTransactionRequest(p models.GetTransactionParams) *transactionpb.GetTransactionRequest {
	return &transactionpb.GetTransactionRequest{TransactionId: p.TransactionID.String()}
}

func toTransaction(t *transactionpb.Transaction) (*models.Transaction, error) {
	if t == nil {
		return nil, errors.New("empty transaction")
	}
	id, err := uuid.Parse(t.GetId())
	if err != nil {
		return nil, fmt.Errorf("parse transaction id: %w", err)
	}
	toAccountID, err := uuid.Parse(t.GetToAccountId())
	if err != nil {
		return nil, fmt.Errorf("parse to_account_id: %w", err)
	}

	tx := &models.Transaction{
		ID:             id,
		ToAccountID:    toAccountID,
		Amount:         t.GetAmount(),
		Currency:       t.GetCurrency(),
		Status:         toTransactionStatus(t.GetStatus()),
		IdempotencyKey: t.GetIdempotencyKey(),
		CreatedAt:      t.GetCreatedAt().AsTime(),
		UpdatedAt:      t.GetUpdatedAt().AsTime(),
	}

	if fromID := t.GetFromAccountId(); fromID != "" {
		parsed, parseErr := uuid.Parse(fromID)
		if parseErr != nil {
			return nil, fmt.Errorf("parse from_account_id: %w", parseErr)
		}
		tx.FromAccountID = &parsed
	}

	return tx, nil
}

func toTransactionStatus(s transactionpb.TransactionStatus) models.TransactionStatus {
	switch s { //nolint:exhaustive // unspecified handled by default
	case transactionpb.TransactionStatus_TRANSACTION_STATUS_PENDING:
		return models.TransactionStatusPending
	case transactionpb.TransactionStatus_TRANSACTION_STATUS_COMPLETED:
		return models.TransactionStatusCompleted
	case transactionpb.TransactionStatus_TRANSACTION_STATUS_FAILED:
		return models.TransactionStatusFailed
	case transactionpb.TransactionStatus_TRANSACTION_STATUS_CANCELLED:
		return models.TransactionStatusCancelled
	default:
		return models.TransactionStatusUnspecified
	}
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
		return fmt.Errorf("transaction service error: %w", err)
	}
}
