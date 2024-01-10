package trace

import (
	"context"

	"github.com/ankorstore/yokai/trace/tracetest"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	otelsdktracetest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	"google.golang.org/grpc"
)

const (
	Stdout   = "stdout"    // processor to send trace spans to the standard output
	OtlpGrpc = "otlp-grpc" // processor to send the trace spans via OTLP/gRPC
	Test     = "test"      // processor to send the trace spans to a test buffer
	Noop     = "noop"      // processor to void the trace spans
)

// NewTestSpanProcessor returns a [OTEL SpanProcessor] using a sync [tracetest.TestTraceExporter].
//
// [OTEL SpanProcessor]: https://github.com/open-telemetry/opentelemetry-go
func NewTestSpanProcessor(testTraceExporter tracetest.TestTraceExporter) trace.SpanProcessor {
	return trace.NewSimpleSpanProcessor(testTraceExporter)
}

// NewNoopSpanProcessor returns a [OTEL SpanProcessor] that voids trace spans via an async [otelsdktracetest.NoopExporter].
//
// [OTEL SpanProcessor]: https://github.com/open-telemetry/opentelemetry-go
func NewNoopSpanProcessor() trace.SpanProcessor {
	return trace.NewBatchSpanProcessor(otelsdktracetest.NewNoopExporter())
}

// NewTestSpanProcessor returns a [OTEL SpanProcessor] using an async [stdouttrace.Exporter].
//
// [OTEL SpanProcessor]: https://github.com/open-telemetry/opentelemetry-go
func NewStdoutSpanProcessor(options ...stdouttrace.Option) trace.SpanProcessor {
	exporter, _ := stdouttrace.New(options...)

	return trace.NewBatchSpanProcessor(exporter)
}

// NewOtlpGrpcSpanProcessor returns a [OTEL SpanProcessor] using an async [otlptracegrpc.Exporter].
//
// [OTEL SpanProcessor]: https://github.com/open-telemetry/opentelemetry-go
func NewOtlpGrpcSpanProcessor(ctx context.Context, conn *grpc.ClientConn) (trace.SpanProcessor, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	return trace.NewBatchSpanProcessor(exporter), nil
}
