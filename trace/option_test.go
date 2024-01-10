package trace_test

import (
	"testing"

	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestGlobal(t *testing.T) {
	t.Parallel()

	var options trace.Options

	trace.Global(true)(&options)
	assert.True(t, options.Global)
}

func TestWithResource(t *testing.T) {
	t.Parallel()

	var options trace.Options

	res := resource.Empty()

	trace.WithResource(res)(&options)
	assert.Equal(t, res, options.Resource)
}

func TestWithSampler(t *testing.T) {
	t.Parallel()

	var options trace.Options

	sampler := otelsdktrace.NeverSample()

	trace.WithSampler(sampler)(&options)
	assert.Equal(t, sampler, options.Sampler)
}

func TestWithSpanProcessor(t *testing.T) {
	t.Parallel()

	var options trace.Options

	p1 := trace.NewNoopSpanProcessor()
	p2 := trace.NewTestSpanProcessor(tracetest.NewDefaultTestTraceExporter())

	trace.WithSpanProcessor(p1)(&options)
	trace.WithSpanProcessor(p2)(&options)
	assert.Equal(t, []otelsdktrace.SpanProcessor{p1, p2}, options.SpanProcessors)
}
