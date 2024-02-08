package fxcron_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const testName = "some test name"
const testExecutionId = "some test execution id"

func TestCtxCronJobName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	assert.Equal(t, "", fxcron.CtxCronJobName(ctx))

	ctx = context.WithValue(context.Background(), fxcron.CtxCronJobNameKey{}, testName)

	assert.Equal(t, testName, fxcron.CtxCronJobName(ctx))
}

func TestCtxCronJobExecutionId(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	assert.Equal(t, "", fxcron.CtxCronJobExecutionId(ctx))

	ctx = context.WithValue(ctx, fxcron.CtxCronJobExecutionIdKey{}, testExecutionId)

	assert.Equal(t, testExecutionId, fxcron.CtxCronJobExecutionId(ctx))
}

func TestCtxLogger(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	logger := log.CtxLogger(ctx)

	assert.Equal(t, logger, fxcron.CtxLogger(ctx))
}

func TestCtxTracer(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
		trace.WithSpanProcessor(fxcron.NewTracerProviderCronJobAnnotator()),
	)
	assert.NoError(t, err)

	ctx := context.WithValue(context.Background(), fxcron.CtxCronJobNameKey{}, testName)
	ctx = context.WithValue(ctx, fxcron.CtxCronJobExecutionIdKey{}, testExecutionId)
	ctx = context.WithValue(ctx, trace.CtxKey{}, tracerProvider)

	_, span := fxcron.CtxTracer(ctx).Start(
		ctx,
		"some span",
		oteltrace.WithAttributes(attribute.String("some attribute", "some value")),
	)
	span.End()

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"some span",
		attribute.String("CronJob", testName),
		attribute.String("CronJobExecutionID", testExecutionId),
		attribute.String("some attribute", "some value"),
	)
}
