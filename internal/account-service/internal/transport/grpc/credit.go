package grpc

import (
	"context"

	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) Credit(_ context.Context, _ *accountpb.CreditRequest) (*accountpb.CreditResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
