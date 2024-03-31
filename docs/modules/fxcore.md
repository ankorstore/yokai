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
- a dedicated HTTP server
- pre-enabled [config](fxconfig.md), [health check](fxhealthcheck.md), [log](fxlog.md), [trace](fxtrace.md), [metrics](fxmetrics.md) and [generate](fxgenerate.md) modules
- an extension system for Yokai `built-in`, [contrib](https://github.com/ankorstore/yokai-contrib) or your `own` modules

The core HTTP server runs automatically on a dedicated port (default `8081`), to serve:

- the dashboard: UI to get an overview of your application
- the debug endpoints: to expose information about your build, config, loaded modules, etc.
- the health check endpoints: to expose the configured [health check probes](fxhealthcheck.md#probes-registration) of your application
- the metrics endpoint: to expose all [collected metrics](fxmetrics.md#metrics-registration) from your application

Whatever your type of application (HTTP, gRPC, worker, etc.), all platform concerns are handled by this
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
      port: 8081                       # core http server port (default 8081)
      errors:              
        obfuscate: false               # to obfuscate error messages on the core http server responses
        stack: false                   # to add error stack trace to error response of the core http server
      dashboard:
        enabled: true                  # to enable the core dashboard
        overview:      
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
          request_path: true          # to normalize http request path, disabled by default
          response_status: true       # to normalize http response status code (2xx, 3xx, ...), disabled by default
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

// Bootstrapper can be used to load modules, options, services and bootstraps your application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// fxhttpserver module loading
	fxhttpserver.FxHttpServerModule,
	// routing
	ProvideRouting(),
	// services
	ProvideServices(),
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

Yokai is built on top of [Fx](https://github.com/uber-go/fx), offering a simple yet powerful dependency injection system.

This means you don't have to worry about injecting dependencies to your structs, your just need to register their constructors, and Yokai will automatically autowire them at runtime.

For example, if you create an `ExampleService` that has the [config](fxconfig.md) as dependency:

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

This will make the `ExampleService` available in Yokai's dependency injection system, with its dependencies autowired.

## Dashboard

The core dashboard is available on the port `8081` if `modules.core.server.dashboard=true`:

![](../../assets/images/dash-core-light.png#only-light)
![](../../assets/images/dash-core-dark.png#only-dark)

From there, you can get:

- an overview of your application
- information and tooling about your application: build, config, metrics, pprof, etc.
- access to the configured health check endpoints
- access to the loaded modules information (when exposed)

The core dashboard is made for development purposes, but since it's served on a dedicated port, you can safely decide to leave it enabled on production, not expose it to the public, and access it via [port forward](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) for example.

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

See [Fx documentation](https://pkg.go.dev/go.uber.org/fx) for available options.