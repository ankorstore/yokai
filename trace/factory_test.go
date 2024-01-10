package trace_test

import (
	"testing"

	"github.com/ankorstore/yokai/trace"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	testResource      = resource.Empty()
	testSampler       = otelsdktrace.AlwaysSample()
	testSpanProcessor = trace.NewNoopSpanProcessor()
)

func TestDefaultTracerProviderFactory(t *testing.T) {
	t.Parallel()

	factory := trace.NewDefaultTracerProviderFactory()

	assert.IsType(t, &trace.DefaultTracerProviderFactory{}, factory)
	assert.Implements(t, (*trace.TracerProviderFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	factory := trace.NewDefaultTracerProviderFactory()

	tracerProvider, err := factory.Create(
		trace.WithResource(testResource),
		trace.WithSampler(testSampler),
		trace.WithSpanProcessor(testSpanProcessor),
	)

	assert.NoError(t, err)
	assert.IsType(t, &otelsdktrace.TracerProvider{}, tracerProvider)
	assert.Equal(t, tracerProvider, otel.GetTracerProvider())
}

func TestCreateAsNotGlobal(t *testing.T) {
	factory := trace.NewDefaultTracerProviderFactory()

	tracerProvider, err := factory.Create(
		trace.Global(false),
		trace.WithResource(testResource),
		trace.WithSampler(testSampler),
		trace.WithSpanProcessor(testSpanProcessor),
	)

	assert.NoError(t, err)
	assert.IsType(t, &otelsdktrace.TracerProvider{}, tracerProvider)
	assert.NotEqual(t, tracerProvider, otel.GetTracerProvider())
}

func TestCreateWithoutSpanProcessor(t *testing.T) {
	factory := trace.NewDefaultTracerProviderFactory()

	tracerProvider, err := factory.Create(
		trace.WithResource(testResource),
		trace.WithSampler(testSampler),
	)

	assert.NoError(t, err)
	assert.IsType(t, &otelsdktrace.TracerProvider{}, tracerProvider)
}
