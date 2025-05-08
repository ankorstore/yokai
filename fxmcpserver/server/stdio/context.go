package stdio

import (
	"context"
	"time"

	fsc "github.com/ankorstore/yokai/fxmcpserver/server/context"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
	ot "go.opentelemetry.io/otel/trace"
)

var _ MCPStdioServerContextHandler = (*DefaultMCPStdioServerContextHandler)(nil)

// MCPStdioServerContextHandler is the interface for MCP Stdio server context handlers.
type MCPStdioServerContextHandler interface {
	Handle() server.StdioContextFunc
}

// DefaultMCPStdioServerContextHandler is the default MCPStdioServerContextHandler implementation.
type DefaultMCPStdioServerContextHandler struct {
	generator      uuid.UuidGenerator
	tracerProvider ot.TracerProvider
	logger         *log.Logger
}

// NewDefaultMCPStdioServerContextHandler returns a new DefaultMCPStdioServerContextHandler instance.
func NewDefaultMCPStdioServerContextHandler(
	generator uuid.UuidGenerator,
	tracerProvider ot.TracerProvider,
	logger *log.Logger,
) *DefaultMCPStdioServerContextHandler {
	return &DefaultMCPStdioServerContextHandler{
		generator:      generator,
		tracerProvider: tracerProvider,
		logger:         logger,
	}
}

// Handle returns the handler func.
func (h *DefaultMCPStdioServerContextHandler) Handle() server.StdioContextFunc {
	return func(ctx context.Context) context.Context {
		// start time propagation
		ctx = fsc.WithStartTime(ctx, time.Now())

		// requestId propagation
		rID := h.generator.Generate()

		ctx = fsc.WithRequestID(ctx, rID)

		// tracer propagation
		ctx = trace.WithContext(ctx, h.tracerProvider)

		ctx, span := trace.CtxTracer(ctx).Start(
			ctx,
			"MCP",
			ot.WithNewRoot(),
			ot.WithSpanKind(ot.SpanKindServer),
			ot.WithAttributes(
				attribute.String("system", "mcpserver"),
				attribute.String("mcp.transport", "stdio"),
				attribute.String("mcp.requestID", rID),
			),
		)

		ctx = fsc.WithRootSpan(ctx, span)

		// logger propagation
		logger := h.logger.
			With().
			Str("system", "mcpserver").
			Str("mcpTransport", "stdio").
			Str("mcpRequestID", rID).
			Logger()

		ctx = logger.WithContext(ctx)

		// cancellation removal propagation
		return context.WithoutCancel(ctx)
	}
}
