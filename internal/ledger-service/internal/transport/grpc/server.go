package grpc

import (
	"context"
	"time"

	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/models"
	ledgerpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/ledger"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type LedgerService interface {
	GetStatement(ctx context.Context, accountID uuid.UUID, from, to time.Time) (models.Statement, error)
}

type serverAPI struct {
	ledgerpb.UnimplementedLedgerServiceServer

	ledgerSvc LedgerService
}

func NewServerAPI(gRPC *grpc.Server, ledgerSvc LedgerService) {
	srv := &serverAPI{
		ledgerSvc: ledgerSvc,
	}
	ledgerpb.RegisterLedgerServiceServer(gRPC, srv)
}
