package sse_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	servercontext "github.com/ankorstore/yokai/fxmcpserver/server/context"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/sdk/trace"
)

type generatorMock struct {
	mock.Mock
}

func (m *generatorMock) Generate() string {
	return m.Called().String(0)
}

//nolint:cyclop
func TestDefaultMCPSSEServerContextHandler_Handle(t *testing.T) {
	t.Parallel()

	t.Run("with provided session id and request id", func(t *testing.T) {
		t.Parallel()

		gm := new(generatorMock)
		gm.AssertNotCalled(t, "Generate")

		tp := trace.NewTracerProvider()

		lb := logtest.NewDefaultTestLogBuffer()
		lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
		assert.NoError(t, err)

		handler := sse.NewDefaultMCPSSEServerContextHandler(gm, tp, lg)

		req := httptest.NewRequest(http.MethodGet, "/sse?sessionId=test-session-id", nil)
		req.Header.Set("X-Request-Id", "test-request-id")

		ctx := handler.Handle()(context.Background(), req)

		assert.Equal(t, "test-session-id", servercontext.CtxSessionID(ctx))
		assert.Equal(t, "test-request-id", servercontext.CtxRequestId(ctx))

		span, ok := servercontext.CtxRootSpan(ctx).(trace.ReadWriteSpan)
		assert.True(t, ok)

		assert.Equal(t, "MCP", span.Name())

		for _, attr := range span.Attributes() {
			if attr.Key == "system" {
				assert.Equal(t, "mcpserver", attr.Value.AsString())
			}
			if attr.Key == "mcp.transport" {
				assert.Equal(t, "sse", attr.Value.AsString())
			}
			if attr.Key == "mcp.sessionID" {
				assert.Equal(t, "test-session-id", attr.Value.AsString())
			}
			if attr.Key == "mcp.requestID" {
				assert.Equal(t, "test-request-id", attr.Value.AsString())
			}
		}

		log.CtxLogger(ctx).Info().Msg("test log")

		logtest.AssertHasLogRecord(t, lb, map[string]any{
			"level":        "info",
			"system":       "mcpserver",
			"mcpTransport": "sse",
			"mcpSessionID": "test-session-id",
			"mcpRequestID": "test-request-id",
			"message":      "test log",
		})

		gm.AssertExpectations(t)
	})

	t.Run("without provided session id and request id", func(t *testing.T) {
		t.Parallel()

		gm := new(generatorMock)
		gm.On("Generate").Return("test-request-id")

		tp := trace.NewTracerProvider()

		lb := logtest.NewDefaultTestLogBuffer()
		lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
		assert.NoError(t, err)

		handler := sse.NewDefaultMCPSSEServerContextHandler(gm, tp, lg)

		req := httptest.NewRequest(http.MethodGet, "/sse", nil)

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
				assert.Equal(t, "sse", attr.Value.AsString())
			}
			if attr.Key == "mcp.sessionID" {
				assert.Equal(t, "", attr.Value.AsString())
			}
			if attr.Key == "mcp.requestID" {
				assert.Equal(t, "test-request-id", attr.Value.AsString())
			}
		}

		log.CtxLogger(ctx).Info().Msg("test log")

		logtest.AssertHasLogRecord(t, lb, map[string]any{
			"level":        "info",
			"system":       "mcpserver",
			"mcpTransport": "sse",
			"mcpSessionID": "",
			"mcpRequestID": "test-request-id",
			"message":      "test log",
		})

		gm.AssertExpectations(t)
	})
}
