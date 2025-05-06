# Fx MCP Server Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxmcpserver-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxmcpserver-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxmcpserver)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxmcpserver)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxmcpserver)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxmcpserver)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxmcpserver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxmcpserver)](https://pkg.go.dev/github.com/ankorstore/yokai/fxmcpserver)

> [Fx](https://uber-go.github.io/fx/) module for [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go).

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
  * [Dependencies](#dependencies)
  * [Loading](#loading)
  * [Configuration](#configuration)
  * [Registration](#registration)
    * [Resources](#resources)
    * [Resource templates](#resource-templates)
    * [Prompts](#prompts)
    * [Tools](#tools)
  * [Testing](#testing)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxmcpserver
```

## Features

This module provides an [MCP server](https://modelcontextprotocol.io/introduction) to your application with:

- automatic panic recovery
- automatic requests logging and tracing (method, target, duration, ...)
- automatic requests metrics (count and duration)
- possibility to register MCP resources, resource templates, prompts and tools
- possibility to expose the MCP server via Stdio (local) and/or HTTP SSE (remote)

## Documentation

### Dependencies

This module is intended to be used alongside:

- the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module
- the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module
- the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module
- the [fxmetrics](https://github.com/ankorstore/yokai/tree/main/fxmetrics) module
- the [fxgenerate](https://github.com/ankorstore/yokai/tree/main/fxgenerate) module

### Loading

To load the module in your application:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,       // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxmcpserver.FxMCPServerModule, // load the module
	).Run()
}
```

### Configuration

Configuration reference:

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: true
modules:
  log:
    level: info
    output: stdout
  trace:
    processor:
      type: stdout
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

Notes:

- the MCP server logging will be based on the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module configuration
- the MCP server tracing will be based on the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module configuration

### Registration

This module offers the possibility to easily register MCP resources, resource templates, prompts and tools.

#### Resources

This module offers an [MCPServerResource](server/registry.go) interface to implement to provide an [MCP resource](https://modelcontextprotocol.io/docs/concepts/resources).

You can use the `AsMCPServerResource()` function to register an MCP resource, or `AsMCPServerResources()` to register several MCP resources at once.

The dependencies of your MCP resources will be autowired.

```go
package main

import (
	"context"
	"os"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
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

func main() {
	fx.New(
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxmcpserver.FxMCPServerModule,
		fx.Options(
			fxmcpserver.AsMCPServerResource(NewReadmeResource), // registers the ReadmeResource as MCP resource
		),
	).Run()
}
```

To expose it, you need to ensure that the MCP server has the `resources` capability enabled:

```yaml
# ./configs/config.yaml
modules:
  mcp:
    server:
      capabilities:
        resources: true # to expose MCP resources & resource templates (disabled by default)
```

#### Resource templates

This module offers an [MCPServerResourceTemplate](server/registry.go) interface to implement to provide an [MCP resource template](https://modelcontextprotocol.io/docs/concepts/resources).

You can use the `AsMCPServerResourceTemplate()` function to register an MCP resource template, or `AsMCPServerResourceTemplates()` to register several MCP resource templates at once.

The dependencies of your MCP resource templates will be autowired.

```go
package main

import (
	"context"
	
	"github.com/foo/bar/internal/user"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
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

func main() {
	fx.New(
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxmcpserver.FxMCPServerModule,
		fx.Options(
			fxmcpserver.AsMCPServerResourceTemplate(NewUserProfileResource), // registers the UserProfileResource as MCP resource template
		),
	).Run()
}
```

To expose it, you need to ensure that the MCP server has the `resources` capability enabled:

```yaml
# ./configs/config.yaml
modules:
  mcp:
    server:
      capabilities:
        resources: true # to expose MCP resources & resource templates (disabled by default)
```

#### Prompts

This module offers an [MCPServerPrompt](server/registry.go) interface to implement to provide an [MCP prompt](https://modelcontextprotocol.io/docs/concepts/prompts).

You can use the `AsMCPServerPrompt()` function to register an MCP prompt, or `AsMCPServerPrompts()` to register several MCP prompts at once.

The dependencies of your MCP prompts will be autowired.

```go
package main

import (
	"context"
	"os"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
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

func main() {
	fx.New(
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxmcpserver.FxMCPServerModule,
		fx.Options(
			fxmcpserver.AsMCPServerPrompt(NewGreetingPrompt), // registers the GreetingPrompt as MCP prompt
		),
	).Run()
}
```

To expose it, you need to ensure that the MCP server has the `prompts` capability enabled:

```yaml
# ./configs/config.yaml
modules:
  mcp:
    server:
      capabilities:
        prompts: true # to expose MCP prompts (disabled by default)
```

#### Tools

This module offers an [MCPServerTool](server/registry.go) interface to implement to provide an [MCP tool](https://modelcontextprotocol.io/docs/concepts/tools).

You can use the `AsMCPServerTool()` function to register an MCP tool, or `AsMCPServerTools()` to register several MCP tools at once.

The dependencies of your MCP tools will be autowired.

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
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

func main() {
	fx.New(
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxmcpserver.FxMCPServerModule,
		fx.Options(
			fxmcpserver.AsMCPServerTool(NewCalculatorTool), // registers the CalculatorTool as MCP tool
		),
	).Run()
}
```

To expose it, you need to ensure that the MCP server has the `tools` capability enabled:

```yaml
# ./configs/config.yaml
modules:
  mcp:
    server:
      capabilities:
        tools: true # to expose MCP tools (disabled by default)
```

### Testing

This module provides a [MCPSSETestServer](fxmcpservertest/server.go) to enable you to easily test your exposed MCP capabilities.

From this server, you can create a ready to use client via `StartClient()` to perform MCP requests, to functionally test your MCP server.

You can then test it, considering `logs`, `traces` and `metrics` are enabled:

```go
package internal_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmcpserver/fxmcpservertest"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestExample(t *testing.T) {
	var testServer *fxmcpservertest.MCPSSETestServer
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
		fxmcpserver.FxMCPServerModule,
		fx.Populate(&testServer, &logBuffer, &traceExporter, &metricsRegistry),
	).RequireStart().RequireStop()

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

You can find more tests examples in this module own [tests](module_test.go).