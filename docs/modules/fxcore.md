---
title: Modules - Core
icon: material/cube-outline
---

# :material-cube-outline: Core Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxcore-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxcore-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxcore)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxcore)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxcore)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxcore)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxcore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxcore)](https://pkg.go.dev/github.com/ankorstore/yokai/fxcore)

## Overview

Yokai provides a [fxcore](https://github.com/ankorstore/yokai/tree/main/fxcore) module, the `heart of your applications`.

It comes with:

- a [bootstrapper](https://github.com/ankorstore/yokai/blob/main/fxcore/bootstrap.go)
- a [dependency injection](https://en.wikipedia.org/wiki/Dependency_injection) system, based on [Fx](https://github.com/uber-go/fx)
- a dedicated core HTTP server
- pre-enabled [config](fxconfig.md), [health check](fxhealthcheck.md), [log](fxlog.md), [trace](fxtrace.md), [metrics](fxmetrics.md) and [generate](fxgenerate.md) modules
- an extension system for Yokai `built-in`, [contrib](https://github.com/ankorstore/yokai-contrib) or your `own` modules

The `core HTTP server` runs automatically on a dedicated port (default `8081`), to serve:

- the dashboard: UI to get an overview of your application
- the debug endpoints: to expose information about your build, config, loaded modules, etc.
- the health check endpoints: to expose the configured [health check probes](fxhealthcheck.md#probes-registration) of your application
- the metrics endpoint: to expose all [collected metrics](fxmetrics.md#metrics-registration) from your application

Whatever your type of application (HTTP, gRPC, worker, etc.), all `platform concerns` are handled by this
dedicated server:

- to avoid to expose sensitive information (health checks, metrics, debug, etc.) to your users
- and most importantly to enable your application to focus on its logic

## Installation

When you use a Yokai `application template`, you have nothing to install, it's ready to use.

## Configuration

```yaml title="configs/config.yaml"
modules:
  core:
    server:
      expose: true                     # to expose the core http server, disabled by default
      address: ":8081"                 # core http server listener address (default :8081)
      errors:              
        obfuscate: false               # to obfuscate error messages on the core http server responses
        stack: false                   # to add error stack trace to error response of the core http server
      dashboard:
        enabled: true                  # to enable the core dashboard
        overview:      
          app_description: true        # to display the app description on the dashboard overview
          app_env: true                # to display the app env on the dashboard overview
          app_debug: true              # to display the app debug on the dashboard overview
          app_version: true            # to display the app version on the dashboard overview
          log_level: true              # to display the log level on the dashboard overview
          log_output: true             # to display the log output on the dashboard overview
          trace_sampler: true          # to display the trace sampler on the dashboard overview
          trace_processor: true        # to display the trace processor on the dashboard overview
      log:
        headers:                       # to log incoming request headers on the core http server
          x-foo: foo                   # to log for example the header x-foo in the log field foo
          x-bar: bar              
        exclude:                       # to exclude specific routes from logging
          - /healthz
          - /livez
          - /readyz
          - /metrics
        level_from_response: true      # to use response status code for log level (ex: 500=error)
      trace:     
        enabled: true                  # to trace incoming request headers on the core http server
        exclude:                       # to exclude specific routes from tracing
          - /healthz     
          - /livez     
          - /readyz     
          - /metrics     
      metrics:     
        expose: true                   # to expose metrics route, disabled by default
        path: /metrics                 # metrics route path (default /metrics)
        collect:       
          enabled: true                # to collect core http server metrics, disabled by default
          namespace: foo               # core http server metrics namespace (empty by default)
        buckets: 0.1, 1, 10            # to override default request duration buckets
        normalize:
          request_path: true           # to normalize http request path, disabled by default
          response_status: true        # to normalize http response status code (2xx, 3xx, ...), disabled by default
      healthcheck:
        startup:
          expose: true                 # to expose health check startup route, disabled by default
          path: /healthz               # health check startup route path (default /healthz)
        readiness:            
          expose: true                 # to expose health check readiness route, disabled by default
          path: /readyz                # health check readiness route path (default /readyz)
        liveness:            
          expose: true                 # to expose health check liveness route, disabled by default
          path: /livez                 # health check liveness route path (default /livez)
      tasks:
        expose: true                   # to expose tasks route, disabled by default
        path: /tasks/:name             # tasks route path (default /tasks/:name)  
      debug:
        config:
          expose: true                 # to expose debug config route
          path: /debug/config          # debug config route path (default /debug/config)
        pprof:
          expose: true                 # to expose debug pprof route
          path: /debug/pprof           # debug pprof route path (default /debug/pprof)
        routes:
          expose: true                 # to expose debug routes route
          path: /debug/routes          # debug routes route path (default /debug/routes)
        stats:
          expose: true                 # to expose debug stats route
          path: /debug/stats           # debug stats route path (default /debug/stats)
        build:
          expose: true                 # to expose debug build route
          path: /debug/build           # debug build route path (default /debug/build)
        modules:
          expose: true                 # to expose debug modules route
          path: /debug/modules/:name   # debug modules route path (default /debug/modules/:name)      
```

Notes:

- the core HTTP server requests logging will be based on the [log](fxlog.md) module configuration
- the core HTTP server requests tracing will be based on the [trace](fxtrace.md) module configuration
- if `app.debug=true` (or env var `APP_DEBUG=true`):
    - the dashboard will be automatically enabled
    - all the debug endpoints will be automatically exposed
    - error responses will not be obfuscated and stack trace will be added

## Usage

### Bootstrap

When you use a Yokai application template, a `internal/bootstrap.go` file is provided.

This is where you can:

- load Yokai `built-in`, [contrib](https://github.com/ankorstore/yokai-contrib) or your `own` modules
- configure the application with any `fx.Option`, at bootstrap on runtime

Example of bootstrap loading the [HTTP server](fxhttpserver.md) module:

```go title="internal/bootstrap.go"
package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

func init() {
	RootDir = fxcore.RootDir(1)
}

// RootDir is the application root directory.
var RootDir string

// Bootstrapper can be used to load modules, options, dependencies, routing and bootstraps your application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxhttpserver.FxHttpServerModule,
	// dependencies registration
	Register(),
	// routing registration
	Router(),

)

// Run starts the application, with a provided [context.Context].
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp()
}

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	tb.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", RootDir))

	Bootstrapper.RunTestApp(tb, fx.Options(options...))
}
```

Notes:

- the `Run()` function is used to start your application.
- the `RunTest()` function can be used in your tests, to start your application in test mode

### Dependency injection

Yokai is built on top of [Fx](https://github.com/uber-go/fx), offering a simple yet powerful `dependency injection system`.

This means you don't have to worry about injecting dependencies to your structs, your just need to register their constructors, and Yokai will automatically autowire them at runtime.

For example, if you create an `ExampleService` that has the [*config.Config](fxconfig.md) as dependency:

```go title="internal/service/example.go"
package service

import (
	"fmt"

	"github.com/ankorstore/yokai/config"
)

type ExampleService struct {
	config *config.Config
}

func NewExampleService(config *config.Config) *ExampleService {
	return &ExampleService{
		config: config,
	}
}

func (s *ExampleService) PrintAppName() {
	fmt.Printf("name: %s", s.config.AppName())
}
```

You then need to register it, by providing its constructor in `internal/register.go`:

```go title="internal/register.go"
package internal

import (
	"github.com/foo/bar/internal/service"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the ExampleService
		fx.Provide(service.NewExampleService),
		// ...
	)
}
```

This will make the `ExampleService` available in Yokai's dependency injection system, with its dependency on `*config.Config` autowired.

The `ExampleService` will also be available for injection in any constructor depending on it.

## Dashboard

If `modules.core.server.dashboard=true`, the core dashboard is available on the port `8081`:

![](../../assets/images/dash-core-light.png#only-light)
![](../../assets/images/dash-core-dark.png#only-dark)

Since it's served on a dedicated port, you can safely decide to
leave it enabled on production, to not expose it to the public, and access it
via [port forward](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/).

### Core

The `Core` section of the dashboard offers you information about:

- `Build`: environment and Go information about your application
- `Config`: resolved configuration
- `Metrics`: exposed metrics
- `Routes`: routes of the core dashboard
- `Pprof`: pprof page
- `Stats`: statistics page

### Health Check

The `Healthcheck` section of the dashboard offers you the possibility to trigger the health check endpoints, depending on their configuration.

You must ensure the health checks are exposed:

```yaml title="configs/config.yaml"
modules:
  core:
    server:
      healthcheck:
        startup:
          expose: true    # to expose health check startup route, disabled by default
          path: /healthz  # health check startup route path (default /healthz)
        readiness:
          expose: true    # to expose health check readiness route, disabled by default
          path: /readyz   # health check readiness route path (default /readyz)
        liveness:
          expose: true    # to expose health check liveness route, disabled by default
          path: /livez    # health check liveness route path (default /livez)
    
```

See the [Health Check](https://ankorstore.github.io/yokai/modules/fxhealthcheck/) module documentation for more information.

### Tasks

If you need to execute one shot / private operations (like flush a cache, trigger an export, etc.) but don't want to expose an endpoint or a command for this, you can create a task.

Yokai will collect them, and make them available in the core dashboard interface, under the `Tasks` section.

This is particularly useful for admin / maintenance purposes, without exposing those to your end users.

First, you must ensure the tasks are exposed:

```yaml title="configs/config.yaml"
modules:
  core:
    server:
      tasks:
        expose: true       # to expose tasks route, disabled by default
        path: /tasks/:name # tasks route path (default /tasks/:name)  
    
```

Then, provide a [Task](https://github.com/ankorstore/yokai/blob/main/fxcore/task.go) implementation:

```go title="internal/tasks/example.go"
package tasks

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcore"
)

var _ fxcore.Task = (*ExampleTask)(nil)

type ExampleTask struct {
	config *config.Config
}

func NewExampleTask(config *config.Config) *ExampleTask {
	return &ExampleTask{
		config: config,
	}
}

func (t *ExampleTask) Name() string {
	return "example"
}

func (t *ExampleTask) Run(ctx context.Context, input []byte) fxcore.TaskResult {
	return fxcore.TaskResult{
		Success: true,                     // task execution status
		Message: "example message",        // task execution message
		Details: map[string]any{           // optional task execution details
			"app":   t.config.AppName(),
			"input": string(input),
		},
	}
}
```

Then, register the task with `AsTask()`:

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/foo/bar/internal/tasks"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the ExampleTask (will auto wire dependencies)
		fxcore.AsTask(tasks.NewExampleTask),
		// ...
	)
}
```

Note: you can also use `AsTasks()` to register several tasks at once.

It'll be then available on the core dashboard for execution:

![](../../assets/images/dash-tasks-light.png#only-light)
![](../../assets/images/dash-tasks-dark.png#only-dark)

### Modules

The `Modules` section of the dashboard offers you the possibility to check the details of the modules exposing information to the core.

If you want your module to expose information in this section, you can provide a [FxModuleInfo](https://github.com/ankorstore/yokai/blob/main/fxcore/info.go) implementation:

```go title="internal/info.go"
package internal

type ExampleModuleInfo struct {}

func (i *ExampleModuleInfo) Name() string {
  return "example"
}

func (i *ExampleModuleInfo) Data() map[string]any {
  return map[string]any{
    "example": "value",
  }
}
```

and then register it in the `core-module-infos` group:

```go title="internal/register.go"
package internal

import (
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the ExampleModuleInfo in the core dashboard
		fx.Provide(
            fx.Annotate(
              ExampleModuleInfo,
              fx.As(new(interface{})),
              fx.ResultTags(`group:"core-module-infos"`),
            ),
		  ),
		// ...
	)
}
```

See [example](https://github.com/ankorstore/yokai/blob/main/fxhttpserver/info.go).


## Testing

You can start your application in test mode, with the `RunTest()` function provided in the [bootstrapper](#bootstrap).

This wil automatically set the env var `APP_ENV=test`, and [merge your test config](fxconfig.md#dynamic-env-overrides).

It accepts a list of `fx.Option`, for example:

- `fx.Populate()` to extract from the test application autowired components for your tests
- `fx.Invoke()` to execute a function at runtime
- `fx.Decorate()` to override components
- `fx.Replace()` to replace components
- etc.

Test example with `fx.Populate()` :

```go title="internal/example_test.go"
package internal_test

import (
	"testing"

	"github.com/foo/bar/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestExample(t *testing.T) {
	var exampleService *service.ExampleService

	// run app in test mode and extract the ExampleService
	internal.RunTest(t, fx.Populate(&exampleService))

	// assertion example
	assert.Equal(t, "foo", exampleService.Foo())
}
```

See [Fx documentation](https://pkg.go.dev/go.uber.org/fx) for the available `fx.Option`.