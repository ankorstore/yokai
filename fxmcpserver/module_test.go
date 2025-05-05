package fxmcpserver_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	fs "github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/tool"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestMCPServerModule(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var mcpServer *server.MCPServer
	var handler sse.MCPSSEServerContextHandler
	var checker *healthcheck.Checker
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxgenerate.FxGenerateModule,
		fxmetrics.FxMetricsModule,
		fxhealthcheck.FxHealthcheckModule,
		fxmcpserver.FxMCPServerModule,
		fx.Options(
			fxmcpserver.AsMCPServerTools(tool.NewAdvancedTestTool),
			fxhealthcheck.AsCheckerProbe(fs.NewMCPServerProbe),
		),
		fx.Supply(fx.Annotate(context.Background(), fx.As(new(context.Context)))),
		fx.Populate(&mcpServer, &handler, &checker, &logBuffer, &traceExporter, &metricsRegistry),
	).RequireStart().RequireStop()

	// create test server
	testServer := server.NewTestServer(mcpServer, server.WithSSEContextFunc(handler.Handle()))
	defer testServer.Close()

	// health check
	checkResult := checker.Check(context.Background(), healthcheck.Readiness)
	assert.False(t, checkResult.Success)
	assert.Equal(t, "MCP SSE server is not running", checkResult.ProbesResults["mcpserver"].Message)

	// create test client
	testClient, err := client.NewSSEMCPClient(testServer.URL + "/sse")
	assert.NoError(t, err)

	// start the client
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = testClient.Start(ctx)
	assert.NoError(t, err)
	defer testClient.Close()

	logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
		"level":   "info",
		"message": "MCP session registered",
	})

	// send initialize request
	initializeRequest := mcp.InitializeRequest{}
	initializeRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initializeRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}

	initializeResult, err := testClient.Initialize(ctx, mcp.InitializeRequest{})
	assert.NoError(t, err)

	assert.Equal(t, "test-server", initializeResult.ServerInfo.Name)
	assert.Equal(t, "1.0.0", initializeResult.ServerInfo.Version)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
		"level":        "info",
		"mcpMethod":    "initialize",
		"mcpTransport": "sse",
		"message":      "MCP request success",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"MCP initialize",
		attribute.String("mcp.method", "initialize"),
		attribute.String("mcp.transport", "sse"),
	)

	expectedMetric := `
		# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
		# TYPE foo_bar_mcp_server_requests_total counter
		foo_bar_mcp_server_requests_total{method="initialize",status="success",target=""} 1
	`
	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"foo_bar_mcp_server_requests_total",
	)
	assert.NoError(t, err)

	// send success tools/call request
	expectedRequest := `{"method":"tools/call","params":{"name":"advanced-test-tool","arguments":{"shouldFail":"false"}}}`
	expectedResponse := `{"content":[{"type":"text","text":"test"}]}`

	callToolRequest := mcp.CallToolRequest{}
	callToolRequest.Params.Name = "advanced-test-tool"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"shouldFail": "false",
	}

	_, err = testClient.CallTool(ctx, callToolRequest)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
		"level":        "info",
		"mcpMethod":    "tools/call",
		"mcpTool":      "advanced-test-tool",
		"mcpRequest":   expectedRequest,
		"mcpResponse":  expectedResponse,
		"mcpTransport": "sse",
		"message":      "MCP request success",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"MCP tools/call advanced-test-tool",
		attribute.String("mcp.method", "tools/call"),
		attribute.String("mcp.tool", "advanced-test-tool"),
		attribute.String("mcp.request", expectedRequest),
		attribute.String("mcp.response", expectedResponse),
		attribute.String("mcp.transport", "sse"),
	)

	expectedMetric = `
		# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
		# TYPE foo_bar_mcp_server_requests_total counter
		foo_bar_mcp_server_requests_total{method="initialize",status="success",target=""} 1
		foo_bar_mcp_server_requests_total{method="tools/call",status="success",target="advanced-test-tool"} 1
	`
	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"foo_bar_mcp_server_requests_total",
	)
	assert.NoError(t, err)

	// send failing tools/call request
	expectedRequest = `{"method":"tools/call","params":{"name":"advanced-test-tool","arguments":{"shouldFail":"true"}}}`

	callToolRequest = mcp.CallToolRequest{}
	callToolRequest.Params.Name = "advanced-test-tool"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"shouldFail": "true",
	}

	_, err = testClient.CallTool(ctx, callToolRequest)
	assert.Error(t, err)
	assert.Equal(t, "advanced tool test failure", err.Error())

	logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
		"level":        "error",
		"mcpError":     "request error: advanced tool test failure",
		"mcpMethod":    "tools/call",
		"mcpTool":      "advanced-test-tool",
		"mcpRequest":   expectedRequest,
		"mcpTransport": "sse",
		"message":      "MCP request error",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"MCP tools/call advanced-test-tool",
		attribute.String("mcp.method", "tools/call"),
		attribute.String("mcp.tool", "advanced-test-tool"),
		attribute.String("mcp.request", expectedRequest),
		attribute.String("mcp.transport", "sse"),
	)

	expectedMetric = `
		# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
		# TYPE foo_bar_mcp_server_requests_total counter
		foo_bar_mcp_server_requests_total{method="initialize",status="success",target=""} 1
		foo_bar_mcp_server_requests_total{method="tools/call",status="success",target="advanced-test-tool"} 1
		foo_bar_mcp_server_requests_total{method="tools/call",status="error",target="advanced-test-tool"} 1
	`
	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"foo_bar_mcp_server_requests_total",
	)
	assert.NoError(t, err)
}
