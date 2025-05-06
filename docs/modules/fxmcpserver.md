---
title: Modules - MCP Server
icon: material/cube-outline
---

# :material-cube-outline: MCP Server Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxmcpserver-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxmcpserver-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxmcpserver)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxmcpserver)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxmcpserver)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxmcpserver)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxmcpserver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxmcpserver)](https://pkg.go.dev/github.com/ankorstore/yokai/fxmcpserver)

## Overview

Yokai provides a [fxmcpserver](https://github.com/ankorstore/yokai/tree/main/fxmcpserver) module, offering an [MCP server](https://modelcontextprotocol.io/introduction) to your application.

It wraps the [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) module.

It comes with:

- automatic panic recovery
- automatic requests logging and tracing (method, target, duration, ...)
- automatic requests metrics (count and duration)
- possibility to register MCP resources, resource templates, prompts and tools
- possibility to expose the MCP server via Stdio (local) and/or HTTP SSE (remote)

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxmcpserver
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxmcpserver"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxmcpserver.FxMCPServerModule,
	// ...
)
```

## Configuration

```yaml title="configs/config.yaml"
modules:
  mcp:
    server:
      name: "MCP Server"                  # server name ("MCP server" by default)
      version: 1.0.0                      # server version (1.0.0 by default)
      capabilities:
        resources: true                   # to expose MCP resources & resource templates (disabled by default)
        prompts: true                     # to expose MCP prompts (disabled by default)
        tools: true                       # to expose MCP tools (disabled by default)
      transport:
        sse:
          expose: true                    # to remotely expose the MCP server via SSE (disabled by default)
          address: ":8082"                # exposition address (":8082" by default)
          base_url: ""                    # base url ("" by default)
          base_path: ""                   # base path ("" by default)
          sse_endpoint: "/sse"            # SSE endpoint ("/sse" by default)
          message_endpoint: "/message"    # message endpoint ("/message" by default)
          keep_alive: true                # to keep connection alive
          keep_alive_interval: 10         # keep alive interval in seconds (10 by default)
        stdio:
          expose: false                   # to locally expose the MCP server via Stdio (disabled by default)
      log:
        request: true                     # to log MCP requests contents (disabled by default)
        response: true                    # to log MCP responses contents (disabled by default)
      trace:
        request: true                     # to trace MCP requests contents (disabled by default)
        response: true                    # to trace MCP responses contents (disabled by default)
      metrics:
        collect:
          enabled: true                   # to collect MCP server metrics (disabled by default)
          namespace: foo                  # MCP server metrics namespace ("" by default)
          subsystem: bar                  # MCP server metrics subsystem ("" by default)
        buckets: 0.1, 1, 10               # to override default request duration buckets
```

## Usage

This module offers the possibility to easily register MCP resources, resource templates, prompts and tools.

### Resources registration

This module offers an [MCPServerResource](https://github.com/ankorstore/yokai/blob/main/fxmcpserver/server/registry.go) interface to implement to provide an [MCP resource](https://modelcontextprotocol.io/docs/concepts/resources).

For example, an MCP resource that reads a file path coming from the configuration:

```go title="internal/mcp/resource/readme.go"
package resource

import (
	"context"
	"os"
	
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ReadmeResource struct {
	config *config.Config
}

func NewReadmeResource(config *config.Config) *ReadmeResource {
	return &ReadmeResource{
		config: config,
	}
}

func (r *ReadmeResource) Name() string {
	return "readme"
}

func (r *ReadmeResource) URI() string {
	return "docs://readme"
}

func (r *ReadmeResource) Options() []mcp.ResourceOption {
	return []mcp.ResourceOption{
		mcp.WithResourceDescription("Project README"),
	}
}

func (r *ReadmeResource) Handle() server.ResourceHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content, err := os.ReadFile(r.config.GetString("config.readme.path"))
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		}, nil
	}
}
```

You can register your MCP resource:

- with `AsMCPServerResource()` to register a single MCP resource
- with `AsMCPServerResources()` to register several MCP resources at once

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/foo/bar/internal/mcp/resource"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// registers ReadmeResource as MCP resource
		fxmcpserver.AsMCPServerResource(resource.NewReadmeResource),
		// ...
	)
}
```

