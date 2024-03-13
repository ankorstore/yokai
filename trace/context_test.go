package trace_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestCtxTracerProviderWithDefaultGlobalTracerProvider(t *testing.T) {
	t.Parallel()

	assert.Equal(t, otel.GetTracerProvider(), trace.CtxTracerProvider(context.Background()))
}

func TestCtxTracerProviderWithCustomTracerProvider(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	ctx := trace.WithContext(context.Background(), tracerProvider)
	assert.Equal(t, tracerProvider, trace.CtxTracerProvider(ctx))
}

func TestCtxTracer(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	ctx := trace.WithContext(context.Background(), tracerProvider)
	assert.Implements(t, (*oteltrace.Tracer)(nil), trace.CtxTracer(ctx))
}
