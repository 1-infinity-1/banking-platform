package grpc

import (
	"context"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type TransactionService interface {
	Transfer(ctx context.Context, req models.TransferRequest) (*models.Transaction, error)
	Replenish(ctx context.Context, req models.ReplenishRequest) (*models.Transaction, error)
	GetHistory(ctx context.Context, req models.GetHistoryRequest) ([]*models.Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
}

type serverAPI struct {
	transactionpb.UnimplementedTransactionServiceServer
	svc TransactionService
}

func NewServerAPI(gRPC *grpc.Server, svc TransactionService) {
	srv := &serverAPI{svc: svc}
	transactionpb.RegisterTransactionServiceServer(gRPC, srv)
}
