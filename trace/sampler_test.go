package trace_test

import (
	"testing"

	"github.com/ankorstore/yokai/trace"
	"github.com/stretchr/testify/assert"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestNewParentBasedAlwaysOnSampler(t *testing.T) {
	t.Parallel()

	sampler := trace.NewParentBasedAlwaysOnSampler()
	assert.Equal(t, otelsdktrace.ParentBased(otelsdktrace.AlwaysSample()), sampler)
}

func TestNewParentBasedAlwaysOffSampler(t *testing.T) {
	t.Parallel()

	sampler := trace.NewParentBasedAlwaysOffSampler()
	assert.Equal(t, otelsdktrace.ParentBased(otelsdktrace.NeverSample()), sampler)
}

func TestNewParentBasedTraceIdRatioSampler(t *testing.T) {
	t.Parallel()

	sampler := trace.NewParentBasedTraceIdRatioSampler(0.5)
	assert.Equal(t, otelsdktrace.ParentBased(otelsdktrace.TraceIDRatioBased(0.5)), sampler)
}

func TestNewAlwaysOnSampler(t *testing.T) {
	t.Parallel()

	sampler := trace.NewAlwaysOnSampler()
	assert.Equal(t, otelsdktrace.AlwaysSample(), sampler)
}

func TestNewAlwaysOffSampler(t *testing.T) {
	t.Parallel()

	sampler := trace.NewAlwaysOffSampler()
	assert.Equal(t, otelsdktrace.NeverSample(), sampler)
}

func TestNewTraceIdRatioSampler(t *testing.T) {
	t.Parallel()

	sampler := trace.NewTraceIdRatioSampler(0.5)
	assert.Equal(t, otelsdktrace.TraceIDRatioBased(0.5), sampler)
}
