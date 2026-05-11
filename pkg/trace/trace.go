package trace

import "context"

type traceKey struct{}

type TraceContext struct {
	TraceID   string
	RequestID string
}

func FromContext(ctx context.Context) TraceContext {
	tc, _ := ctx.Value(traceKey{}).(TraceContext)
	return tc
}

func ToContext(ctx context.Context, tc TraceContext) context.Context {
	return context.WithValue(ctx, traceKey{}, tc)
}
