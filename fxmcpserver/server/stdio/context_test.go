package stdio_test

import (
	"context"
	"testing"

	servercontext "github.com/ankorstore/yokai/fxmcpserver/server/context"
	"github.com/ankorstore/yokai/fxmcpserver/server/stdio"
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

func TestDefaultMCPStdioServerContextHandler_Handle(t *testing.T) {
	t.Parallel()

	gm := new(generatorMock)
	gm.On("Generate").Return("test-request-id")

	tp := trace.NewTracerProvider()

	lb := logtest.NewDefaultTestLogBuffer()
	lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
	assert.NoError(t, err)

	handler := stdio.NewDefaultMCPStdioServerContextHandler(gm, tp, lg)

	ctx := handler.Handle()(context.Background())

	assert.Equal(t, "test-request-id", servercontext.CtxRequestId(ctx))

	span, ok := servercontext.CtxRootSpan(ctx).(trace.ReadWriteSpan)
	assert.True(t, ok)

	assert.Equal(t, "MCP", span.Name())

	for _, attr := range span.Attributes() {
		if attr.Key == "system" {
			assert.Equal(t, "mcpserver", attr.Value.AsString())
		}
		if attr.Key == "mcp.transport" {
			assert.Equal(t, "stdio", attr.Value.AsString())
		}
		if attr.Key == "mcp.requestID" {
			assert.Equal(t, "test-request-id", attr.Value.AsString())
		}
	}

	log.CtxLogger(ctx).Info().Msg("test log")

	logtest.AssertHasLogRecord(t, lb, map[string]any{
		"level":        "info",
		"system":       "mcpserver",
		"mcpTransport": "stdio",
		"mcpRequestID": "test-request-id",
		"message":      "test log",
	})

	gm.AssertExpectations(t)
}
