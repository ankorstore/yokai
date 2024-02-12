package fxcron

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// CtxCronJobNameKey is a contextual struct key.
type CtxCronJobNameKey struct{}

// CtxCronJobExecutionIdKey is a contextual struct key.
type CtxCronJobExecutionIdKey struct{}

// CtxCronJobName returns the contextual cron job name.
func CtxCronJobName(ctx context.Context) string {
	if name, ok := ctx.Value(CtxCronJobNameKey{}).(string); ok {
		return name
	} else {
		return ""
	}
}

// CtxCronJobExecutionId returns the contextual cron job execution id.
func CtxCronJobExecutionId(ctx context.Context) string {
	if id, ok := ctx.Value(CtxCronJobExecutionIdKey{}).(string); ok {
		return id
	} else {
		return ""
	}
}

// CtxLogger returns the contextual logger.
func CtxLogger(ctx context.Context) *log.Logger {
	return log.CtxLogger(ctx)
}

// CtxTracer returns the contextual tracer.
func CtxTracer(ctx context.Context) oteltrace.Tracer {
	return trace.CtxTracerProvider(ctx).Tracer(ModuleName)
}
