package sse

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

var _ MCPSSEServerContextHandler = (*DefaultMCPSSEServerContextHandler)(nil)

// MCPSSEServerContextHook is the interface for MCP SSE server context hooks.
type MCPSSEServerContextHook interface {
	Handle() server.SSEContextFunc
}

// MCPSSEServerContextHandler is the interface for MCP SSE server context handlers.
type MCPSSEServerContextHandler interface {
	Handle() server.SSEContextFunc
}

// DefaultMCPSSEServerContextHandler is the default MCPSSEServerContextHandler implementation.
type DefaultMCPSSEServerContextHandler struct {
	generator         uuid.UuidGenerator
	tracerProvider    ot.TracerProvider
	textMapPropagator propagation.TextMapPropagator
	logger            *log.Logger
	contextHooks      []MCPSSEServerContextHook
}

// NewDefaultMCPSSEServerContextHandler returns a new DefaultMCPSSEServerContextHandler instance.
func NewDefaultMCPSSEServerContextHandler(
	generator uuid.UuidGenerator,
	tracerProvider ot.TracerProvider,
	textMapPropagator propagation.TextMapPropagator,
	logger *log.Logger,
	contextHooks ...MCPSSEServerContextHook,
) *DefaultMCPSSEServerContextHandler {
	return &DefaultMCPSSEServerContextHandler{
		generator:         generator,
		tracerProvider:    tracerProvider,
		textMapPropagator: textMapPropagator,
		logger:            logger,
		contextHooks:      contextHooks,
	}
}

// Handle returns the handler func.
func (h *DefaultMCPSSEServerContextHandler) Handle() server.SSEContextFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		// start time propagation
		ctx = fsc.WithStartTime(ctx, time.Now())

		// sessionId propagation
		sID := req.URL.Query().Get("sessionId")

		ctx = fsc.WithSessionID(ctx, sID)

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
				attribute.String("mcp.transport", "sse"),
				attribute.String("mcp.sessionID", sID),
				attribute.String("mcp.requestID", rID),
			),
		)

		ctx = fsc.WithRootSpan(ctx, span)

		// logger propagation
		logger := h.logger.
			With().
			Str("system", "mcpserver").
			Str("mcpTransport", "sse").
			Str("mcpSessionID", sID).
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
