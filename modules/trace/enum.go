package trace

import "strings"

// SpanProcessor is an enum for the supported span processors.
type SpanProcessor int

const (
	NoopSpanProcessor SpanProcessor = iota
	StdoutSpanProcessor
	TestSpanProcessor
	OtlpGrpcSpanProcessor
)

// String returns a string representation of the [SpanProcessor].
//
//nolint:exhaustive
func (p SpanProcessor) String() string {
	switch p {
	case StdoutSpanProcessor:
		return Stdout
	case TestSpanProcessor:
		return Test
	case OtlpGrpcSpanProcessor:
		return OtlpGrpc
	default:
		return Noop
	}
}

// FetchSpanProcessor returns a [SpanProcessor] for a given value.
func FetchSpanProcessor(p string) SpanProcessor {
	switch strings.ToLower(p) {
	case Stdout:
		return StdoutSpanProcessor
	case Test:
		return TestSpanProcessor
	case OtlpGrpc:
		return OtlpGrpcSpanProcessor
	default:
		return NoopSpanProcessor
	}
}

// Sampler is an enum for the supported samplers.
type Sampler int

const (
	ParentBasedAlwaysOnSampler Sampler = iota
	ParentBasedAlwaysOffSampler
	ParentBasedTraceIdRatioSampler
	AlwaysOnSampler
	AlwaysOffSampler
	TraceIdRatioSampler
)

// String returns a string representation of the [Sampler].
//
//nolint:exhaustive
func (s Sampler) String() string {
	switch s {
	case ParentBasedAlwaysOffSampler:
		return ParentBasedAlwaysOff
	case ParentBasedTraceIdRatioSampler:
		return ParentBasedTraceIdRatio
	case AlwaysOnSampler:
		return AlwaysOn
	case AlwaysOffSampler:
		return AlwaysOff
	case TraceIdRatioSampler:
		return TraceIdRatio
	default:
		return ParentBasedAlwaysOn
	}
}

// FetchSampler returns a [Sampler] for a given value.
func FetchSampler(s string) Sampler {
	switch strings.ToLower(s) {
	case ParentBasedAlwaysOff:
		return ParentBasedAlwaysOffSampler
	case ParentBasedTraceIdRatio:
		return ParentBasedTraceIdRatioSampler
	case AlwaysOn:
		return AlwaysOnSampler
	case AlwaysOff:
		return AlwaysOffSampler
	case TraceIdRatio:
		return TraceIdRatioSampler
	default:
		return ParentBasedAlwaysOnSampler
	}
}