The dependencies of your MCP resources will be autowired.

To expose it, you need to ensure that the MCP server has the `resources` capability enabled:

```yaml title="configs/config.yaml"
modules:
  mcp:
    server:
      capabilities:
        resources: true # to expose MCP resources & resource templates (disabled by default)
```

### Resource templates registration

This module offers an [MCPServerResourceTemplate](https://github.com/ankorstore/yokai/blob/main/fxmcpserver/server/registry.go) interface to implement to provide an [MCP resource template](https://modelcontextprotocol.io/docs/concepts/resources).

For example, an MCP resource template that retrieves a user profile for a given id:

```go title="internal/mcp/resource/readme.go"
package resource

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/foo/bar/internal/user"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type UserProfileResource struct {
	repository *user.Respository
}

func NewUserProfileResource(repository *user.Respository) *UserProfileResource {
	return &UserProfileResource{
		repository: repository,
	}
}

func (r *UserProfileResource) Name() string {
	return "user-profile"
}

func (r *UserProfileResource) URI() string {
	return "users://{id}/profile"
}

func (r *UserProfileResource) Options() []mcp.ResourceTemplateOption {
	return []mcp.ResourceTemplateOption{
		mcp.WithTemplateDescription("User profile"),
	}
}

func (r *UserProfileResource) Handle() server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// some user id extraction logic
		userID := extractUserIDFromURI(request.Params.URI)

		// find user profile by user id
		user, err := r.repository.Find(userID)
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     user,
			},
		}, nil
	}
}
```

You can register your MCP resource template:

- with `AsMCPServerResourceTemplate()` to register a single MCP resource template
- with `AsMCPServerResourceTemplates()` to register several MCP resource templates at once

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/foo/bar/internal/mcp/resource"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// registers UserProfileResource as MCP resource
		fxmcpserver.AsMCPServerResourceTemplate(resource.NewUserProfileResource),
		// ...
	)
}
```

The dependencies of your MCP resource templates will be autowired.

To expose it, you need to ensure that the MCP server has the `resources` capability enabled:

```yaml title="configs/config.yaml"
modules:
  mcp:
    server:
      capabilities:
        resources: true # to expose MCP resources & resource templates (disabled by default)
```

### Prompts registration

This module offers an [MCPServerPrompt](https://github.com/ankorstore/yokai/blob/main/fxmcpserver/server/registry.go) interface to implement to provide an [MCP prompt](https://modelcontextprotocol.io/docs/concepts/prompts).

For example, an MCP prompt that greets a provided user name:

```go title="internal/mcp/prompt/greet.go"
package prompt

import (
	"context"
	"fmt"
	
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GreetingPrompt struct {
	config *config.Config
}

func NewGreetingPrompt(config *config.Config) *GreetingPrompt {
	return &GreetingPrompt{
		config: config,
	}
}

func (p *GreetingPrompt) Name() string {
	return "greeting"
}

func (p *GreetingPrompt) Options() []mcp.PromptOption {
	return []mcp.PromptOption{
		mcp.WithPromptDescription("A friendly greeting prompt"),
		mcp.WithArgument(
			"name",
			mcp.ArgumentDescription("Name of the person to greet"),
		),
	}
}

func (p *GreetingPrompt) Handle() server.PromptHandlerFunc {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		name := request.Params.Arguments["name"]
		if name == "" {
			name = "friend"
		}

		return mcp.NewGetPromptResult(
			"A friendly greeting",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent(fmt.Sprintf("Hello, %s! I am %s. How can I help you today?", name, p.config.GetString("config.assistant.name"))),
				),
			},
		), nil
	}
}
```

You can register your MCP prompt:

- with `AsMCPServerPrompt()` to register a single MCP prompt
- with `AsMCPServerPrompts()` to register several MCP prompts at once

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/foo/bar/internal/mcp/prompt"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// registers GreetingPrompt as MCP prompt
		fxmcpserver.AsMCPServerPrompt(prompt.NewGreetingPrompt),
		// ...
	)
}
```

