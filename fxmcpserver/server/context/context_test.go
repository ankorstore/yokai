package context_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	servercontext "github.com/ankorstore/yokai/fxmcpserver/server/context"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestCtxRequestId(t *testing.T) {
	t.Parallel()

	t.Run("with existing context entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		ctx = servercontext.WithRequestID(ctx, "test-request-id")

		assert.Equal(t, "test-request-id", servercontext.CtxRequestId(ctx))
	})

	t.Run("without existing context entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		assert.Equal(t, "", servercontext.CtxRequestId(ctx))
	})
}

func TestCtxSessionId(t *testing.T) {
	t.Parallel()

	t.Run("with existing context entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		ctx = servercontext.WithSessionID(ctx, "test-session-id")

		assert.Equal(t, "test-session-id", servercontext.CtxSessionID(ctx))
	})

	t.Run("without existing context entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		assert.Equal(t, "", servercontext.CtxSessionID(ctx))
	})
}

func TestCtxRootSpan(t *testing.T) {
	t.Parallel()

	t.Run("with existing context entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		_, span := trace.NewTracerProvider().Tracer("test-tracer").Start(ctx, "test-span")

		ctx = servercontext.WithRootSpan(ctx, span)

		assert.Equal(t, "*trace.recordingSpan", fmt.Sprintf("%T", servercontext.CtxRootSpan(ctx)))
	})

	t.Run("without existing context entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		assert.Equal(t, "trace.noopSpan", fmt.Sprintf("%T", servercontext.CtxRootSpan(ctx)))
	})
}

func TestCtxStartTime(t *testing.T) {
	t.Parallel()

	startTime, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.NoError(t, err)

	t.Run("with existing context entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		ctx = servercontext.WithStartTime(ctx, startTime)

		assert.Equal(t, startTime, servercontext.CtxStartTime(ctx))
	})

	t.Run("without existing context entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		assert.NotEqual(t, startTime, servercontext.CtxStartTime(ctx))
	})
}
