package grpc

import (
	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
)

//nolint:unused // scaffold: used when implementing TODO handlers
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
