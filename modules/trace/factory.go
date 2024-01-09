package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

// TracerProviderFactory is the interface for [OTEL TracerProvider] factories.
//
// [OTEL TracerProvider]: https://github.com/open-telemetry/opentelemetry-go
type TracerProviderFactory interface {
	Create(options ...TracerProviderOption) (*trace.TracerProvider, error)
}

// DefaultTracerProviderFactory is the default [TracerProviderFactory] implementation.
type DefaultTracerProviderFactory struct{}

// NewDefaultTracerProviderFactory returns a [DefaultTracerProviderFactory], implementing [TracerProviderFactory].
func NewDefaultTracerProviderFactory() TracerProviderFactory {
	return &DefaultTracerProviderFactory{}
}

// Create returns a new [OTEL TracerProvider], and accepts a list of [TracerProviderOption].
//
// For example:
//
//	tp, _ := trace.NewDefaultTracerProviderFactory().Create()
//
//	is equivalent to:
//	tp, _ = trace.NewDefaultTracerProviderFactory().Create(
//		trace.Global(true),                                       // set the tracer provider as global
//		trace.WithResource(resource.Default()),                   // use the default resource
//		trace.WithSampler(trace.NewParentBasedAlwaysOnSampler()), // use parent based always on sampling
//		trace.WithSpanProcessor(trace.NewNoopSpanProcessor()),    // use noop processor (void trace spans)
//	)
//
// [OTEL TracerProvider]: https://github.com/open-telemetry/opentelemetry-go
func (f *DefaultTracerProviderFactory) Create(options ...TracerProviderOption) (*trace.TracerProvider, error) {
	appliedOptions := DefaultTracerProviderOptions()
	for _, opt := range options {
		opt(&appliedOptions)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(appliedOptions.Resource),
		trace.WithSampler(appliedOptions.Sampler),
	)

	if len(appliedOptions.SpanProcessors) == 0 {
		tracerProvider.RegisterSpanProcessor(NewNoopSpanProcessor())
	} else {
		for _, processor := range appliedOptions.SpanProcessors {
			tracerProvider.RegisterSpanProcessor(processor)
		}
	}

	if appliedOptions.Global {
		otel.SetTracerProvider(tracerProvider)

		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		)
	}

	return tracerProvider, nil
}
