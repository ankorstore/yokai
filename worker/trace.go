package worker

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// AnnotateTracerProvider extends a provided [oteltrace.TracerProvider] spans with worker execution attributes.
func AnnotateTracerProvider(base oteltrace.TracerProvider) oteltrace.TracerProvider {
	if tp, ok := base.(*otelsdktrace.TracerProvider); ok {
		tp.RegisterSpanProcessor(NewTracerProviderWorkerAnnotator())

		return tp
	}

	return base
}

// TracerProviderWorkerAnnotator is the [oteltrace.TracerProvider] workers annotator, implementing [otelsdktrace.SpanProcessor].
type TracerProviderWorkerAnnotator struct{}

// NewTracerProviderWorkerAnnotator returns a new [TracerProviderWorkerAnnotator].
func NewTracerProviderWorkerAnnotator() *TracerProviderWorkerAnnotator {
	return &TracerProviderWorkerAnnotator{}
}

// OnStart adds worker execution attributes to a given [otelsdktrace.ReadWriteSpan].
func (a *TracerProviderWorkerAnnotator) OnStart(ctx context.Context, s otelsdktrace.ReadWriteSpan) {
	name := CtxWorkerName(ctx)
	if name != "" {
		s.SetAttributes(attribute.String(TraceSpanAttributeWorkerName, name))
	}

	executionId := CtxWorkerExecutionId(ctx)
	if executionId != "" {
		s.SetAttributes(attribute.String(TraceSpanAttributeWorkerExecutionId, executionId))
	}
}

// Shutdown is just for [otelsdktrace.SpanProcessor] compliance.
func (a *TracerProviderWorkerAnnotator) Shutdown(context.Context) error {
	return nil
}

// ForceFlush is just for [otelsdktrace.SpanProcessor] compliance.
func (a *TracerProviderWorkerAnnotator) ForceFlush(context.Context) error {
	return nil
}

// OnEnd is just for [otelsdktrace.SpanProcessor] compliance.
func (a *TracerProviderWorkerAnnotator) OnEnd(otelsdktrace.ReadOnlySpan) {
	// noop
}
