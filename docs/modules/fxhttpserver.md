# HTTP Server Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxhttpserver-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxhttpserver-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxhttpserver)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxhttpserver)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxhttpserver)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxhttpserver)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxhttpserver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxhttpserver)](https://pkg.go.dev/github.com/ankorstore/yokai/fxhttpserver)

## Overview

Yokai provides a [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module, providing an HTTP server to your application.

It wraps the [httpserver](https://github.com/ankorstore/yokai/tree/main/httpserver) module, based on [Echo](https://echo.labstack.com/).

It comes with:

- automatic panic recovery
- automatic requests logging and tracing (method, path, duration, ...)
- automatic requests metrics (count and duration)
- possibility to register handlers, groups and middlewares
- possibility to render HTML templates

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxhttpserver
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpserver"
)


var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxhttpserver module
	fxhttpserver.FxHttpServerModule,
	// routing
	ProvideRouting(),
	// ...
)
```

Then create, if not existing, the `internal/routing.go` file for your registrations:

```go title="internal/routing.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

func ProvideRouting() fx.Option {
	return fx.Options()
}
```

It is recommended to keep routing definition separated from services definitions, for better maintainability. If you use the [HTTP application template](../../applications/templates/#http-application-template), this is already done for you.

## Usage

This module offers the possibility to easily register HTTP handlers, groups and middlewares.

### Middleware registration

You can use the `AsMiddleware()` function to register global middlewares on your HTTP server:

- any [Middleware](https://github.com/ankorstore/yokai/blob/main/fxhttpserver/registry.go) implementation
- or any `echo.MiddlewareFunc`, for example [Echo built-in middlewares]([echo.MiddlewareFunc](https://echo.labstack.com/docs/category/middleware))

For example, you can create a middleware:

```go title="internal/middleware/example.go"
package middleware

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type ExampleMiddleware struct {
	config *config.Config
}

func NewExampleMiddleware(config *config.Config) *ExampleMiddleware {
	return &ExampleMiddleware{
		config: config,
	}
}

func (m *ExampleMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// example of correlated log
			log.CtxLogger(c.Request().Context()).Info().Msg("in example middleware")

			// injected dependency example usage
			c.Response().Header().Add("app-name", m.config.AppName())

			return next(c)
		}
	}
}
```

You can then register your middlewares:

```go title="internal/routing.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/foo/bar/internal/middleware"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

func ProvideRouting() fx.Option {
	return fx.Options(
		// registers the Echo CORS middleware via echo.Use()
		fxhttpserver.AsMiddleware(echomiddleware.CORS(), fxhttpserver.GlobalUse),
		// registers and autowire the ExampleMiddleware via echo.Pre()
		fxhttpserver.AsMiddleware(middleware.NewExampleMiddleware, fxhttpserver.GlobalPre), 
		// ...
	)
}
```

### Handler registration

You can use the `AsHandler()` function to register handlers and their middlewares on your HTTP server:

- any [Handler](https://github.com/ankorstore/yokai/blob/main/fxhttpserver/registry.go) implementation
- or any `echo.HandlerFunc`

For example, you can create a handler:

```go title="internal/handler/example.go"
package handler

import (
	"net/http"
	
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type ExampleHandler struct {
	config *config.Config
}

func NewExampleHandler(config *config.Config) *ExampleHandler {
	return &ExampleHandler{
		config: config,
	}
}

func (h *ExampleHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		// example of correlated trace span
		ctx, span := trace.CtxTracer(c.Request().Context()).Start(c.Request().Context(), "example span")
		defer span.End()

		// example of correlated log
		log.CtxLogger(ctx).Info().Msg("in example handler")

		// injected dependency example usage
		return c.String(http.StatusOK, fmt.Sprintf("app name: %s", h.config.AppName()))
	}
}
```

You can then register it:

```go title="internal/routing.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/foo/bar/internal/handler"
	"github.com/foo/bar/internal/middleware"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

