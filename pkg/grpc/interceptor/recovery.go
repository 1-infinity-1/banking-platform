package interceptor

import (
	"context"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryUnaryServerInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("grpc panic recovered",
					slog.String("protocol", "grpc"),
					slog.String("method", info.FullMethod),
					slog.Any("panic", r),
					slog.String("stack_trace", string(debug.Stack())),
				)

				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
