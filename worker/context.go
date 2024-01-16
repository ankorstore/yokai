package worker

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const TracerName = "worker"

type CtxWorkerNameKey struct{}
type CtxWorkerExecutionIdKey struct{}

func CtxWorkerName(ctx context.Context) string {
	if name, ok := ctx.Value(CtxWorkerNameKey{}).(string); ok {
		return name
	} else {
		return ""
	}
}

func CtxWorkerExecutionId(ctx context.Context) string {
	if id, ok := ctx.Value(CtxWorkerExecutionIdKey{}).(string); ok {
		return id
	} else {
		return ""
	}
}

func CtxLogger(ctx context.Context) *log.Logger {
	return log.CtxLogger(ctx)
}

func CtxTracer(ctx context.Context) oteltrace.Tracer {
	return trace.CtxTracerProvider(ctx).Tracer(TracerName)
}