func ProvideRouting() fx.Option {
	return fx.Options(
		// registers and autowire the ExampleHandler for [GET] /example, with the ExampleMiddleware and Echo CORS() middlewares
		fxhttpserver.AsHandler("GET", "/example", handler.NewExampleHandler, middleware.NewExampleMiddleware, echomiddleware.CORS()),
		// ...
	)
}
```

### Handlers group registration

You can use the `AsHandlersGroup()` function to register handlers groups and their middlewares on your HTTP
server:

- any [Handler](https://github.com/ankorstore/yokai/blob/main/fxhttpserver/registry.go) implementation or any `echo.HandlerFunc`, with their middlewares
- and group them
	- under a common route `prefix`
	- with common any [Middleware](https://github.com/ankorstore/yokai/blob/main/fxhttpserver/registry.go) implementations or any `echo.MiddlewareFunc`

For example, you can create another handler:

```go title="internal/handler/other.go"
package handler

import (
	"net/http"
	
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

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
		// example of correlated trace span
		ctx, span := trace.CtxTracer(c.Request().Context()).Start(c.Request().Context(), "other span")
		defer span.End()

		// example of correlated log
		log.CtxLogger(ctx).Info().Msg("in other handler")

		// injected dependency example usage
		return c.String(http.StatusOK, fmt.Sprintf("app name: %s", h.config.AppName()))
	}
}
```

You can then register your handlers in a group:

```go title="internal/routing.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/foo/bar/internal/handler"
	"github.com/foo/bar/internal/middleware"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

