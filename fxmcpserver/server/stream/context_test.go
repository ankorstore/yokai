package stream_test

import (
	"context"
	"github.com/ankorstore/yokai/fxmcpserver/server/stream"
	"net/http"
	"net/http/httptest"
	"testing"

	servercontext "github.com/ankorstore/yokai/fxmcpserver/server/context"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/hook"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type generatorMock struct {
	mock.Mock
}

func (m *generatorMock) Generate() string {
	return m.Called().String(0)
}

func TestDefaultMCPStreamableHTTPServerContextHandler_Handle(t *testing.T) {
	t.Parallel()

	t.Run("with defaults", func(t *testing.T) {
		t.Parallel()

		gm := new(generatorMock)
		gm.On("Generate").Return("test-request-id")

		tp := trace.NewTracerProvider()

		tmp := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

		lb := logtest.NewDefaultTestLogBuffer()
		lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
		assert.NoError(t, err)

		handler := stream.NewDefaultMCPStreamableHTTPServerContextHandler(gm, tp, tmp, lg)

		req := httptest.NewRequest(http.MethodGet, "/mcp", nil)

		ctx := handler.Handle()(context.Background(), req)

		assert.Equal(t, "", servercontext.CtxSessionID(ctx))
		assert.Equal(t, "test-request-id", servercontext.CtxRequestId(ctx))

		span, ok := servercontext.CtxRootSpan(ctx).(trace.ReadWriteSpan)
		assert.True(t, ok)

		assert.Equal(t, "MCP", span.Name())

		for _, attr := range span.Attributes() {
			if attr.Key == "system" {
				assert.Equal(t, "mcpserver", attr.Value.AsString())
			}
			if attr.Key == "mcp.transport" {
				assert.Equal(t, "streamable-http", attr.Value.AsString())
			}
			if attr.Key == "mcp.requestID" {
				assert.Equal(t, "test-request-id", attr.Value.AsString())
			}
		}

		log.CtxLogger(ctx).Info().Msg("test log")

		logtest.AssertHasLogRecord(t, lb, map[string]any{
			"level":        "info",
			"system":       "mcpserver",
			"mcpTransport": "streamable-http",
			"mcpRequestID": "test-request-id",
			"message":      "test log",
		})

		gm.AssertExpectations(t)
	})

	t.Run("with provided request id and hook", func(t *testing.T) {
		t.Parallel()

		gm := new(generatorMock)
		gm.AssertNotCalled(t, "Generate")

		tp := trace.NewTracerProvider()

		tmp := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

		lb := logtest.NewDefaultTestLogBuffer()
		lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
		assert.NoError(t, err)

		hk := hook.NewSimpleMCPStreamableHTTPServerContextHook()

		handler := stream.NewDefaultMCPStreamableHTTPServerContextHandler(gm, tp, tmp, lg, hk)

		req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
		req.Header.Set("X-Request-Id", "test-request-id")

		ctx := handler.Handle()(context.Background(), req)

		assert.Equal(t, "test-request-id", servercontext.CtxRequestId(ctx))

		span, ok := servercontext.CtxRootSpan(ctx).(trace.ReadWriteSpan)
		assert.True(t, ok)

		assert.Equal(t, "MCP", span.Name())

		for _, attr := range span.Attributes() {
			if attr.Key == "system" {
				assert.Equal(t, "mcpserver", attr.Value.AsString())
			}
			if attr.Key == "mcp.transport" {
				assert.Equal(t, "streamable-http", attr.Value.AsString())
			}
			if attr.Key == "mcp.requestID" {
				assert.Equal(t, "test-request-id", attr.Value.AsString())
			}
		}

		log.CtxLogger(ctx).Info().Msg("test log")

		logtest.AssertHasLogRecord(t, lb, map[string]any{
			"level":        "info",
			"system":       "mcpserver",
			"mcpTransport": "streamable-http",
			"mcpRequestID": "test-request-id",
			"message":      "test log",
		})

		//nolint:forcetypeassert
		assert.Equal(t, "bar", ctx.Value("foo").(string))

		gm.AssertExpectations(t)
	})
}
