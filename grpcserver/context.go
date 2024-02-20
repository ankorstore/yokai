package grpcserver

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TracerName is the grpcserver tracer name.
const TracerName = "grpcserver"

// CtxLogger returns the contextual [log.Logger].
func CtxLogger(ctx context.Context) *log.Logger {
	return log.CtxLogger(ctx)
}

// CtxTracer returns the contextual [oteltrace.Tracer].
func CtxTracer(ctx context.Context) oteltrace.Tracer {
	return trace.CtxTracerProvider(ctx).Tracer(TracerName)
}
