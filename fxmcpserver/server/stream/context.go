package stream

import (
	"context"
	"net/http"
	"time"

	fsc "github.com/ankorstore/yokai/fxmcpserver/server/context"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	ot "go.opentelemetry.io/otel/trace"
)

var _ MCPStreamableHTTPServerContextHandler = (*DefaultMCPStreamableHTTPServerContextHandler)(nil)

// MCPStreamableHTTPServerContextHook is the interface for MCP StreamableHTTP server context hooks.
type MCPStreamableHTTPServerContextHook interface {
	Handle() server.HTTPContextFunc
}

// MCPStreamableHTTPServerContextHandler is the interface for MCP StreamableHTTP server context handlers.
type MCPStreamableHTTPServerContextHandler interface {
	Handle() server.HTTPContextFunc
}

// DefaultMCPStreamableHTTPServerContextHandler is the default MCPStreamableHTTPServerContextHandler implementation.
type DefaultMCPStreamableHTTPServerContextHandler struct {
	generator         uuid.UuidGenerator
	tracerProvider    ot.TracerProvider
	textMapPropagator propagation.TextMapPropagator
	logger            *log.Logger
	contextHooks      []MCPStreamableHTTPServerContextHook
}

// NewDefaultMCPStreamableHTTPServerContextHandler returns a new DefaultMCPStreamableHTTPServerContextHandler instance.
func NewDefaultMCPStreamableHTTPServerContextHandler(
	generator uuid.UuidGenerator,
	tracerProvider ot.TracerProvider,
	textMapPropagator propagation.TextMapPropagator,
	logger *log.Logger,
	contextHooks ...MCPStreamableHTTPServerContextHook,
) *DefaultMCPStreamableHTTPServerContextHandler {
	return &DefaultMCPStreamableHTTPServerContextHandler{
		generator:         generator,
		tracerProvider:    tracerProvider,
		textMapPropagator: textMapPropagator,
		logger:            logger,
		contextHooks:      contextHooks,
	}
}

// Handle returns the handler func.
func (h *DefaultMCPStreamableHTTPServerContextHandler) Handle() server.HTTPContextFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		// start time propagation
		ctx = fsc.WithStartTime(ctx, time.Now())

		// requestId propagation
		rID := req.Header.Get("X-Request-Id")

		if rID == "" {
			rID = h.generator.Generate()
			req.Header.Set("X-Request-Id", rID)
		}

		ctx = fsc.WithRequestID(ctx, rID)

		// tracer propagation
		ctx = h.textMapPropagator.Extract(ctx, propagation.HeaderCarrier(req.Header))

		ctx = trace.WithContext(ctx, h.tracerProvider)

		ctx, span := trace.CtxTracer(ctx).Start(
			ctx,
			"MCP",
			ot.WithSpanKind(ot.SpanKindServer),
			ot.WithAttributes(
				attribute.String("system", "mcpserver"),
				attribute.String("mcp.transport", "streamable-http"),
				attribute.String("mcp.requestID", rID),
			),
		)

		ctx = fsc.WithRootSpan(ctx, span)

		// logger propagation
		logger := h.logger.
			With().
			Str("system", "mcpserver").
			Str("mcpTransport", "streamable-http").
			Str("mcpRequestID", rID).
			Logger()

		ctx = logger.WithContext(ctx)

		// cancellation removal propagation
		ctx = context.WithoutCancel(ctx)

		// hooks propagation
		for _, hook := range h.contextHooks {
			ctx = hook.Handle()(ctx, req)
		}

		return ctx
	}
}
