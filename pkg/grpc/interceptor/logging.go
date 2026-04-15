package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggingUnaryServerInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fields := make([]any, 0, 8)
		fields = append(fields,
			slog.String("protocol", "grpc"),
			slog.String("method", info.FullMethod),
		)

		if peerAddr, ok := peer.FromContext(ctx); ok {
			fields = append(fields, slog.String("peer", peerAddr.String()))
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("x-request-id"); len(vals) > 0 {
				fields = append(fields, slog.String("request_id", vals[0]))
			}
			if vals := md.Get("x-trace-id"); len(vals) > 0 {
				fields = append(fields, slog.String("trace_id", vals[0]))
			}
		}

		resp, err := handler(ctx, req)
		if err != nil {
			statusCode := status.Convert(err)

			fields = append(fields,
				slog.String("grpc_code", statusCode.Code().String()),
				slog.Duration("duration", time.Since(time.Now().UTC())),
				slog.String("message", statusCode.Message()),
			)

			log.Error("grpc request failed", fields...)

			return resp, err
		}

		log.Info("grpc request handled", fields...)

		return resp, nil
	}
}
