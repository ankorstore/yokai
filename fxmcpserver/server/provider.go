package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ankorstore/yokai/config"
	fsc "github.com/ankorstore/yokai/fxmcpserver/server/context"
	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var _ MCPServerHooksProvider = (*DefaultMCPServerHooksProvider)(nil)

// MCPServerHooksProvider is the interface for the MCP server hooks provider.
type MCPServerHooksProvider interface {
	Provide() *server.Hooks
}

// DefaultMCPServerHooksProvider is the default MCPServerHooksProvider implementation.
type DefaultMCPServerHooksProvider struct {
	config           *config.Config
	requestsCounter  *prometheus.CounterVec
	requestsDuration *prometheus.HistogramVec
}

// NewDefaultMCPServerHooksProvider returns a new DefaultMCPServerHooksProvider instance.
func NewDefaultMCPServerHooksProvider(registry prometheus.Registerer, config *config.Config) *DefaultMCPServerHooksProvider {
	namespace := Sanitize(config.GetString("modules.mcp.server.metrics.collect.namespace"))
	subsystem := Sanitize(config.GetString("modules.mcp.server.metrics.collect.subsystem"))

	buckets := prometheus.DefBuckets
	if bucketsConfig := config.GetString("modules.mcp.server.metrics.buckets"); bucketsConfig != "" {
		buckets = []float64{}

		for _, s := range Split(bucketsConfig) {
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				buckets = append(buckets, f)
			}
		}
	}

	requestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "mcp_server_requests_total",
			Help:      "Number of processed MCP requests",
		},
		[]string{
			"method",
			"target",
			"status",
		},
	)

	requestsDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "mcp_server_requests_duration_seconds",
			Help:      "Time spent processing MCP requests",
			Buckets:   buckets,
		},
		[]string{
			"method",
			"target",
		},
	)

	registry.MustRegister(requestsCounter, requestsDuration)

	return &DefaultMCPServerHooksProvider{
		config:           config,
		requestsCounter:  requestsCounter,
		requestsDuration: requestsDuration,
	}
}

