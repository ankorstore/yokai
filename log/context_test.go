package log_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

const (
	testTraceId = "c4ca71e03e42c2c3d54293a6e2608bfa"
	testSpanId  = "8d0fdc8a74baaaea"
)

func TestCtxLogger(t *testing.T) {
	t.Parallel()

	testLogBuffer := logtest.NewDefaultTestLogBuffer()

	ctx := context.Background()

	zeroLogger := zerolog.New(testLogBuffer).With().Str("test", "some value").Logger()
	ctx = zeroLogger.WithContext(ctx)

	logger := log.CtxLogger(ctx)
	assert.NotNil(t, logger)

	logger.Info().Msg("some message")

	hasRecord, err := testLogBuffer.HasRecord(map[string]interface{}{
		"level":   "info",
		"test":    "some value",
		"message": "some message",
	})
	assert.NoError(t, err)
	assert.True(t, hasRecord)

	containRecord, err := testLogBuffer.ContainRecord(map[string]interface{}{
		"level":   "info",
		"test":    "some value",
		"message": "ome mess",
	})
	assert.NoError(t, err)
	assert.True(t, containRecord)
}

func TestCtxLoggerWithSpanContext(t *testing.T) {
	t.Parallel()

	testLogBuffer := logtest.NewDefaultTestLogBuffer()

	traceId, err := trace.TraceIDFromHex(testTraceId)
	assert.NoError(t, err)

	spanId, err := trace.SpanIDFromHex(testSpanId)
	assert.NoError(t, err)

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
		SpanID:  spanId,
	})

	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	zeroLogger := zerolog.New(testLogBuffer).With().Str("test", "some value").Logger()
	ctx = zeroLogger.WithContext(ctx)

	logger := log.CtxLogger(ctx)
	assert.NotNil(t, logger)

	logger.Info().Msg("some message")

	hasRecord, err := testLogBuffer.HasRecord(map[string]interface{}{
		"level":   "info",
		"test":    "some value",
		"message": "some message",
		"traceID": testTraceId,
		"spanID":  testSpanId,
	})
	assert.NoError(t, err)
	assert.True(t, hasRecord)

	containRecord, err := testLogBuffer.ContainRecord(map[string]interface{}{
		"level":   "info",
		"test":    "some value",
		"message": "ome mess",
		"traceID": testTraceId,
		"spanID":  testSpanId,
	})
	assert.NoError(t, err)
	assert.True(t, containRecord)
}
