package tracetest_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	oteltracetest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestNewDefaultTestTraceExporter(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	assert.IsType(t, &tracetest.DefaultTestTraceExporter{}, exporter)
	assert.Implements(t, (*tracetest.TestTraceExporter)(nil), exporter)

	assert.IsType(t, &oteltracetest.InMemoryExporter{}, exporter.Exporter())
}

func TestLifecycle(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	assert.Len(t, exporter.Spans().Snapshots(), 0)

	tracer := tracerProvider.Tracer("test")

	ctx, span := tracer.Start(
		context.Background(),
		"test span",
		oteltrace.WithAttributes(
			attribute.String("string attribute name", "string attribute value"),
			attribute.Int("int attribute name", 42),
		),
	)
	span.End()

	_, span2 := tracer.Start(
		ctx,
		"test span without attributes",
	)
	span2.End()

	assert.Len(t, exporter.Spans().Snapshots(), 2)

	exportedSpan, err := exporter.Span("test span")
	assert.NoError(t, err)
	assert.Equal(t, "test span", exportedSpan.Name)

	exportedSpan2, err := exporter.Span("test span without attributes")
	assert.NoError(t, err)
	assert.Equal(t, "test span without attributes", exportedSpan2.Name)

	_, err = exporter.Span("invalid")
	assert.Error(t, err)
	assert.Equal(t, "span with name invalid cannot be found", err.Error())

	exporter.Reset()
	assert.Len(t, exporter.Spans().Snapshots(), 0)
}

func TestHasSpan(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	tracer := tracerProvider.Tracer("test")

	ctx, span := tracer.Start(
		context.Background(),
		"test span",
		oteltrace.WithAttributes(
			attribute.String("string attribute name", "string attribute value"),
			attribute.Int("int attribute name", 42),
		),
	)
	span.End()

	_, span2 := tracer.Start(
		ctx,
		"test span without attributes",
	)
	span2.End()

	assert.True(
		t,
		exporter.HasSpan(
			"test span", // valid
		),
	)
	assert.True(
		t,
		exporter.HasSpan(
			"test span", // valid
			attribute.String("string attribute name", "string attribute value"), // valid
		),
	)
	assert.True(
		t,
		exporter.HasSpan(
			"test span",                             // valid
			attribute.Int("int attribute name", 42), // valid
		),
	)
	assert.True(
		t, exporter.HasSpan(
			"test span", // valid
			attribute.String("string attribute name", "string attribute value"), // valid
			attribute.Int("int attribute name", 42),                             // valid
		),
	)
	assert.True(
		t,
		exporter.HasSpan(
			"test span without attributes", // valid
		),
	)

	assert.False(
		t,
		exporter.HasSpan(
			"test span", // valid
			attribute.String("string attribute name", "invalid attribute value"), // invalid
		),
	)
	assert.False(
		t,
		exporter.HasSpan(
			"test span", // valid
			attribute.String("invalid attribute name", "string attribute value"), // invalid
		),
	)
	assert.False(
		t,
		exporter.HasSpan(
			"test span",                             // valid
			attribute.Int("int attribute name", 24), // invalid
		),
	)
	assert.False(
		t,
		exporter.HasSpan(
			"test span", // valid
			attribute.Int("invalid attribute name", 42), // invalid
		),
	)
	assert.False(
		t, exporter.HasSpan(
			"test span", // valid
			attribute.String("string attribute name", "invalid attribute value"), // invalid
			attribute.Int("int attribute name", 24),                              // invalid
		),
	)
	assert.False(
		t, exporter.HasSpan(
			"test span", // valid
			attribute.String("string attribute name", "string attribute value"),   // valid
			attribute.String("invalid attribute name", "invalid attribute value"), // invalid
			attribute.Int("int attribute name", 42),                               // valid
		),
	)
	assert.False(
		t,
		exporter.HasSpan(
			"test span without attributes",                                        // valid
			attribute.String("invalid attribute name", "invalid attribute value"), // invalid
		),
	)
}

func TestContainSpan(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	tracer := tracerProvider.Tracer("test")

	ctx, span := tracer.Start(
		context.Background(),
		"test span",
		oteltrace.WithAttributes(
			attribute.String("string attribute name", "string attribute value"),
			attribute.Int("int attribute name", 42),
		),
	)
	span.End()

	_, span2 := tracer.Start(
		ctx,
		"test span without attributes",
	)
	span2.End()

	assert.True(
		t,
		exporter.ContainSpan(
			"test span", // valid
		),
	)
	assert.True(
		t,
		exporter.ContainSpan(
			"test span", // valid
			attribute.String("string attribute name", "string attribute value"), // valid
		),
	)
	assert.True(
		t,
		exporter.ContainSpan(
			"test span", // valid
			attribute.String("string attribute name", "attribute"), // valid
		),
	)
	assert.True(
		t,
		exporter.ContainSpan(
			"test span",                             // valid
			attribute.Int("int attribute name", 42), // valid
		),
	)
	assert.True(
		t, exporter.ContainSpan(
			"test span", // valid
			attribute.String("string attribute name", "string attribute value"), // valid
			attribute.Int("int attribute name", 42),                             // valid
		),
	)
	assert.True(
		t, exporter.ContainSpan(
			"test span", // valid
			attribute.String("string attribute name", "attribute"), // valid
			attribute.Int("int attribute name", 42),                // valid
		),
	)
	assert.True(
		t,
		exporter.ContainSpan(
			"test span without attributes", // valid
		),
	)

	assert.False(
		t,
		exporter.ContainSpan(
			"test span", // valid
			attribute.String("string attribute name", "invalid attribute value"), // invalid
		),
	)
	assert.False(
		t,
		exporter.ContainSpan(
			"test span", // valid
			attribute.String("invalid attribute name", "string attribute value"), // invalid
		),
	)
	assert.False(
		t,
		exporter.ContainSpan(
			"test span",                             // valid
			attribute.Int("int attribute name", 24), // invalid
		),
	)
	assert.False(
		t,
		exporter.ContainSpan(
			"test span", // valid
			attribute.Int("invalid attribute name", 42), // invalid
		),
	)
	assert.False(
		t, exporter.ContainSpan(
			"test span", // valid
			attribute.String("string attribute name", "invalid attribute value"), // invalid
			attribute.Int("int attribute name", 24),                              // invalid
		),
	)
	assert.False(
		t, exporter.ContainSpan(
			"test span", // valid
			attribute.String("string attribute name", "attribute"),                // valid
			attribute.String("invalid attribute name", "invalid attribute value"), // invalid
			attribute.Int("int attribute name", 42),                               // valid
		),
	)
	assert.False(
		t,
		exporter.ContainSpan(
			"test span without attributes",                                        // valid
			attribute.String("invalid attribute name", "invalid attribute value"), // invalid
		),
	)
}
