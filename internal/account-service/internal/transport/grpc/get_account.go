package grpc

import (
	"context"

	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) GetAccount(_ context.Context, _ *accountpb.GetAccountRequest) (*accountpb.Account, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
