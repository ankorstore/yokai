package trace

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TraceHook is a hook.Hook implementation for SQL tracing.
type TraceHook struct {
	options Options
}

// NewTraceHook returns a new TraceHook, for a provided list of TraceHookOption.
func NewTraceHook(options ...TraceHookOption) *TraceHook {
	appliedOpts := DefaultTraceHookOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	return &TraceHook{
		options: appliedOpts,
	}
}

// Before executes SQL tracing logic before SQL operations.
func (h *TraceHook) Before(ctx context.Context, event *sql.HookEvent) context.Context {
	if sql.ContainsOperation(h.options.ExcludedOperations, event.Operation()) {
		return ctx
	}

	attributes := []attribute.KeyValue{
		semconv.DBSystemKey.String(event.System().String()),
	}

	if query := event.Query(); query != "" {
		attributes = append(
			attributes,
			semconv.DBStatementKey.String(query),
		)
	}

	if h.options.Arguments && event.Arguments() != nil {
		attributes = append(
			attributes,
			attribute.String("db.statement.arguments", fmt.Sprintf("%+v", event.Arguments())),
		)
	}

	ctx, _ = trace.CtxTracerProvider(ctx).Tracer("yokai-sql").Start(
		ctx,
		fmt.Sprintf("SQL %s", event.Operation().String()),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(attributes...),
	)

	return ctx
}

// After executes SQL tracing logic after SQL operations.
func (h *TraceHook) After(ctx context.Context, event *sql.HookEvent) {
	if sql.ContainsOperation(h.options.ExcludedOperations, event.Operation()) {
		return
	}

	span := oteltrace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	code := codes.Ok
	if event.Error() != nil {
		if !errors.Is(event.Error(), driver.ErrSkip) {
			code = codes.Error
			span.RecordError(event.Error())
		}
	}
	span.SetStatus(code, code.String())

	attributes := []attribute.KeyValue{
		attribute.Int64("db.lastInsertId", event.LastInsertId()),
		attribute.Int64("db.rowsAffected", event.RowsAffected()),
	}

	latency, err := event.Latency()
	if err == nil {
		attributes = append(attributes, attribute.String("db.latency", latency.String()))
	}

	span.SetAttributes(attributes...)

	span.End()
}
