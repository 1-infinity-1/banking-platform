package grpc

import (
	"context"

	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) Replenish(
	_ context.Context,
	_ *transactionpb.ReplenishRequest,
) (*transactionpb.Transaction, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
