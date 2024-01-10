package tracetest_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestAssertHasTraceSpan(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	_, span := tracerProvider.Tracer("test").Start(
		context.Background(),
		"test span",
		oteltrace.WithAttributes(
			attribute.String("string attribute name", "string attribute value"),
			attribute.Int("int attribute name", 42),
		),
	)
	span.End()

	mt := new(testing.T)
	tracetest.AssertHasTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.Int("int attribute name", 42),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("invalid attribute name", "invalid attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 24),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 42),
		attribute.String("invalid attribute name", "invalid attribute value"),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasTraceSpan(
		mt,
		exporter,
		"invalid span",
	)
	assert.True(t, mt.Failed())
}

func TestAssertHasNotTraceSpan(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	_, span := tracerProvider.Tracer("test").Start(
		context.Background(),
		"test span",
		oteltrace.WithAttributes(
			attribute.String("string attribute name", "string attribute value"),
			attribute.Int("int attribute name", 42),
		),
	)
	span.End()

	mt := new(testing.T)
	tracetest.AssertHasNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.Int("int attribute name", 42),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("invalid attribute name", "invalid attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 24),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 42),
		attribute.String("invalid attribute name", "invalid attribute value"),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertHasNotTraceSpan(
		mt,
		exporter,
		"invalid span",
	)
	assert.False(t, mt.Failed())
}

func TestAssertContainTraceSpan(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	_, span := tracerProvider.Tracer("test").Start(
		context.Background(),
		"test span",
		oteltrace.WithAttributes(
			attribute.String("string attribute name", "string attribute value"),
			attribute.Int("int attribute name", 42),
		),
	)
	span.End()

	mt := new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "attribute value"),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.Int("int attribute name", 42),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("invalid attribute name", "invalid attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 24),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 42),
		attribute.String("invalid attribute name", "invalid attribute value"),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainTraceSpan(
		mt,
		exporter,
		"invalid span",
	)
	assert.True(t, mt.Failed())
}

func TestAssertContainNotTraceSpan(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	_, span := tracerProvider.Tracer("test").Start(
		context.Background(),
		"test span",
		oteltrace.WithAttributes(
			attribute.String("string attribute name", "string attribute value"),
			attribute.Int("int attribute name", 42),
		),
	)
	span.End()

	mt := new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "attribute value"),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.Int("int attribute name", 42),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.True(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("invalid attribute name", "invalid attribute value"),
		attribute.Int("int attribute name", 42),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 24),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"test span",
		attribute.String("string attribute name", "string attribute value"),
		attribute.Int("int attribute name", 42),
		attribute.String("invalid attribute name", "invalid attribute value"),
	)
	assert.False(t, mt.Failed())

	mt = new(testing.T)
	tracetest.AssertContainNotTraceSpan(
		mt,
		exporter,
		"invalid span",
	)
	assert.False(t, mt.Failed())
}
