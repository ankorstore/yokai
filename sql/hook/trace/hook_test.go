package trace_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/sql/hook"
	"github.com/ankorstore/yokai/sql/hook/trace"
	yokaitrace "github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestTraceHookWithDefaults(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	tp, err := yokaitrace.NewDefaultTracerProviderFactory().Create(
		yokaitrace.WithSpanProcessor(yokaitrace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	h := trace.NewTraceHook()

	ctx := yokaitrace.WithContext(context.Background(), tp)
	event := hook.NewHookEvent("system", "operation", "query", "argument")

	ctx = h.Before(ctx, event)

	event.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		Stop()

	h.After(ctx, event)

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"operation",
		semconv.DBSystemKey.String("system"),
		semconv.DBStatementKey.String("query"),
		attribute.Int64("db.lastInsertId", int64(1)),
		attribute.Int64("db.rowsAffected", int64(2)),
	)
}

func TestTraceHookWithOptions(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	tp, err := yokaitrace.NewDefaultTracerProviderFactory().Create(
		yokaitrace.WithSpanProcessor(yokaitrace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	h := trace.NewTraceHook(
		trace.WithArguments(true),
		trace.WithExcludedOperations("excludedOperation"),
	)

	ctx := yokaitrace.WithContext(context.Background(), tp)

	// regular event
	event := hook.NewHookEvent("system", "regularOperation", "query", "argument")

	ctx = h.Before(ctx, event)

	event.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		Stop()

	h.After(ctx, event)

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"regularOperation",
		semconv.DBSystemKey.String("system"),
		semconv.DBStatementKey.String("query"),
		attribute.String("db.statement.arguments", fmt.Sprintf("%+v", "argument")),
		attribute.Int64("db.lastInsertId", int64(1)),
		attribute.Int64("db.rowsAffected", int64(2)),
	)

	// excluded operation event
	excludedOperationEvent := hook.NewHookEvent("system", "excludedOperation", "query", "argument")

	ctx = h.Before(ctx, excludedOperationEvent)

	excludedOperationEvent.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		Stop()

	h.After(ctx, excludedOperationEvent)

	tracetest.AssertHasNotTraceSpan(
		t,
		exporter,
		"excludedOperation",
		semconv.DBSystemKey.String("system"),
		semconv.DBStatementKey.String("query"),
		attribute.String("db.statement.arguments", fmt.Sprintf("%+v", "argument")),
		attribute.Int64("db.lastInsertId", int64(1)),
		attribute.Int64("db.rowsAffected", int64(2)),
	)

	// error event
	errorEvent := hook.NewHookEvent("system", "errorOperation", "query", "argument")

	ctx = h.Before(ctx, errorEvent)

	errorEvent.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		SetError(fmt.Errorf("test error")).
		Stop()

	h.After(ctx, errorEvent)

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"errorOperation",
		semconv.DBSystemKey.String("system"),
		semconv.DBStatementKey.String("query"),
		attribute.String("db.statement.arguments", fmt.Sprintf("%+v", "argument")),
		attribute.Int64("db.lastInsertId", int64(1)),
		attribute.Int64("db.rowsAffected", int64(2)),
	)

	span, err := exporter.Span("errorOperation")
	assert.NoError(t, err)

	assert.Equal(t, codes.Error, span.Status.Code)
	assert.Len(t, span.Events, 1)
	assert.Equal(
		t, []attribute.KeyValue{
			semconv.ExceptionType("*errors.errorString"),
			semconv.ExceptionMessage("test error"),
		},
		span.Events[0].Attributes,
	)
}

func TestTraceHookAfterWithNonRecordingSpan(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	tp, err := yokaitrace.NewDefaultTracerProviderFactory().Create(
		yokaitrace.WithSpanProcessor(yokaitrace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	h := trace.NewTraceHook()

	ctx := yokaitrace.WithContext(context.Background(), tp)
	event := hook.NewHookEvent("system", "operation", "query", "argument")

	ctx = h.Before(ctx, event)

	span := oteltrace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("attribute.name", "value"))
	span.End()

	h.After(ctx, event)

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"operation",
		semconv.DBSystemKey.String("system"),
		semconv.DBStatementKey.String("query"),
		attribute.String("attribute.name", "value"),
	)
}
