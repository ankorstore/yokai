package trace

import (
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Options are options for the [TracerProviderFactory] implementations.
type Options struct {
	Global         bool
	Resource       *resource.Resource
	Sampler        trace.Sampler
	SpanProcessors []trace.SpanProcessor
}

// DefaultTracerProviderOptions are the default options used in the [TracerProviderFactory].
func DefaultTracerProviderOptions() Options {
	return Options{
		Global:         true,
		Resource:       resource.Default(),
		Sampler:        NewParentBasedAlwaysOnSampler(),
		SpanProcessors: []trace.SpanProcessor{},
	}
}

// TracerProviderOption are functional options for the [TracerProviderFactory] implementations.
type TracerProviderOption func(o *Options)

// Global is used to set the [OTEL TracerProvider] as global.
//
// [OTEL TracerProvider]: https://github.com/open-telemetry/opentelemetry-go
func Global(b bool) TracerProviderOption {
	return func(o *Options) {
		o.Global = b
	}
}

// WithResource is used to set the resource to use by the [OTEL TracerProvider].
//
// [OTEL TracerProvider]: https://github.com/open-telemetry/opentelemetry-go
func WithResource(r *resource.Resource) TracerProviderOption {
	return func(o *Options) {
		o.Resource = r
	}
}

// WithSampler is used to set the sampler to use by the [OTEL TracerProvider].
//
// [OTEL TracerProvider]: https://github.com/open-telemetry/opentelemetry-go
func WithSampler(s trace.Sampler) TracerProviderOption {
	return func(o *Options) {
		o.Sampler = s
	}
}

// WithSpanProcessor is used to set the span processor to use by the [OTEL TracerProvider].
//
// [OTEL TracerProvider]: https://github.com/open-telemetry/opentelemetry-go
func WithSpanProcessor(spanProcessor trace.SpanProcessor) TracerProviderOption {
	return func(o *Options) {
		o.SpanProcessors = append(o.SpanProcessors, spanProcessor)
	}
}
