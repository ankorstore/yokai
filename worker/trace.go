package worker

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func AnnotateTracerProvider(base oteltrace.TracerProvider) oteltrace.TracerProvider {
	if tp, ok := base.(*otelsdktrace.TracerProvider); ok {
		tp.RegisterSpanProcessor(NewTracerProviderWorkerAnnotator())

		return tp
	}

	return base
}

type TracerProviderWorkerAnnotator struct{}

func NewTracerProviderWorkerAnnotator() *TracerProviderWorkerAnnotator {
	return &TracerProviderWorkerAnnotator{}
}

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
func (a *TracerProviderWorkerAnnotator) Shutdown(context.Context) error {
	return nil
}

func (a *TracerProviderWorkerAnnotator) ForceFlush(context.Context) error {
	return nil
}

func (a *TracerProviderWorkerAnnotator) OnEnd(otelsdktrace.ReadOnlySpan) {
	// noop
}
