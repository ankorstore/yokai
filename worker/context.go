package worker

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TracerName is the workers tracer name.
const TracerName = "worker"

// CtxWorkerNameKey is a contextual struct key for the current worker name.
type CtxWorkerNameKey struct{}

// CtxWorkerExecutionIdKey is a contextual struct key for the current worker execution id.
type CtxWorkerExecutionIdKey struct{}

// CtxWorkerName returns the contextual [Worker] name.
func CtxWorkerName(ctx context.Context) string {
	if name, ok := ctx.Value(CtxWorkerNameKey{}).(string); ok {
		return name
	} else {
		return ""
	}
}

// CtxWorkerExecutionId returns the contextual [Worker] execution id.
func CtxWorkerExecutionId(ctx context.Context) string {
	if id, ok := ctx.Value(CtxWorkerExecutionIdKey{}).(string); ok {
		return id
	} else {
		return ""
	}
}

// CtxLogger returns the contextual [log.Logger].
func CtxLogger(ctx context.Context) *log.Logger {
	return log.CtxLogger(ctx)
}

// CtxTracer returns the contextual [oteltrace.Tracer].
func CtxTracer(ctx context.Context) oteltrace.Tracer {
	return trace.CtxTracerProvider(ctx).Tracer(TracerName)
}
