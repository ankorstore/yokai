package fxmcpserver_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/prompt"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resource"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/tool"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
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
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var mcpServer *server.MCPServer
	var handler sse.MCPSSEServerContextHandler
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry
	var info *fxmcpserver.MCPServerModuleInfo

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxgenerate.FxGenerateModule,
		fxmetrics.FxMetricsModule,
		fxmcpserver.FxMCPServerModule,
		fx.Options(
			fxmcpserver.AsMCPServerTools(tool.NewTestTool),
			fxmcpserver.AsMCPServerPrompts(prompt.NewTestPrompt),
			fxmcpserver.AsMCPServerResources(resource.NewTestResource),
			fxmcpserver.AsMCPServerResourceTemplates(resource.NewTestResourceTemplate),
		),
		fx.Populate(&mcpServer, &handler, &logBuffer, &traceExporter, &metricsRegistry, &info),
	).RequireStart().RequireStop()

	// test server
	testServer := server.NewTestServer(mcpServer, server.WithSSEContextFunc(handler.Handle()))
	defer testServer.Close()

	// test client
	testClient, err := client.NewSSEMCPClient(testServer.URL + "/sse")
	assert.NoError(t, err)

	// start the client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = testClient.Start(ctx)
	assert.NoError(t, err)
	defer testClient.Close()

	logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
		"level":   "info",
		"message": "MCP session registered",
	})

	// send initialize request
	expectedRequest := `{"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}`
	expectedResponse := `{"protocolVersion":"2024-11-05","capabilities":{"logging":{},"prompts":{},"resources":{},"tools":{}},"serverInfo":{"name":"test-server","version":"1.0.0"}}`

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}

	initResult, err := testClient.Initialize(ctx, initRequest)
	assert.NoError(t, err)
	assert.Equal(t, "test-server", initResult.ServerInfo.Name)
	assert.Equal(t, "1.0.0", initResult.ServerInfo.Version)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
		"level":       "info",
		"mcpMethod":   "initialize",
		"mcpRequest":  expectedRequest,
		"mcpResponse": expectedResponse,
		"message":     "MCP request success",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"MCP initialize",
		attribute.String("mcp.method", "initialize"),
		attribute.String("mcp.request", expectedRequest),
		attribute.String("mcp.response", expectedResponse),
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

	// send tool/list request
	expectedRequest = `{"method":"tools/list","params":{}}`
	expectedResponse = `{"tools":[{"annotations":{"destructiveHint":true,"openWorldHint":true},"description":"Test tool.","inputSchema":{"properties":{},"type":"object"},"name":"test-tool"}]}`

	toolsListRequest := mcp.ListToolsRequest{}

	toolsListResult, err := testClient.ListTools(ctx, toolsListRequest)
	assert.NoError(t, err)
	assert.Len(t, toolsListResult.Tools, 1)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
		"level":       "info",
		"mcpMethod":   "tools/list",
		"mcpRequest":  expectedRequest,
		"mcpResponse": expectedResponse,
		"message":     "MCP request success",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"MCP tools/list",
		attribute.String("mcp.method", "tools/list"),
		attribute.String("mcp.request", expectedRequest),
		attribute.String("mcp.response", expectedResponse),
	)

	expectedMetric = `
		# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
		# TYPE foo_bar_mcp_server_requests_total counter
        foo_bar_mcp_server_requests_total{method="initialize",status="success",target=""} 1
        foo_bar_mcp_server_requests_total{method="tools/list",status="success",target=""} 1
	`
	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"foo_bar_mcp_server_requests_total",
	)
	assert.NoError(t, err)
}
