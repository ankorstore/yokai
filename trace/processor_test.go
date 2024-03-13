package trace_test

import (
	"context"
	"testing"
	"time"

	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

func TestNewTestSpanProcessor(t *testing.T) {
	t.Parallel()

	testExporter := tracetest.NewDefaultTestTraceExporter()
	spanProcessor := trace.NewTestSpanProcessor(testExporter)

	assert.Implements(t, (*otelsdktrace.SpanProcessor)(nil), spanProcessor)
}

func TestNewNoopSpanProcessor(t *testing.T) {
	t.Parallel()

	spanProcessor := trace.NewNoopSpanProcessor()

	assert.Implements(t, (*otelsdktrace.SpanProcessor)(nil), spanProcessor)
}

func TestNewStdoutSpanProcessor(t *testing.T) {
	t.Parallel()

	spanProcessor := trace.NewStdoutSpanProcessor()

	assert.Implements(t, (*otelsdktrace.SpanProcessor)(nil), spanProcessor)
}

func TestNewOtlpGrpcSpanProcessorSuccess(t *testing.T) {
	t.Parallel()

	spanProcessor, err := trace.NewOtlpGrpcSpanProcessor(context.Background(), &grpc.ClientConn{})

	assert.NoError(t, err)
	assert.Implements(t, (*otelsdktrace.SpanProcessor)(nil), spanProcessor)
}

func TestNewOtlpGrpcSpanProcessorFailure(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Microsecond))
	defer cancel()

	conn, _ := trace.NewOtlpGrpcClientConnection(ctx, "https://example.com")

	spanProcessor, err := trace.NewOtlpGrpcSpanProcessor(ctx, conn)

	assert.NoError(t, err)
	assert.Implements(t, (*otelsdktrace.SpanProcessor)(nil), spanProcessor)
}
