package trace_test

import (
	"testing"

	"github.com/ankorstore/yokai/trace"
	"github.com/stretchr/testify/assert"
)

func TestSpanProcessorAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		processor trace.SpanProcessor
		expected  string
	}{
		{trace.StdoutSpanProcessor, trace.Stdout},
		{trace.TestSpanProcessor, trace.Test},
		{trace.OtlpGrpcSpanProcessor, trace.OtlpGrpc},
		{trace.NoopSpanProcessor, trace.Noop},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.processor.String())
	}
}

func TestFetchSpanProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected trace.SpanProcessor
	}{
		{trace.Stdout, trace.StdoutSpanProcessor},
		{trace.Test, trace.TestSpanProcessor},
		{trace.OtlpGrpc, trace.OtlpGrpcSpanProcessor},
		{"default", trace.NoopSpanProcessor},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, trace.FetchSpanProcessor(tt.input))
	}
}

func TestSamplerAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		sampler  trace.Sampler
		expected string
	}{
		{trace.ParentBasedAlwaysOffSampler, trace.ParentBasedAlwaysOff},
		{trace.ParentBasedTraceIdRatioSampler, trace.ParentBasedTraceIdRatio},
		{trace.AlwaysOnSampler, trace.AlwaysOn},
		{trace.AlwaysOffSampler, trace.AlwaysOff},
		{trace.TraceIdRatioSampler, trace.TraceIdRatio},
		{trace.ParentBasedAlwaysOnSampler, trace.ParentBasedAlwaysOn},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.sampler.String())
	}
}

func TestFetchSampler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected trace.Sampler
	}{
		{trace.ParentBasedAlwaysOff, trace.ParentBasedAlwaysOffSampler},
		{trace.ParentBasedTraceIdRatio, trace.ParentBasedTraceIdRatioSampler},
		{trace.AlwaysOn, trace.AlwaysOnSampler},
		{trace.AlwaysOff, trace.AlwaysOffSampler},
		{trace.TraceIdRatio, trace.TraceIdRatioSampler},
		{"default", trace.ParentBasedAlwaysOnSampler},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, trace.FetchSampler(tt.input))
	}
}
