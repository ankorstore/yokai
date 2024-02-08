package fxcron

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func AnnotateTracerProvider(base oteltrace.TracerProvider) oteltrace.TracerProvider {
	if tp, ok := base.(*otelsdktrace.TracerProvider); ok {
		tp.RegisterSpanProcessor(NewTracerProviderCronJobAnnotator())

		return tp
	}

	return base
}

type TracerProviderCronJobAnnotator struct{}

func NewTracerProviderCronJobAnnotator() *TracerProviderCronJobAnnotator {
	return &TracerProviderCronJobAnnotator{}
}

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
func (a *TracerProviderCronJobAnnotator) Shutdown(context.Context) error {
	return nil
}

func (a *TracerProviderCronJobAnnotator) ForceFlush(context.Context) error {
	return nil
}

func (a *TracerProviderCronJobAnnotator) OnEnd(otelsdktrace.ReadOnlySpan) {
	// noop
}
