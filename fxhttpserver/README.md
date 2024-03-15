# Fx Http Server Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxhttpserver-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxhttpserver-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxhttpserver)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxhttpserver)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxhttpserver)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxhttpserver)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxhttpserver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxhttpserver)](https://pkg.go.dev/github.com/ankorstore/yokai/fxhttpserver)

> [Fx](https://uber-go.github.io/fx/) module for [httpserver](https://github.com/ankorstore/yokai/tree/main/httpserver).

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
	* [Dependencies](#dependencies)
	* [Loading](#loading)
	* [Configuration](#configuration)
	* [Registration](#registration)
		* [Middlewares](#middlewares)
		* [Handlers](#handlers)
		* [Handlers groups](#handlers-groups)
	* [Templates](#templates)
	* [Override](#override)
	* [Testing](#testing)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxhttpserver
```

## Features

This module provides a http server to your Fx application with:

- automatic panic recovery
- automatic requests logging and tracing (method, path, duration, ...)
- automatic requests metrics (count and duration)
- possibility to register handlers, groups and middlewares
- possibility to render HTML templates

## Documentation

### Dependencies

This module is intended to be used alongside:

- the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module
- the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module
- the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module
- the [fxmetrics](https://github.com/ankorstore/yokai/tree/main/fxmetrics) module
- the [fxgenerate](https://github.com/ankorstore/yokai/tree/main/fxgenerate) module

### Loading

To load the module in your Fx application:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,         // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxhttpserver.FxHttpServerModule, // load the module
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
  http:
    server:
      port: 8080                      # http server port (default 8080)
      errors:
        obfuscate: false              # to obfuscate error messages on the http server responses
        stack: false                  # to add error stack trace to error response of the http server
      log:
        headers:                      # to log incoming request headers on the http server
          x-foo: foo                  # to log for example the header x-foo in the log field foo
          x-bar: bar
        exclude:                      # to exclude specific routes from logging
          - /foo
          - /bar
        level_from_response: true     # to use response status code for log level (ex: 500=error)
      trace:
        enabled: true                 # to trace incoming request headers on the http server
        exclude:                      # to exclude specific routes from tracing
          - /foo
          - /bar
      metrics:
        collect:
          enabled: true               # to collect http server metrics
          namespace: foo              # http server metrics namespace (empty by default)
          subsystem: bar              # http server metrics subsystem (empty by default)
        buckets: 0.1, 1, 10           # to override default request duration buckets
        normalize:               
          request_path: true          # to normalize http request path, disabled by default
          response_status: true       # to normalize http response status code (2xx, 3xx, ...), disabled by default
      templates:
        enabled: true                 # disabled by default
        path: templates/*.html        # templates path lookup pattern
```

Notes:

- the http server requests logging will be based on the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog)
  module configuration
- the http server requests tracing will be based on the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace)
  module configuration
- if `app.debug=true` (or env var `APP_DEBUG=true`), error responses will not be obfuscated and stack trace will be
  added

### Registration

This module offers the possibility to easily register handlers, groups and middlewares.

#### Middlewares

You can use the `AsMiddleware()` function to register global middlewares on your http server:

- you can provide any [Middleware](registry.go) interface implementation (will be autowired from Fx container)
- or any `echo.MiddlewareFunc`, for example any
  built-in [Echo middleware](https://echo.labstack.com/docs/category/middleware)

```go
package main

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

type SomeMiddleware struct {
	config *config.Config
}

func NewSomeMiddleware(config *config.Config) *SomeMiddleware {
	return &SomeMiddleware{
		config: config,
	}
}

func (m *SomeMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// request correlated log
			httpserver.CtxLogger(c).Info().Msg("in some middleware")

			// use injected dependency
			c.Response().Header().Add("app-name", m.config.AppName())

			return next(c)
		}
	}
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,         // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxhttpserver.FxHttpServerModule, // load the module
		fx.Provide(
			fxhttpserver.AsMiddleware(middleware.CORS(), fxhttpserver.GlobalUse), // register echo CORS middleware via echo.Use()
			fxhttpserver.AsMiddleware(NewSomeMiddleware, fxhttpserver.GlobalPre), // register and autowire the SomeMiddleware via echo.Pre()
		),
	).Run()
}
```

#### Handlers

You can use the `AsHandler()` function to register handlers and their middlewares on your http server:

- you can provide any [Handler](registry.go) interface implementation (will be autowired from Fx container)
- or any `echo.HandlerFunc`

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

type SomeMiddleware struct {
	config *config.Config
}

func NewSomeMiddleware(config *config.Config) *SomeMiddleware {
	return &SomeMiddleware{
		config: config,
	}
}

func (m *SomeMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// request correlated log
			httpserver.CtxLogger(c).Info().Msg("in some middleware")

			// use injected dependency
			c.Response().Header().Add("app-name", m.config.AppName())

			return next(c)
		}
	}
}

type SomeHandler struct {
	config *config.Config
}

func NewSomeHandler(config *config.Config) *SomeHandler {
	return &SomeHandler{
		config: config,
	}
}

func (h *SomeHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		// request correlated trace span
		ctx, span := httpserver.CtxTracer(c).Start(c.Request().Context(), "some span")
		defer span.End()

		// request correlated log
		log.CtxLogger(ctx).Info().Msg("in some handler")

		// use injected dependency
		return c.String(http.StatusOK, fmt.Sprintf("app name: %s", h.config.AppName()))
	}
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,         // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxhttpserver.FxHttpServerModule, // load the module
		fx.Provide(
			// register and autowire the SomeHandler handler for [GET] /some-path, with the SomeMiddleware and echo CORS() middlewares
			fxhttpserver.AsHandler("GET", "/some-path", NewSomeHandler, NewSomeMiddleware, middleware.CORS()),
		),
	).Run()
}
```

#### Handlers groups

You can use the `AsHandlersGroup()` function to register handlers groups and their middlewares on your http
server:

- you can provide any [Handler](registry.go) interface implementation (will be autowired from Fx container) or
  any `echo.HandlerFunc`, with their middlewares
- and group them
	- under a common route `prefix`
	- with common [Middleware](registry.go) interface implementation (will be autowired from Fx container) or
	  any `echo.MiddlewareFunc`

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

type SomeMiddleware struct {
	config *config.Config
}

func NewSomeMiddleware(config *config.Config) *SomeMiddleware {
	return &SomeMiddleware{
		config: config,
	}
}

func (m *SomeMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// request correlated log
			httpserver.CtxLogger(c).Info().Msg("in some middleware")

			// use injected dependency
			c.Response().Header().Add("app-name", m.config.AppName())

			return next(c)
		}
	}
}

type SomeHandler struct {
	config *config.Config
}

func NewSomeHandler(config *config.Config) *SomeHandler {
	return &SomeHandler{
		config: config,
	}
}

func (h *SomeHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		// request correlated trace span
		ctx, span := httpserver.CtxTracer(c).Start(c.Request().Context(), "some span")
		defer span.End()

		// request correlated log
		log.CtxLogger(ctx).Info().Msg("in some handler")

		// use injected dependency
		return c.String(http.StatusOK, fmt.Sprintf("app name: %s", h.config.AppName()))
	}
}

type OtherHandler struct {
	config *config.Config
}

func NewOtherHandler(config *config.Config) *OtherHandler {
	return &OtherHandler{
		config: config,
	}
}

func (h *OtherHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		// use injected dependency
		return c.String(http.StatusOK, fmt.Sprintf("app version: %s", h.config.AppVersion()))
	}
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,         // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxhttpserver.FxHttpServerModule, // load the module
		fx.Provide(
			// register and autowire the SomeHandler handler with NewSomeMiddleware middleware for [GET] /group/some-path
			// register and autowire the OtherHandler handler with echo CORS middleware for [POST] /group/other-path
			// register the echo CSRF middleware for all handlers of this group
			fxhttpserver.AsHandlersGroup(
				"/group",
				[]*fxhttpserver.HandlerRegistration{
					fxhttpserver.NewHandlerRegistration("GET", "/some-path", NewSomeHandler, NewSomeMiddleware),
					fxhttpserver.NewHandlerRegistration("POST", "/other-path", NewOtherHandler, middleware.CORS()),
				},
				middleware.CSRF(),
			),
		),
	).Run()
}
```

### Templates

The module will look up HTML templates to render if `modules.http.server.templates.enabled=true`.

The HTML templates will be loaded from a path matching the pattern specified in `modules.http.server.templates.path`.

Considering the configuration:

```yaml
# ./configs/config.yaml
app:
  name: app
modules:
  http:
    server:
      templates:
        enabled: true
        path: templates/*.html
```

And the template:

```html
<!-- templates/app.html -->
<html>
    <body>
        <h1>App name is {{index . "name"}}</h1>
    </body>
</html>
```

To render it:

```go
package main

import (
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type TemplateHandler struct {
	config *config.Config
}

func NewTemplateHandler(cfg *config.Config) *TemplateHandler {
	return &TemplateHandler{
		config: cfg,
	}
}

func (h *TemplateHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		// output: "App name is app"
		return c.Render(http.StatusOK, "app.html", map[string]interface{}{
			"name": h.config.AppName(),
		})
	}
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,         // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxhttpserver.FxHttpServerModule, // load the module
		fx.Provide(
			fxhttpserver.AsHandler("GET", "/app", NewTemplateHandler),
		),
	).Run()
}
```

### Override

By default, the `echo.Echo` is created by
the [DefaultHttpServerFactory](https://github.com/ankorstore/yokai/blob/main/httpserver/factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type CustomHttpServerFactory struct{}

func NewCustomHttpServerFactory() httpserver.HttpServerFactory {
	return &CustomHttpServerFactory{}
}

func (f *CustomHttpServerFactory) Create(options ...httpserver.HttpServerOption) (*echo.Echo, error) {
	return echo.New(), nil
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,                 // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxhttpserver.FxHttpServerModule,         // load the module
		fx.Decorate(NewCustomHttpServerFactory), // override the module with a custom factory
		fx.Invoke(func(httpServer *echo.Echo) {  // invoke the custom http server
			// ...
		}),
	).Run()
}
```

### Testing

This module allows you to easily provide `functional` tests for your handlers.

For example, considering this handler:

```go
package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SomeHandler struct{}

func NewSomeHandler() *SomeHandler {
	return &SomeHandler{}
}

func (h *SomeHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}
}
```

You can then test it, considering `logs`, `traces` and `metrics` are enabled:

```go
package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"handler"
)

func TestSomeHandler(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxhttpserver.FxHttpServerModule,
		fx.Options(
			fxhttpserver.AsHandler("GET", "/test", handler.NewSomeHandler),
		),
		fx.Populate(&httpServer, &logBuffer, &traceExporter, &metricsRegistry), // extract components
	).RequireStart().RequireStop()

	// http call [GET] /test on the server
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	// assertions on http response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Body.String(), "ok")

	// assertion on the logs buffer
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test",
		"module":  "httpserver",
		"method":  "GET",
		"uri":     "/test",
		"status":  200,
		"message": "request logger",
	})

	// assertion on the traces exporter
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /test",
		semconv.HTTPRoute("/test"),
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPStatusCode(http.StatusOK),
	)

	// assertion on the metrics registry
	expectedMetric := `
		# HELP app_httpserver_requests_total Number of processed HTTP requests
		# TYPE app_httpserver_requests_total counter
		app_httpserver_requests_total{path="/test",method="GET",status="2xx"} 1
	`

	err := testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"app_httpserver_requests_total",
	)
	assert.NoError(t, err)
}
```

You can find more tests examples in this module own [tests](module_test.go).