The dependencies of your MCP prompts will be autowired.

To expose it, you need to ensure that the MCP server has the `prompts` capability enabled:

```yaml title="configs/config.yaml"
modules:
  mcp:
    server:
      capabilities:
        prompts: true # to expose MCP prompts (disabled by default)
```

### Tools registration

This module offers an [MCPServerTool](https://github.com/ankorstore/yokai/blob/main/fxmcpserver/server/registry.go) interface to implement to provide an [MCP tool](https://modelcontextprotocol.io/docs/concepts/tools).

For example, an MCP tool that performs basic arithmetic calculations:

```go title="internal/mcp/tool/calculator.go"
package tool

import (
	"context"
	"fmt"
	
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type CalculatorTool struct {
	config *config.Config
}

func NewCalculatorTool(config *config.Config) *CalculatorTool {
	return &CalculatorTool{
		config: config,
	}
}

func (t *CalculatorTool) Name() string {
	return "calculator"
}

func (t *CalculatorTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Perform basic arithmetic calculations"),
		mcp.WithString(
			"operation",
			mcp.Required(),
			mcp.Description("The arithmetic operation to perform"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber(
			"x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber(
			"y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	}
}

func (t *CalculatorTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		log.CtxLogger(c.Request().Context()).Info().Msg("in calculator tool")
		
		if !t.config.GetBool("config.calculator.enabled") {
			return nil, fmt.Errorf("calculator is not enabled")
		}

		op := request.Params.Arguments["operation"].(string)
		x := request.Params.Arguments["x"].(float64)
		y := request.Params.Arguments["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultError("cannot divide by zero"), nil
			}

			result = x / y
		}

		return mcp.FormatNumberResult(result), nil
	}
}
```

You can register your MCP tool:

- with `AsMCPServerTool()` to register a single MCP tool
- with `AsMCPServerTools()` to register several MCP tools at once

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/foo/bar/internal/mcp/tool"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// registers CalculatorTool as MCP tool
		fxmcpserver.AsMCPServerTool(tool.NewCalculatorTool),
		// ...
	)
}
```

The dependencies of your MCP tools will be autowired.

To expose it, you need to ensure that the MCP server has the `tools` capability enabled:

```yaml title="configs/config.yaml"
modules:
  mcp:
    server:
      capabilities:
        tools: true # to expose MCP tools (disabled by default)
```

## Logging

You can configure the MCP server requests and responses automatic logging:

```yaml title="configs/config.yaml"
modules:
  mcp:
    server:
      log:
        request: true   # to log MCP requests contents (disabled by default)
        response: true  # to log MCP responses contents (disabled by default)
```

As a result, in your application logs:

```
INF in calculator tool mcpRequestID=460aab37-e16e-4464-9956-54fce47746e7 mcpSessionID=8f617d54-e4c9-4459-bb26-76b4d96e2b72 mcpTransport=sse service=yokai-mcp spanID=0f536ffa84fb8800 system=mcpserver traceID=594a9585cbfd5362c03968cd6d7d786c
INF MCP request success mcpLatency=4.869308ms mcpMethod=tools/call mcpRequest="..." mcpResponse="..." mcpRequestID=460aab37-e16e-4464-9956-54fce47746e7 mcpSessionID=8f617d54-e4c9-4459-bb26-76b4d96e2b72 mcpTool=calculator mcpTransport=sse service=yokai-mcp spanID=0f536ffa84fb8800 system=mcpserver traceID=594a9585cbfd5362c03968cd6d7d786c
```

If both HTTP server logging and tracing are enabled, log records will automatically have the current `traceID` and `spanID` to be able to correlate logs and trace spans.

To get logs correlation in your MCP registrations, you need to retrieve the logger from the context with `log.CtxLogger()`:

```go
log.CtxLogger(c.Request().Context()).Info().Msg("in calculator tool")
```

The MCP server logging will be based on the [log](fxlog.md) module configuration.

## Tracing

You can configure the MCP server requests and responses automatic tracing:

```yaml title="configs/config.yaml"
modules:
  mcp:
    server:
      trace:
        request: true   # to trace MCP requests contents (disabled by default)
        response: true  # to trace MCP responses contents (disabled by default)
