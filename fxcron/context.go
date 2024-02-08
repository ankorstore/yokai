package fxcron

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type CtxCronJobNameKey struct{}
type CtxCronJobExecutionIdKey struct{}

func CtxCronJobName(ctx context.Context) string {
	if name, ok := ctx.Value(CtxCronJobNameKey{}).(string); ok {
		return name
	} else {
		return ""
	}
}

func CtxCronJobExecutionId(ctx context.Context) string {
	if id, ok := ctx.Value(CtxCronJobExecutionIdKey{}).(string); ok {
		return id
	} else {
		return ""
	}
}

func CtxLogger(ctx context.Context) *log.Logger {
	return log.CtxLogger(ctx)
}

func CtxTracer(ctx context.Context) oteltrace.Tracer {
	return trace.CtxTracerProvider(ctx).Tracer(ModuleName)
}
