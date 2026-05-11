package interceptor

import (
	"context"
	"log/slog"
	"time"

	"github.com/1-infinity-1/banking-platform/pkg/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const initialFieldsCapacity = 8

func LoggingUnaryServerInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now().UTC()

		fields := make([]any, 0, initialFieldsCapacity)
		fields = append(fields,
			slog.String("protocol", "grpc"),
			slog.String("method", info.FullMethod),
		)

		if tc := trace.FromContext(ctx); tc.TraceID != "" || tc.RequestID != "" {
			if tc.TraceID != "" {
				fields = append(fields, slog.String("trace_id", tc.TraceID))
			}
			if tc.RequestID != "" {
				fields = append(fields, slog.String("request_id", tc.RequestID))
			}
		}

		resp, err := handler(ctx, req)
		if err != nil {
			statusCode := status.Convert(err)

			fields = append(fields,
				slog.String("grpc_code", statusCode.Code().String()),
				slog.Duration("duration", time.Since(startTime)),
				slog.String("message", statusCode.Message()),
			)

			log.Error("grpc request failed", fields...)

			return resp, err
		}

		fields = append(fields, slog.Duration("duration", time.Since(startTime)))
		log.Info("grpc request handled", fields...)

		return resp, nil
	}
}
