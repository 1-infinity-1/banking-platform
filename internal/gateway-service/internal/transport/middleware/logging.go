package middleware

import (
	"log/slog"
	"net/http"
	"time"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/pkg/trace"
	"github.com/ogen-go/ogen/middleware"
)

const initialFieldsCapacity = 8

func Logging(log *slog.Logger) api.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		startTime := time.Now().UTC()

		fields := make([]any, 0, initialFieldsCapacity)
		fields = append(fields,
			slog.String("protocol", "http"),
			slog.String("operation_id", req.OperationName),
		)

		if tc := trace.FromContext(req.Context); tc.TraceID != "" || tc.RequestID != "" {
			if tc.TraceID != "" {
				fields = append(fields, slog.String("trace_id", tc.TraceID))
			}
			if tc.RequestID != "" {
				fields = append(fields, slog.String("request_id", tc.RequestID))
			}
		}

		resp, err := next(req)
		fields = append(fields, slog.Duration("duration", time.Since(startTime)))

		if err != nil {
			fields = append(fields, slog.String("message", err.Error()))
			log.Error("http request failed", fields...)
			return resp, err
		}

		if errResp, ok := resp.Type.(*api.ErrorStatusCode); ok {
			fields = append(fields,
				slog.Int("status_code", errResp.StatusCode),
				slog.String("code", errResp.Response.Code),
				slog.String("message", errResp.Response.Message),
			)
			if errResp.StatusCode >= http.StatusInternalServerError {
				log.Error("http request failed", fields...)
			} else {
				log.Info("http request handled", fields...)
			}
			return resp, nil
		}

		log.Info("http request handled", fields...)
		return resp, nil
	}
}
