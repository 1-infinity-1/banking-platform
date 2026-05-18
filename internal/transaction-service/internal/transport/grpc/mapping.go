package grpc

import (
	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoTransactionStatus(s models.TransactionStatus) transactionpb.TransactionStatus {
	switch s {
	case models.TransactionStatusPending:
		return transactionpb.TransactionStatus_TRANSACTION_STATUS_PENDING
	case models.TransactionStatusCompleted:
		return transactionpb.TransactionStatus_TRANSACTION_STATUS_COMPLETED
	case models.TransactionStatusFailed:
		return transactionpb.TransactionStatus_TRANSACTION_STATUS_FAILED
	case models.TransactionStatusCancelled:
		return transactionpb.TransactionStatus_TRANSACTION_STATUS_CANCELLED
	case models.TransactionStatusUnspecified:
		return transactionpb.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED
	}
	return transactionpb.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED
}

func toProtoTransaction(t *models.Transaction) *transactionpb.Transaction {
	proto := &transactionpb.Transaction{
		Id:             t.PublicID.String(),
		ToAccountId:    t.ToAccountID.String(),
		Amount:         t.Amount.String(),
		Currency:       t.Currency,
		Status:         toProtoTransactionStatus(t.Status),
		IdempotencyKey: t.IdempotencyKey,
		CreatedAt:      timestamppb.New(t.CreatedAt),
		UpdatedAt:      timestamppb.New(t.UpdatedAt),
	}
	if t.FromAccountID != nil {
		from := t.FromAccountID.String()
		proto.FromAccountId = &from
	}
	return proto
}

func toProtoTransactionsList(transactions []*models.Transaction) *transactionpb.TransactionsList {
	list := &transactionpb.TransactionsList{
		Transactions: make([]*transactionpb.Transaction, 0, len(transactions)),
	}
	for _, t := range transactions {
		list.Transactions = append(list.Transactions, toProtoTransaction(t))
	}
	return list
}
