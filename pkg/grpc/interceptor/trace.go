package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type traceKey struct{}

type TraceContext struct {
	TraceID   string
	RequestID string
}

func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tc := TraceContext{}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("x-trace-id"); len(vals) > 0 {
				tc.TraceID = vals[0]
			}
			if vals := md.Get("x-request-id"); len(vals) > 0 {
				tc.RequestID = vals[0]
			}
		}

		ctx = context.WithValue(ctx, traceKey{}, tc)

		return handler(ctx, req)
	}
}

func TraceFromContext(ctx context.Context) TraceContext {
	tc, _ := ctx.Value(traceKey{}).(TraceContext)
	return tc
}
