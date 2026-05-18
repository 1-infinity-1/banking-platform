package transaction

import (
	"context"
	"fmt"

	transactionpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/transaction"
	"github.com/1-infinity-1/banking-platform/pkg/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn *grpc.ClientConn
	svc  transactionpb.TransactionServiceClient
}

func NewClient(host, port string) (*Client, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(traceInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}
	return &Client{
		conn: conn,
		svc:  transactionpb.NewTransactionServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func traceInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tc := trace.FromContext(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx,
			"x-trace-id", tc.TraceID,
			"x-request-id", tc.RequestID,
		)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
