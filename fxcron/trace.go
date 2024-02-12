package fxcron

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// AnnotateTracerProvider extends a provided [oteltrace.TracerProvider] spans with cron jobs execution attributes.
func AnnotateTracerProvider(base oteltrace.TracerProvider) oteltrace.TracerProvider {
	if tp, ok := base.(*otelsdktrace.TracerProvider); ok {
		tp.RegisterSpanProcessor(NewTracerProviderCronJobAnnotator())

		return tp
	}

	return base
}

// TracerProviderCronJobAnnotator is the [oteltrace.TracerProvider] cron jobs annotator, implementing [otelsdktrace.SpanProcessor].
type TracerProviderCronJobAnnotator struct{}

// NewTracerProviderCronJobAnnotator returns a new [TracerProviderWorkerAnnotator].
func NewTracerProviderCronJobAnnotator() *TracerProviderCronJobAnnotator {
	return &TracerProviderCronJobAnnotator{}
}

// OnStart adds cron job execution attributes to a given [otelsdktrace.ReadWriteSpan].
func (a *TracerProviderCronJobAnnotator) OnStart(ctx context.Context, s otelsdktrace.ReadWriteSpan) {
	name := CtxCronJobName(ctx)
	if name != "" {
		s.SetAttributes(attribute.String(TraceSpanAttributeCronJobName, name))
	}

	executionId := CtxCronJobExecutionId(ctx)
	if executionId != "" {
		s.SetAttributes(attribute.String(TraceSpanAttributeCronJobExecutionId, executionId))
	}
}

// Shutdown is just for [otelsdktrace.SpanProcessor] compliance.
func (a *TracerProviderCronJobAnnotator) Shutdown(context.Context) error {
	return nil
}

// ForceFlush is just for [otelsdktrace.SpanProcessor] compliance.
func (a *TracerProviderCronJobAnnotator) ForceFlush(context.Context) error {
	return nil
}

// OnEnd is just for [otelsdktrace.SpanProcessor] compliance.
func (a *TracerProviderCronJobAnnotator) OnEnd(otelsdktrace.ReadOnlySpan) {
	// noop
}
