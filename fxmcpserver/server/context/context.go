package context

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type CtxRequestIdKey struct{}
type CtxSessionIdKey struct{}
type CtxRootSpanKey struct{}
type CtxStartTimeKey struct{}

// WithRequestID adds a given request id to a given context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, CtxRequestIdKey{}, requestID)
}

// CtxRequestId returns the request id from a given context.
func CtxRequestId(ctx context.Context) string {
	if rid, ok := ctx.Value(CtxRequestIdKey{}).(string); ok {
		return rid
	}

	return ""
}

// WithSessionID adds a given session id to a given context.
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, CtxSessionIdKey{}, sessionID)
}

// CtxSessionID returns the session id from a given context.
func CtxSessionID(ctx context.Context) string {
	if sid, ok := ctx.Value(CtxSessionIdKey{}).(string); ok {
		return sid
	}

	return ""
}

// WithRootSpan adds a root span to a given context.
func WithRootSpan(ctx context.Context, span trace.Span) context.Context {
	return context.WithValue(ctx, CtxRootSpanKey{}, span)
}

// CtxRootSpan returns the root span from a given context.
func CtxRootSpan(ctx context.Context) trace.Span {
	if span, ok := ctx.Value(CtxRootSpanKey{}).(trace.Span); ok {
		return span
	}

	return trace.SpanFromContext(ctx)
}

// WithStartTime adds a start time to a given context.
func WithStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, CtxStartTimeKey{}, t)
}

// CtxStartTime returns the start time from a given context.
func CtxStartTime(ctx context.Context) time.Time {
	if t, ok := ctx.Value(CtxStartTimeKey{}).(time.Time); ok {
		return t
	}

	return time.Now()
}
