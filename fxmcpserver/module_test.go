package fxmcpserver_test

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmcpserver/fxmcpservertest"
	fs "github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/hook"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/prompt"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resource"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resourcetemplate"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/tool"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

//nolint:maintidx,forcetypeassert
func TestMCPServerModule(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var testMCPStreamableHTTPServer *fxmcpservertest.MCPStreamableHTTPTestServer
	var testMCPSSEServer *fxmcpservertest.MCPSSETestServer
	var provider fs.MCPServerHooksProvider
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
			fxmcpserver.AsMCPServerTools(tool.NewSimpleTestTool, tool.NewAdvancedTestTool),
			fxmcpserver.AsMCPServerPrompts(prompt.NewSimpleTestPrompt),
			fxmcpserver.AsMCPServerResources(resource.NewSimpleTestResource),
			fxmcpserver.AsMCPServerResourceTemplates(resourcetemplate.NewSimpleTestResourceTemplate),
			fxmcpserver.AsMCPSSEServerContextHooks(hook.NewSimpleMCPSSEServerContextHook),
			fxmcpserver.AsMCPStreamableHTTPServerContextHooks(hook.NewSimpleMCPStreamableHTTPServerContextHook),
			fxhealthcheck.AsCheckerProbe(fs.NewMCPServerProbe),
		),
		fx.Supply(fx.Annotate(context.Background(), fx.As(new(context.Context)))),
		fx.Populate(
			&testMCPStreamableHTTPServer,
			&testMCPSSEServer,
			&provider,
			&checker,
			&logBuffer,
			&traceExporter,
			&metricsRegistry,
		),
	).RequireStart().RequireStop()

	// ensure test servers closure
	defer func() {
		testMCPStreamableHTTPServer.Close()
		testMCPSSEServer.Close()
	}()

	ctx := context.Background()

	// health check
	checkResult := checker.Check(context.Background(), healthcheck.Readiness)
	assert.False(t, checkResult.Success)
	assert.Equal(
		t,
		"MCP StreamableHTTP server is not running, MCP SSE server is not running",
		checkResult.ProbesResults["mcpserver"].Message,
	)

	// start test clients
	testMCPStreamableHTTPClient, err := testMCPStreamableHTTPServer.StartClient(ctx)
	assert.NoError(t, err)

	testMCPSSEClient, err := testMCPSSEServer.StartClient(ctx)
	assert.NoError(t, err)
	defer testMCPSSEClient.Close()

	// hooks provider
	defaultProvider, ok := provider.(*fs.DefaultMCPServerHooksProvider)
	assert.True(t, ok)

	tests := []struct {
		name      string
		client    *client.Client
		transport string
	}{
		{
			"with StreamableHTTP transport",
			testMCPStreamableHTTPClient,
			"streamable-http",
		},
		{
			"with SSE transport",
			testMCPSSEClient,
			"sse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset o11y
			logBuffer.Reset()
			traceExporter.Reset()
			defaultProvider.Reset()

			// send success tools/call request
			expectedRequest := `{"method":"tools/call","params":{"name":"advanced-test-tool","arguments":{"shouldFail":"false"}}}`
			expectedResponse := `{"content":[{"type":"text","text":"test"}]}`

			callToolRequest := mcp.CallToolRequest{}
			callToolRequest.Params.Name = "advanced-test-tool"
			callToolRequest.Params.Arguments = map[string]interface{}{
				"shouldFail": "false",
			}

			callToolResult, err := tt.client.CallTool(ctx, callToolRequest)
			assert.NoError(t, err)
			assert.False(t, callToolResult.IsError)

			logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
				"level":        "info",
				"mcpMethod":    "tools/call",
				"mcpTool":      "advanced-test-tool",
				"mcpRequest":   expectedRequest,
				"mcpResponse":  expectedResponse,
				"mcpTransport": tt.transport,
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
				attribute.String("mcp.transport", tt.transport),
			)

			expectedMetric := `
				# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
				# TYPE foo_bar_mcp_server_requests_total counter
				foo_bar_mcp_server_requests_total{method="tools/call",status="success",target="advanced-test-tool"} 1
			`
			err = testutil.GatherAndCompare(
				metricsRegistry,
				strings.NewReader(expectedMetric),
				"foo_bar_mcp_server_requests_total",
			)
			assert.NoError(t, err)

			// send error tools/call request
			expectedRequest = `{"method":"tools/call","params":{"name":"advanced-test-tool","arguments":{"shouldFail":"true"}}}`

			callToolRequest = mcp.CallToolRequest{}
			callToolRequest.Params.Name = "advanced-test-tool"
			callToolRequest.Params.Arguments = map[string]interface{}{
				"shouldFail": "true",
			}

			_, err = tt.client.CallTool(ctx, callToolRequest)
			assert.Error(t, err)
			assert.Equal(t, "advanced tool test failure", err.Error())

			logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
				"level":        "error",
				"mcpError":     "request error: advanced tool test failure",
				"mcpMethod":    "tools/call",
				"mcpTool":      "advanced-test-tool",
				"mcpRequest":   expectedRequest,
				"mcpTransport": tt.transport,
				"message":      "MCP request error",
			})

			tracetest.AssertHasTraceSpan(
				t,
				traceExporter,
				"MCP tools/call advanced-test-tool",
				attribute.String("mcp.method", "tools/call"),
				attribute.String("mcp.tool", "advanced-test-tool"),
				attribute.String("mcp.request", expectedRequest),
				attribute.String("mcp.transport", tt.transport),
			)

			expectedMetric = `
				# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
				# TYPE foo_bar_mcp_server_requests_total counter
				foo_bar_mcp_server_requests_total{method="tools/call",status="success",target="advanced-test-tool"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="error",target="advanced-test-tool"} 1
			`
			err = testutil.GatherAndCompare(
				metricsRegistry,
				strings.NewReader(expectedMetric),
				"foo_bar_mcp_server_requests_total",
			)
			assert.NoError(t, err)

			// send success prompts/get request
			expectedRequest = `{"method":"prompts/get","params":{"name":"simple-test-prompt"}}`
			expectedResponse = `{"description":"ok","messages":[{"role":"assistant","content":{"type":"text","text":"context hook value: bar"}}]}`

			getPromptRequest := mcp.GetPromptRequest{}
			getPromptRequest.Params.Name = "simple-test-prompt"

			getPromptResult, err := tt.client.GetPrompt(ctx, getPromptRequest)
			assert.NoError(t, err)
			assert.Equal(t, mcp.RoleAssistant, getPromptResult.Messages[0].Role)
			assert.Equal(t, "context hook value: bar", getPromptResult.Messages[0].Content.(mcp.TextContent).Text)

			logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
				"level":        "info",
				"mcpMethod":    "prompts/get",
				"mcpPrompt":    "simple-test-prompt",
				"mcpRequest":   expectedRequest,
				"mcpResponse":  expectedResponse,
				"mcpTransport": tt.transport,
				"message":      "MCP request success",
			})

			tracetest.AssertHasTraceSpan(
				t,
				traceExporter,
				"MCP prompts/get simple-test-prompt",
				attribute.String("mcp.method", "prompts/get"),
				attribute.String("mcp.prompt", "simple-test-prompt"),
				attribute.String("mcp.request", expectedRequest),
				attribute.String("mcp.response", expectedResponse),
				attribute.String("mcp.transport", tt.transport),
			)

			expectedMetric = `
				# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
				# TYPE foo_bar_mcp_server_requests_total counter
				foo_bar_mcp_server_requests_total{method="prompts/get",status="success",target="simple-test-prompt"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="success",target="advanced-test-tool"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="error",target="advanced-test-tool"} 1
			`
			err = testutil.GatherAndCompare(
				metricsRegistry,
				strings.NewReader(expectedMetric),
				"foo_bar_mcp_server_requests_total",
			)
			assert.NoError(t, err)

			// send error prompts/get request
			expectedRequest = `{"method":"prompts/get","params":{"name":"invalid-test-prompt"}}`

			getPromptRequest = mcp.GetPromptRequest{}
			getPromptRequest.Params.Name = "invalid-test-prompt"

			_, err = tt.client.GetPrompt(ctx, getPromptRequest)
			assert.Error(t, err)
			assert.Equal(t, "prompt 'invalid-test-prompt' not found: prompt not found", err.Error())

			logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
				"level":        "error",
				"mcpError":     "request error: prompt 'invalid-test-prompt' not found: prompt not found",
				"mcpMethod":    "prompts/get",
				"mcpPrompt":    "invalid-test-prompt",
				"mcpRequest":   expectedRequest,
				"mcpTransport": tt.transport,
				"message":      "MCP request error",
			})

			tracetest.AssertHasTraceSpan(
				t,
				traceExporter,
				"MCP prompts/get invalid-test-prompt",
				attribute.String("mcp.method", "prompts/get"),
				attribute.String("mcp.prompt", "invalid-test-prompt"),
				attribute.String("mcp.request", expectedRequest),
				attribute.String("mcp.transport", tt.transport),
			)

			expectedMetric = `
				# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
				# TYPE foo_bar_mcp_server_requests_total counter
				foo_bar_mcp_server_requests_total{method="prompts/get",status="error",target="invalid-test-prompt"} 1
				foo_bar_mcp_server_requests_total{method="prompts/get",status="success",target="simple-test-prompt"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="success",target="advanced-test-tool"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="error",target="advanced-test-tool"} 1
			`
			err = testutil.GatherAndCompare(
				metricsRegistry,
				strings.NewReader(expectedMetric),
				"foo_bar_mcp_server_requests_total",
			)
			assert.NoError(t, err)

			// send success resources/get request
			expectedRequest = `{"method":"resources/read","params":{"uri":"simple-test://resources"}}`
			expectedResponse = `{"contents":[{"uri":"simple-test://resources","mimeType":"text/plain","text":"simple test resource"}]}`

			readResourceRequest := mcp.ReadResourceRequest{}
			readResourceRequest.Params.URI = "simple-test://resources"

			readResourceResult, err := tt.client.ReadResource(ctx, readResourceRequest)
			assert.NoError(t, err)
			assert.Equal(t, "simple test resource", readResourceResult.Contents[0].(mcp.TextResourceContents).Text)

			logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
				"level":          "info",
				"mcpMethod":      "resources/read",
				"mcpResourceURI": "simple-test://resources",
				"mcpRequest":     expectedRequest,
				"mcpResponse":    expectedResponse,
				"mcpTransport":   tt.transport,
				"message":        "MCP request success",
			})

			tracetest.AssertHasTraceSpan(
				t,
				traceExporter,
				"MCP resources/read simple-test://resources",
				attribute.String("mcp.method", "resources/read"),
				attribute.String("mcp.resourceURI", "simple-test://resources"),
				attribute.String("mcp.request", expectedRequest),
				attribute.String("mcp.response", expectedResponse),
				attribute.String("mcp.transport", tt.transport),
			)

			expectedMetric = `
				# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
				# TYPE foo_bar_mcp_server_requests_total counter
				foo_bar_mcp_server_requests_total{method="prompts/get",status="error",target="invalid-test-prompt"} 1
				foo_bar_mcp_server_requests_total{method="prompts/get",status="success",target="simple-test-prompt"} 1
				foo_bar_mcp_server_requests_total{method="resources/read",status="success",target="simple-test://resources"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="success",target="advanced-test-tool"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="error",target="advanced-test-tool"} 1
			`
			err = testutil.GatherAndCompare(
				metricsRegistry,
				strings.NewReader(expectedMetric),
				"foo_bar_mcp_server_requests_total",
			)
			assert.NoError(t, err)

			// send error resources/get request
			expectedRequest = `{"method":"resources/read","params":{"uri":"simple-test://invalid"}}`

			readResourceRequest = mcp.ReadResourceRequest{}
			readResourceRequest.Params.URI = "simple-test://invalid"

			_, err = tt.client.ReadResource(ctx, readResourceRequest)
			assert.Error(t, err)
			assert.Equal(t, "handler not found for resource URI 'simple-test://invalid': resource not found", err.Error())

			logtest.AssertHasLogRecord(t, logBuffer, map[string]any{
				"level":          "error",
				"mcpError":       "request error: handler not found for resource URI 'simple-test://invalid': resource not found",
				"mcpMethod":      "resources/read",
				"mcpResourceURI": "simple-test://invalid",
				"mcpRequest":     expectedRequest,
				"mcpTransport":   tt.transport,
				"message":        "MCP request error",
			})

			tracetest.AssertHasTraceSpan(
				t,
				traceExporter,
				"MCP resources/read simple-test://invalid",
				attribute.String("mcp.method", "resources/read"),
				attribute.String("mcp.resourceURI", "simple-test://invalid"),
				attribute.String("mcp.request", expectedRequest),
				attribute.String("mcp.transport", tt.transport),
			)

			expectedMetric = `
				# HELP foo_bar_mcp_server_requests_total Number of processed MCP requests
				# TYPE foo_bar_mcp_server_requests_total counter
				foo_bar_mcp_server_requests_total{method="prompts/get",status="error",target="invalid-test-prompt"} 1
				foo_bar_mcp_server_requests_total{method="prompts/get",status="success",target="simple-test-prompt"} 1
				foo_bar_mcp_server_requests_total{method="resources/read",status="error",target="simple-test://invalid"} 1
				foo_bar_mcp_server_requests_total{method="resources/read",status="success",target="simple-test://resources"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="success",target="advanced-test-tool"} 1
				foo_bar_mcp_server_requests_total{method="tools/call",status="error",target="advanced-test-tool"} 1
			`
			err = testutil.GatherAndCompare(
				metricsRegistry,
				strings.NewReader(expectedMetric),
				"foo_bar_mcp_server_requests_total",
			)
			assert.NoError(t, err)
		})
	}
}
