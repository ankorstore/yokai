package httpserver

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TraceSpanAttributeHttpRequestId is a span attribute representing the request id.
const TraceSpanAttributeHttpRequestId = "guid:x-request-id"

// AnnotateTracerProvider extend adds the [TracerProviderRequestIdAnnotator] span processor to the provided [oteltrace.TracerProvider].
func AnnotateTracerProvider(base oteltrace.TracerProvider) oteltrace.TracerProvider {
	if tp, ok := base.(*otelsdktrace.TracerProvider); ok {
		tp.RegisterSpanProcessor(NewTracerProviderRequestIdAnnotator())

		return tp
	}

	return base
}

// TracerProviderRequestIdAnnotator is a span processor to add the contextual request id as span attribute.
type TracerProviderRequestIdAnnotator struct{}

// NewTracerProviderRequestIdAnnotator returns a new [TracerProviderRequestIdAnnotator].
func NewTracerProviderRequestIdAnnotator() *TracerProviderRequestIdAnnotator {
	return &TracerProviderRequestIdAnnotator{}
}

// OnStart adds the contextual request id as span attribute.
func (a *TracerProviderRequestIdAnnotator) OnStart(ctx context.Context, s otelsdktrace.ReadWriteSpan) {
	if rid, ok := ctx.Value(CtxRequestIdKey{}).(string); ok {
		s.SetAttributes(attribute.String(TraceSpanAttributeHttpRequestId, rid))
	}
}

// Shutdown performs no operations.
func (a *TracerProviderRequestIdAnnotator) Shutdown(context.Context) error {
	return nil
}

// ForceFlush performs no operations.
func (a *TracerProviderRequestIdAnnotator) ForceFlush(context.Context) error {
	return nil
}

// OnEnd performs no operations.
func (a *TracerProviderRequestIdAnnotator) OnEnd(otelsdktrace.ReadOnlySpan) {
	// noop
}
