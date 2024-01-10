package log

import (
	"context"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// CtxLogger retrieves a [Logger] from a provided context (or creates and appends a new one if missing).
//
// It automatically adds the traceID and spanID log fields depending on current tracing context.
func CtxLogger(ctx context.Context) *Logger {
	fields := make(map[string]interface{})

	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		fields["traceID"] = spanContext.TraceID().String()
	}
	if spanContext.HasSpanID() {
		fields["spanID"] = spanContext.SpanID().String()
	}

	if len(fields) > 0 {
		logger := zerolog.Ctx(ctx).With().Fields(fields).Logger()

		return &Logger{&logger}
	}

	return &Logger{zerolog.Ctx(ctx)}
}
