# Fx Core Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxcore-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxcore-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxcore)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxcore)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxcore)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxcore)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxcore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxcore)](https://pkg.go.dev/github.com/ankorstore/yokai/fxcore)

> [Fx](https://uber-go.github.io/fx/) core module.

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
	* [Preloaded modules](#preloaded-modules)
	* [Configuration](#configuration)
	* [Bootstrap](#bootstrap)
		* [Application](#application)
		* [Test application](#test-application)
		* [Root dir](#root-dir)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxcore
```

## Features

The fxcore module provides the foundation of your application:

- a bootstrapper
- a dependency injection system
- a dedicated core http server
- ready to use config, health check, logger and tracer and metrics components
- an extension system for Yokai built-in, [contrib](https://github.com/ankorstore/yokai-contrib) or your own modules

The `core http server` runs automatically on a dedicated port (default `8081`), to serve:

- the dashboard: UI to get an overview of your application
- the metrics endpoint: to expose all collected metrics from your application
- the health check endpoints: to expose all configured health check probes of your application
- the debug endpoints: to expose various information about your config, modules, build, etc.

Whatever your type of application (httpserver, gRPC server, worker, etc.), all `platform concerns` are handled by this
dedicated server:

- to avoid to expose sensitive information (health checks, metrics, debug, etc) to your users
- and most importantly to enable your application to focus on its logic

## Documentation

### Preloaded modules

This core module preloads:

- the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module
- the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module
- the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module
- the [fxmetrics](https://github.com/ankorstore/yokai/tree/main/fxmetrics) module
- the [fxgenerate](https://github.com/ankorstore/yokai/tree/main/fxgenerate) module
- the [fxhealthcheck](https://github.com/ankorstore/yokai/tree/main/fxhealthcheck) module

### Configuration

Configuration reference:

```yaml
# ./configs/config.yaml
app:
  name: app
  description: app description
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

- the core http server requests logging will be based on the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module configuration
- the core http server requests tracing will be based on the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module configuration
- if `app.debug=true` (or env var `APP_DEBUG=true`):
	- the dashboard will be automatically enabled
    - all the debug endpoints will be automatically exposed
	- error responses will not be obfuscated and stack trace will be added

Check the [configuration files documentation](https://github.com/ankorstore/yokai/tree/main/config#configuration-files) for more details.

### Bootstrap

The core module provides a bootstrapper:

- to plug in all the [Fx modules](https://github.com/ankorstore/yokai#fx-modules) required by your application
- to provide your own application modules and services
- to start your application (real or test runtime)

#### Application

Create an application service, for example depending on a database connection:

```go
package service

import (
	"gorm.io/gorm"
)

type ExampleService struct {
	db *gorm.DB
}

func NewExampleService(db *gorm.DB) *ExampleService {
	return &ExampleService{
		db: db,
	}
}

func (s *ExampleService) Ping() bool {
	return s.db.Ping() // simplification
}
```

Create your application [Bootstrapper](bootstrap.go) with your bootstrap options:

```go
package bootstrap

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxorm"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"path/to/service"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	fxorm.FxOrmModule,                     // load the ORM module (provides *gorm.DB)
	fx.Provide(service.NewExampleService), // autowire your service (*gorm.DB auto injection)
	fxcore.AsCoreExtraInfo("foo", "bar"),  // register extra information to display on core dashboard
)
```

You can use the bootstrapper to start your application:

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/fxcore"
	"path/to/bootstrap"
)

func main() {
	// run the application
	bootstrap.Bootstrapper.RunApp()

	// or you can also run the application with a specific root context
	bootstrap.Bootstrapper.WithContext(context.Background()).RunApp()

	// or you can also bootstrap and run it on your own
	app := bootstrap.Bootstrapper.BootstrapApp()
	app.Run()
}
```

#### Test application

You can reuse your [Bootstrapper](bootstrap.go) to run your application in test mode:

```go
package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"path/to/bootstrap"
	"path/to/service"
)

func TestExampleService(t *testing.T) {
	// *service.ExampleService instance to extract from your application
	var svc *service.ExampleService

	// run the app in test mode and populate the service
	bootstrap.Bootstrapper.RunTestApp(t, fx.Populate(&svc))

	// assertion example
	assert.True(t, svc.Ping())
}
```

You can also use `BootstrapTestApp()` to bootstrap in test mode and run it on your own:

```go
testApp := bootstrap.Bootstrapper.BootstrapTestApp(t, ...)
testApp.RequireStart().RequireStop()
```

Note: bootstrapping your application in test mode will set `APP_ENV=test`, automatically loading your testing
configuration.

#### Root dir

The core module provides the possibility to retrieve the root dir with `RootDir()`, useful for setting relative
path to templates or configs.

```go
package bootstrap

import (
	"github.com/ankorstore/yokai/fxcore"
)

var RootDir string

func init() {
	RootDir = fxcore.RootDir(0) // configure number of stack frames to ascend
}
```

Then you can then use the global `RootDir` variable in any packages:

```go
package main

import (
	"fmt"

	"path/to/bootstrap"
)

func main() {
	fmt.Printf("root dir: %s", bootstrap.RootDir)
}
```

Or in any tests:

```go
package main_test

import (
	"fmt"
	"testing"

	"path/to/bootstrap"
)

func TestSomething(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", bootstrap.RootDir))

	//...
}
```
