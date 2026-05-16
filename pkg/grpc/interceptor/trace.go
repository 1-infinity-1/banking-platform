package interceptor

import (
	"context"

	"github.com/1-infinity-1/banking-platform/pkg/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		tc := trace.Context{}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("x-trace-id"); len(vals) > 0 {
				tc.TraceID = vals[0]
			}
			if vals := md.Get("x-request-id"); len(vals) > 0 {
				tc.RequestID = vals[0]
			}
		}

		ctx = trace.ToContext(ctx, tc)

		return handler(ctx, req)
	}
}