func ProvideRouting() fx.Option {
	return fx.Options(
		fxhttpserver.AsHandlersGroup(
			// common route prefix
			"/group",
			[]*fxhttpserver.HandlerRegistration{
				// registers and autowire the ExampleHandler for [GET] /group/example, with the ExampleMiddleware
				fxhttpserver.NewHandlerRegistration("GET", "/example", handler.NewExampleHandler, middleware.NewExampleMiddleware),
				// registers and autowire the OtherHandler for [GET] /group/other, with the Echo CORS middleware
				fxhttpserver.NewHandlerRegistration("GET", "/other", handler.NewOtherHandler, echomiddleware.CORS()),
			},
			// common Echo CSRF middleware, applied to both handlers
			echomiddleware.CSRF(),
		),
		// ...
	)
}
```

## Configuration

Configuration reference:

```yaml title="configs/config.yaml"
modules:
  http:
    server:
      port: 8080                  # http server port (default 8080)
      errors:
        obfuscate: false          # to obfuscate error messages on the http server responses
        stack: false              # to add error stack trace to error response of the http server
      log:
        headers:                  # to log incoming request headers on the http server
          x-foo: foo              # to log for example the header x-foo in the log field foo
          x-bar: bar
        exclude:                  # to exclude specific routes from logging
          - /foo
          - /bar
        level_from_response: true # to use response status code for log level (ex: 500=error)
      trace:
        enabled: true             # to trace incoming request headers on the http server
        exclude:                  # to exclude specific routes from tracing
          - /foo
          - /bar
      metrics:
        collect:
          enabled: true           # to collect http server metrics
          namespace: app          # http server metrics namespace (default app.name value)
          subsystem: httpserver   # http server metrics subsystem (default httpserver)
        buckets: 0.1, 1, 10       # to override default request duration buckets
        normalize: true           # to normalize http status code (2xx, 3xx, ...)
      templates:
        enabled: true             # disabled by default
        path: templates/*.html    # templates path lookup pattern
```


If `app.debug=true` (or env var `APP_DEBUG=true`), error responses will not be obfuscated and stack trace will be added.

## Templates

The module will look up HTML templates to render if `modules.http.server.templates.enabled=true`.

The HTML templates will be loaded from a path matching the pattern specified in `modules.http.server.templates.path`.

For example, considering the following configuration:

```yaml title="configs/config.yaml"
app:
  name: app
modules:
  http:
    server:
      templates:
        enabled: true
        path: templates/*.html
```

And the following template:

```html title="templates/app.html"
<html>
    <body>
        <h1>App name is {{index . "name"}}</h1>
    </body>
</html>
```

You can then render it from your handler, with the `Render()` function:

```go title="internal/handler/template.go"
package handler

import (
	"net/http"
	
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
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
		// will render (HTML output): "App name is app"
		return c.Render(http.StatusOK, "app.html", map[string]interface{}{
			"name": h.config.AppName(),
		})
	}
}
```

## Logging

You can configure HTTP requests automatic logging:

```yaml title="configs/config.yaml"
modules:
  http:
    server:
      log:
        headers:                  # to log incoming request headers on the http server
          x-foo: foo              # to log for example the header x-foo in the log field foo
          x-bar: bar
        exclude:                  # to exclude specific routes from logging
          - /foo
          - /bar
        level_from_response: true # to use response status code for log level (ex: 500=error)
```

To get logs correlation in your handlers, you need to retrieve the logger from the context with `log.CtxLogger()`:

```go
log.CtxLogger(c.Request().Context()).Info().Msg("example message")
```

As a result, in your application logs:

```
INT service=app example message requestID=0f507e36-ea56-4842-b2f5-a53467e227e5 spanID=950c48301f39d2e3 traceID=d69d972b00302ec3e5369c8d439c4fac
INF service=app request logger latency="12.34µs" method=GET uri=/example status=200 module=httpserver requestID=0f507e36-ea56-4842-b2f5-a53467e227e5 spanID=950c48301f39d2e3 traceID=d69d972b00302ec3e5369c8d439c4fac
```

If both HTTP server logging and tracing are enabled, log records will automatically have the current `traceID` and `spanID` to be able to correlate logs and trace spans.

The HTTP server logging will be based on the [fxlog](fxlog.md) module configuration.

## Tracing

You can enable HTTP requests automatic tracing with `modules.http.server.trace.enable=true`:

```yaml title="configs/config.yaml"
modules:
  http:
    server:
      trace:
        enabled: true # to trace incoming request headers on the http server
        exclude:      # to exclude specific routes from tracing
          - /foo
          - /bar
```

To get traces correlation in your handlers, you need to retrieve the tracer from the context with `trace.CtxTracer()`:

```go
ctx, span := trace.CtxTracer(c.Request().Context()).Start(c.Request().Context(), "example span")
defer span.End()
```

As a result, in your application trace spans attributes:

```
service.name: app
http.method: GET
http.route: /example
http.status_code: 200
...
```

## Metrics

You can enable HTTP requests automatic metrics with `modules.http.server.trace.enable=true`:

```yaml title="configs/config.yaml"
modules:
  http:
    server:
      metrics:
        collect:
          enabled: true          # to collect http server metrics
          namespace: app         # http server metrics namespace (default app.name value)
          subsystem: httpserver  # http server metrics subsystem (default httpserver)
        buckets: 0.1, 1, 10      # to override default request duration buckets
        normalize: true          # to normalize http status code (2xx, 3xx, ...)
```

For example, after calling `[GET] /example`, the [fxcore](https://github.com/ankorstore/yokai/tree/main/fxcore) HTTP server will expose in the configured metrics endpoint:

```makefile title="[GET] /metrics"
# ...
# HELP app_httpserver_request_duration_seconds Time spent processing HTTP requests
# TYPE app_httpserver_request_duration_seconds histogram
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="0.005"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="0.01"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="0.025"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="0.05"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="0.1"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="0.25"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="0.5"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="1"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="2.5"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="5"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="10"} 1
app_httpserver_request_duration_seconds_bucket{handler="/example",method="GET",le="+Inf"} 1
app_httpserver_request_duration_seconds_sum{handler="/",method="GET"} 0.0014433150000000001
# HELP app_httpserver_requests_total Number of processed HTTP requests
# TYPE app_httpserver_requests_total counter
app_httpserver_requests_total{handler="/example",method="GET",status="200"} 1
```

## Testing

This module provides the possibility to perform functional testing, by calling your application endpoints from your tests.

You can easily assert on:

- HTTP response
- logs
- traces
- metrics

For example, if we want to test the [ExampleHandler](#handler-registration):

```go title="internal/handler/example_test.go"
package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/foo/bar/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func TestExampleHandler(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter, &metricsRegistry))

	// call [GET] /example
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	// HTTP response example
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Body.String(), "app name: app")

	// logs assertion example
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "app",
		"message": "in example handler",
	})
	
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "app",
		"module":  "httpserver",
		"method":  "GET",
		"uri":     "/example",
		"status":  http.StatusOK,
		"message": "request logger",
	})

	// traces assertion example
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"example span",
	)
	
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /example",
		semconv.HTTPRoute("/test"),
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPStatusCode(http.StatusOK),
	)

	// metrics assertion example
	expectedMetric := `
		# HELP app_httpserver_requests_total Number of processed HTTP requests
		# TYPE app_httpserver_requests_total counter
		app_httpserver_requests_total{handler="/example",method="GET",status="2xx"} 1
	`

	err := testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"app_httpserver_requests_total",
	)
	assert.NoError(t, err)
}
```
