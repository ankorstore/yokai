package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// TracerName is the default tracer name.
const TracerName = "yokai"

// CtxKey is a contextual struct key.
type CtxKey struct{}

// WithContext appends to a given context a [OTEL TracerProvider].
//
// [OTEL TracerProvider]: https://github.com/open-telemetry/opentelemetry-go
func WithContext(ctx context.Context, tp trace.TracerProvider) context.Context {
	return context.WithValue(ctx, CtxKey{}, tp)
}

// CtxTracerProvider retrieves an [OTEL TracerProvider] from a provided context (or returns the default one if missing).
//
// [OTEL TracerProvider]: https://github.com/open-telemetry/opentelemetry-go
func CtxTracerProvider(ctx context.Context) trace.TracerProvider {
	if tp, ok := ctx.Value(CtxKey{}).(trace.TracerProvider); ok {
		return tp
	} else {
		return otel.GetTracerProvider()
	}
}

// CtxTracer returns the default contextual [OTEL Tracer] from a provided context (or returns the default one if missing).
//
// [OTEL Tracer]: https://go.opentelemetry.io/otel/trace
func CtxTracer(ctx context.Context) trace.Tracer {
	return CtxTracerProvider(ctx).Tracer(TracerName)
}
