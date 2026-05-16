package trace

import "context"

type traceKey struct{}

type Context struct {
	TraceID   string
	RequestID string
}

func FromContext(ctx context.Context) Context {
	tc, _ := ctx.Value(traceKey{}).(Context)
	return tc
}

func ToContext(ctx context.Context, tc Context) context.Context {
	return context.WithValue(ctx, traceKey{}, tc)
}
