package middleware

import (
	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/pkg/trace"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/middleware"
)

func Trace() api.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		traceID := req.Raw.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		requestID := req.Raw.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := trace.ToContext(req.Context, trace.TraceContext{
			TraceID:   traceID,
			RequestID: requestID,
		})
		req.SetContext(ctx)

		return next(req)
	}
}
