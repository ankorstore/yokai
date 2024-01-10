package trace

import (
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	ParentBasedAlwaysOn     = "parent-based-always-on"      // parent based always on sampling
	ParentBasedAlwaysOff    = "parent-based-always-off"     // parent based always off sampling
	ParentBasedTraceIdRatio = "parent-based-trace-id-ratio" // parent based trace id ratio sampling
	AlwaysOn                = "always-on"                   // always on sampling
	AlwaysOff               = "always-off"                  // always off sampling
	TraceIdRatio            = "trace-id-ratio"              // trace id ratio sampling
)

// NewParentBasedAlwaysOnSampler returns a [otelsdktrace.Sampler] with parent based always on sampling.
func NewParentBasedAlwaysOnSampler() otelsdktrace.Sampler {
	return otelsdktrace.ParentBased(otelsdktrace.AlwaysSample())
}

// NewParentBasedAlwaysOffSampler returns a [otelsdktrace.Sampler] with parent based always off sampling.
func NewParentBasedAlwaysOffSampler() otelsdktrace.Sampler {
	return otelsdktrace.ParentBased(otelsdktrace.NeverSample())
}

// NewParentBasedTraceIdRatioSampler returns a [otelsdktrace.Sampler] with parent based trace id ratio sampling.
func NewParentBasedTraceIdRatioSampler(ratio float64) otelsdktrace.Sampler {
	return otelsdktrace.ParentBased(otelsdktrace.TraceIDRatioBased(ratio))
}

// NewAlwaysOnSampler returns a [otelsdktrace.Sampler] with always on sampling.
func NewAlwaysOnSampler() otelsdktrace.Sampler {
	return otelsdktrace.AlwaysSample()
}

// NewAlwaysOffSampler returns a [otelsdktrace.Sampler] with always off sampling.
func NewAlwaysOffSampler() otelsdktrace.Sampler {
	return otelsdktrace.NeverSample()
}

// NewTraceIdRatioSampler returns a [otelsdktrace.Sampler] with trace id ratio sampling.
func NewTraceIdRatioSampler(ratio float64) otelsdktrace.Sampler {
	return otelsdktrace.TraceIDRatioBased(ratio)
}
