package grpc

import (
	"context"

	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req models.CreateAccountRequest) (*models.Account, error)
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*models.Account, error)
	GetAccount(ctx context.Context, accountID uuid.UUID) (*models.Account, error)
	GetBalance(ctx context.Context, accountID uuid.UUID) (*models.Balance, error)
	UpdateStatus(ctx context.Context, req models.UpdateStatusRequest) (*models.Account, error)
	Debit(ctx context.Context, req models.DebitRequest) (*models.DebitResult, error)
	Credit(ctx context.Context, req models.CreditRequest) (*models.CreditResult, error)
}

type serverAPI struct {
	accountpb.UnimplementedAccountServiceServer

	svc AccountService
}

func NewServerAPI(gRPC *grpc.Server, svc AccountService) {
	srv := &serverAPI{svc: svc}
	accountpb.RegisterAccountServiceServer(gRPC, srv)
}
