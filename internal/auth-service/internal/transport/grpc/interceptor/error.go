package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/models"
	"github.com/1-infinity-1/banking-platform/pkg/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		res, err := handler(ctx, req)
		if err == nil {
			return res, nil
		}

		if _, ok := status.FromError(err); ok {
			return res, err
		}

		var (
			notFoundErr   *models.NotFoundError
			invalidParams *models.InvalidParamsError
			businessErr   *models.BusinessError
		)
		switch {
		case errors.As(err, &notFoundErr):
			return res, status.Error(codes.NotFound, notFoundErr.Error())
		case errors.As(err, &invalidParams):
			return res, status.Error(codes.InvalidArgument, invalidParams.Error())
		case errors.As(err, &businessErr):
			return res, status.Error(codes.FailedPrecondition, businessErr.Error())
		default:
			tc := interceptor.TraceFromContext(ctx)
			log.Error("unexpected internal error",
				slog.String("method", info.FullMethod),
				slog.String("trace_id", tc.TraceID),
				slog.String("request_id", tc.RequestID),
				slog.Any("error", err),
			)

			return res, status.Error(codes.Internal, models.ErrInternal.Error())
		}

	}
}