```

As a result, in your application trace spans attributes:

```
service.name: yokai-mcp
mcp.method: tools/call
mcp.tool: calculator
mcp.transport: sse
mcp.request: ...
mcp.response: ...
...
```

To get traces correlation in your MCP registrations, you need to retrieve the tracer from the context with `trace.CtxTracer()`:

```go
ctx, span := trace.CtxTracer(ctx).Start(ctx, "in calculator tool")
defer span.End()
```

The MCP server tracing will be based on the [fxtrace](trace.md) module configuration.

## Metrics

You can enable MCP requests automatic metrics with `modules.mcp.server.metrics.collect.enable=true`:

```yaml title="configs/config.yaml"
modules:
  mcp:
    server:
      metrics:
        collect:
          enabled: true      # to collect MCP server metrics (disabled by default)
          namespace: foo     # MCP server metrics namespace ("" by default)
          subsystem: bar     # MCP server metrics subsystem ("" by default)
        buckets: 0.1, 1, 10  # to override default request duration buckets
```

For example, after calling the `calculator` MCP tool, the [core](fxcore.md) HTTP server will expose in the configured metrics endpoint:

```makefile title="[GET] /metrics"
# ...
# HELP mcp_server_requests_duration_seconds Time spent processing MCP requests
# TYPE mcp_server_requests_duration_seconds histogram
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="0.005"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="0.01"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="0.025"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="0.05"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="0.1"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="0.25"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="0.5"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="1"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="2.5"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="5"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="10"} 1
mcp_server_requests_duration_seconds_bucket{method="tools/call",target="calculator",le="+Inf"} 1
mcp_server_requests_duration_seconds_sum{method="tools/call",target="calculator"} 0.004869308
mcp_server_requests_duration_seconds_count{method="tools/call",target="calculator"} 1
# HELP mcp_server_requests_total Number of processed MCP requests
# TYPE mcp_server_requests_total counter
mcp_server_requests_total{method="tools/call",status="success",target="calculator"} 1
```

## Testing

This module provides a [MCPSSETestServer](https://github.com/ankorstore/yokai/blob/main/fxmcpserver/fxmcpservertest/server.go) to enable you to easily test your exposed MCP registrations.

From this server, you can create a ready to use client via `StartClient()` to perform MCP requests, to functionally test your MCP server.

You can easily assert on:

- MCP responses
- logs
- traces
- metrics

For example, a test an `MCP ping`:

```go title="internal/mcp/ping_test.go"
package handler_test

import (
	"testing"

	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/foo/bar/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func TestMCPPing(t *testing.T) {
	var testServer *fxmcpservertest.MCPSSETestServer
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	internal.RunTest(t, fx.Populate(&testServer, &logBuffer, &traceExporter, &metricsRegistry))

	// close the test server once done
	defer testServer.Close()

	// start test client
	testClient, err := testServer.StartClient(context.Background())
	assert.NoError(t, err)

	// close the test client once done
	defer testClient.Close()

	// send MCP ping request
	err = testClient.Ping(context.Background())
	assert.NoError(t, err)

	// assertion on the logs buffer
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "info",
		"mcpMethod":    "ping",
		"mcpTransport": "sse",
		"message":      "MCP request success",
	})

	// assertion on the traces exporter
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"MCP ping",
		attribute.String("mcp.method", "ping"),
		attribute.String("mcp.transport", "sse"),
	)

	// assertion on the metrics registry
	expectedMetric := `
		# HELP mcp_server_requests_total Number of processed HTTP requests
		# TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="ping",status="success",target=""} 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"mcp_server_requests_total",
	)
	assert.NoError(t, err)
}
```