// Provide provides the MCP server hooks.
//
//nolint:cyclop,gocognit
func (p *DefaultMCPServerHooksProvider) Provide() *server.Hooks {
	hooks := &server.Hooks{}

	traceRequest := p.config.GetBool("modules.mcp.server.trace.request")
	traceResponse := p.config.GetBool("modules.mcp.server.trace.response")

	logRequest := p.config.GetBool("modules.mcp.server.log.request")
	logResponse := p.config.GetBool("modules.mcp.server.log.response")

	metricsEnabled := p.config.GetBool("modules.mcp.server.metrics.collect.enabled")

	hooks.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		log.CtxLogger(ctx).Info().Str("mcpSessionID", session.SessionID()).Msg("MCP session registered")
	})

	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		latency := time.Since(fsc.CtxStartTime(ctx))

		mcpMethod := string(method)

		spanNameSuffix := mcpMethod

		spanAttributes := []attribute.KeyValue{
			attribute.String("mcp.latency", latency.String()),
			attribute.String("mcp.method", mcpMethod),
		}

		logFields := map[string]any{
			"mcpLatency": latency.String(),
			"mcpMethod":  mcpMethod,
		}

		metricTarget := ""

		jsonMessage, err := json.Marshal(message)
		if err == nil {
			if traceRequest {
				spanAttributes = append(spanAttributes, attribute.String("mcp.request", string(jsonMessage)))
			}

			if logRequest {
				logFields["mcpRequest"] = string(jsonMessage)
			}
		}

		jsonResult, err := json.Marshal(result)
		if err == nil {
			if traceResponse {
				spanAttributes = append(spanAttributes, attribute.String("mcp.response", string(jsonResult)))
			}

			if logResponse {
				logFields["mcpResponse"] = string(jsonResult)
			}
		}

		//nolint:exhaustive
		switch method {
		case mcp.MethodResourcesRead:
			if req, ok := message.(*mcp.ReadResourceRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.URI)
				spanAttributes = append(spanAttributes, attribute.String("mcp.resourceURI", req.Params.URI))
				logFields["mcpResourceURI"] = req.Params.URI
				metricTarget = req.Params.URI
			}
		case mcp.MethodPromptsGet:
			if req, ok := message.(*mcp.GetPromptRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.Name)
				spanAttributes = append(spanAttributes, attribute.String("mcp.prompt", req.Params.Name))
				logFields["mcpPrompt"] = req.Params.Name
				metricTarget = req.Params.Name
			}
		case mcp.MethodToolsCall:
			if req, ok := message.(*mcp.CallToolRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.Name)
				spanAttributes = append(spanAttributes, attribute.String("mcp.tool", req.Params.Name))
				logFields["mcpTool"] = req.Params.Name
				metricTarget = req.Params.Name
			}
		}

		if rwSpan, ok := fsc.CtxRootSpan(ctx).(otelsdktrace.ReadWriteSpan); ok {
			rwSpan.SetName(fmt.Sprintf("%s %s", rwSpan.Name(), spanNameSuffix))
			rwSpan.SetStatus(codes.Ok, "MCP request success")
			rwSpan.SetAttributes(spanAttributes...)
			rwSpan.End()
		}

		log.CtxLogger(ctx).Info().Fields(logFields).Msg("MCP request success")

		if metricsEnabled {
			p.requestsCounter.WithLabelValues(mcpMethod, metricTarget, "success").Inc()
			p.requestsDuration.WithLabelValues(mcpMethod, metricTarget).Observe(latency.Seconds())
		}
	})

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		latency := time.Since(fsc.CtxStartTime(ctx))

		mcpMethod := string(method)

		errMessage := fmt.Sprintf("%v", err)

		spanNameSuffix := mcpMethod

		spanAttributes := []attribute.KeyValue{
			attribute.String("mcp.latency", latency.String()),
			attribute.String("mcp.method", mcpMethod),
			attribute.String("mcp.error", errMessage),
		}

		logFields := map[string]any{
			"mcpLatency": latency.String(),
			"mcpMethod":  mcpMethod,
			"mcpError":   errMessage,
		}

		metricTarget := ""

		jsonMessage, err := json.Marshal(message)
		if err == nil {
			if traceRequest {
				spanAttributes = append(spanAttributes, attribute.String("mcp.request", string(jsonMessage)))
			}

			if logRequest {
				logFields["mcpRequest"] = string(jsonMessage)
			}
		}

		//nolint:exhaustive
		switch method {
		case mcp.MethodResourcesRead:
			if req, ok := message.(*mcp.ReadResourceRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.URI)
				spanAttributes = append(spanAttributes, attribute.String("mcp.resourceURI", req.Params.URI))
				logFields["mcpResourceURI"] = req.Params.URI
				metricTarget = req.Params.URI
			}
		case mcp.MethodPromptsGet:
			if req, ok := message.(*mcp.GetPromptRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.Name)
				spanAttributes = append(spanAttributes, attribute.String("mcp.prompt", req.Params.Name))
				logFields["mcpPrompt"] = req.Params.Name
				metricTarget = req.Params.Name
			}
		case mcp.MethodToolsCall:
			if req, ok := message.(*mcp.CallToolRequest); ok {
				spanNameSuffix = fmt.Sprintf("%s %s", spanNameSuffix, req.Params.Name)
				spanAttributes = append(spanAttributes, attribute.String("mcp.tool", req.Params.Name))
				logFields["mcpTool"] = req.Params.Name
				metricTarget = req.Params.Name
			}
		}

		if rwSpan, ok := fsc.CtxRootSpan(ctx).(otelsdktrace.ReadWriteSpan); ok {
			rwSpan.SetName(fmt.Sprintf("%s %s", rwSpan.Name(), spanNameSuffix))
			rwSpan.RecordError(err)
			rwSpan.SetStatus(codes.Error, errMessage)
			rwSpan.SetAttributes(spanAttributes...)
			rwSpan.End()
		}

		log.CtxLogger(ctx).Error().Fields(logFields).Msg("MCP request error")

		if metricsEnabled {
			p.requestsCounter.WithLabelValues(mcpMethod, metricTarget, "error").Inc()
			p.requestsDuration.WithLabelValues(mcpMethod, metricTarget).Observe(latency.Seconds())
		}
	})

	return hooks
}

// Reset resets the MCP requests metrics.
func (p *DefaultMCPServerHooksProvider) Reset() {
	p.requestsCounter.Reset()
	p.requestsDuration.Reset()
}
