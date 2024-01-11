# Http Server Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/httpserver-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/httpserver-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/httpserver)](https://goreportcard.com/report/github.com/ankorstore/yokai/httpserver)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=httpserver)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/httpserver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/httpserver)](https://pkg.go.dev/github.com/ankorstore/yokai/httpserver)


> Http server module based on [Echo](https://echo.labstack.com/).

<!-- TOC -->

* [Installation](#installation)
* [Documentation](#documentation)
	* [Usage](#usage)
	* [Add-ons](#add-ons)
		* [Logger](#logger)
		* [Error handler](#error-handler)
		* [Http Handlers](#http-handlers)
			* [Debug handlers](#debug-handlers)
			* [Pprof handlers](#pprof-handlers)
			* [Healthcheck handlers](#healthcheck-handlers)
		* [Middlewares](#middlewares)
			* [Request id middleware](#request-id-middleware)
			* [Request logger middleware](#request-logger-middleware)
			* [Request tracer middleware](#request-tracer-middleware)
			* [Request metrics middleware](#request-metrics-middleware)
		* [HTML Templates](#html-templates)

<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/httpserver
```

## Documentation

### Usage

This module provides a [HttpServerFactory](factory.go), allowing to build an `echo.Echo` instance.

```go
package main

import (
	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var server, _ = httpserver.NewDefaultHttpServerFactory().Create()

// equivalent to:
var server, _ = httpserver.NewDefaultHttpServerFactory().Create(
	httpserver.WithDebug(false),                                  // debug disabled by default
	httpserver.WithBanner(false),                                 // banner disabled by default
	httpserver.WithRecovery(true),                                // panic recovery middleware enabled by default
	httpserver.WithLogger(log.New("default")),                    // echo default logger
	httpserver.WithBinder(&echo.DefaultBinder{}),                 // echo default binder
	httpserver.WithJsonSerializer(&echo.DefaultJSONSerializer{}), // echo default json serializer
	httpserver.WithHttpErrorHandler(nil),                         // echo default error handler
)

server.Start(...)
```

See [Echo documentation](https://echo.labstack.com/docs) for more details.

### Add-ons

This module provides several add-ons ready to use to enrich your http server.

#### Logger

This module provides an [EchoLogger](logger.go), compatible with
the [log module](https://github.com/ankorstore/yokai/tree/main/log):

```go
package main

import (
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
)

func main() {
	logger, _ := log.NewDefaultLoggerFactory().Create()

	server, _ := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
	)
}
```

#### Error handler

This module provides a [JsonErrorHandler](error.go), with configurable error obfuscation and call stack:

```go
package main

import (
	"github.com/ankorstore/yokai/httpserver"
)

func main() {
	server, _ := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(
			false, // without error details obfuscation
			false, // without error call stack
		)),
	)
}
```

You can set the parameters:

- `obfuscate=true` to obfuscate the error details from the response message, i.e. will use for
  example `Internal Server Error` for a response code 500 (recommended for production)
- `stack=true` to add the error call stack to the log and response (not suitable for production)

This will make a call to `[GET] https://example.com` and forward automatically the `authorization`, `x-request-id`
and `traceparent` headers from the handler request.

#### Http Handlers

##### Debug handlers

This module provides several [debug handlers](handler), compatible with
the [config module](https://github.com/ankorstore/yokai/tree/main/config):

- `DebugBuildHandler` to dump current build information
- `DebugConfigHandler` to dump current config values
- `DebugRoutesHandler` to dump current registered routes on the server
- `DebugVersionHandler` to dump current version

```go
package main

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/handler"
)

func main() {
	cfg, _ := config.NewDefaultConfigFactory().Create()

	server, _ := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithDebug(true),
	)

	if server.Debug {
		server.GET("/debug/build", handler.DebugBuildHandler())
		server.GET("/debug/config", handler.DebugConfigHandler(cfg))
		server.GET("/debug/routes", handler.DebugRoutesHandler(server))
		server.GET("/debug/version", handler.DebugVersionHandler(cfg))
	}
}
```

This will expose `[GET] /debug/*` endpoints (not suitable for production).

##### Pprof handlers

This module provides [pprof handlers](handler), compatible with the [net/http/pprof](https://pkg.go.dev/net/http/pprof)
package:

- `PprofIndexHandler` to offer pprof index dashboard
- `PprofAllocsHandler` to provide a sampling of all past memory allocations
- `PprofBlockHandler` to provide stack traces that led to blocking on synchronization primitives
- `PprofCmdlineHandler` to provide the command line invocation of the current program
- `PprofGoroutineHandler` to provide the stack traces of all current goroutines
- `PprofHeapHandler` to provide a sampling of memory allocations of live objects
- `PprofMutexHandler` to provide stack traces of holders of contended mutexes
- `PprofProfileHandler` to provide CPU profile
- `PprofSymbolHandler` to look up the program counters listed in the request
- `PprofThreadCreateHandler` to provide stack traces that led to the creation of new OS threads
- `PprofTraceHandler` to provide a trace of execution of the current program

```go
package main

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/handler"
)

func main() {
	server, _ := httpserver.NewDefaultHttpServerFactory().Create()

	// index dashboard
	server.GET("/debug/pprof/", handler.PprofIndexHandler())

	// linked from index dashboard
	server.GET("/debug/pprof/allocs", handler.PprofAllocsHandler())
	server.GET("/debug/pprof/block", handler.PprofBlockHandler())
	server.GET("/debug/pprof/cmdline", handler.PprofCmdlineHandler())
	server.GET("/debug/pprof/goroutine", handler.PprofGoroutineHandler())
	server.GET("/debug/pprof/heap", handler.PprofHeapHandler())
	server.GET("/debug/pprof/mutex", handler.PprofMutexHandler())
	server.GET("/debug/pprof/profile", handler.PprofProfileHandler())
	server.GET("/debug/pprof/symbol", handler.PprofSymbolHandler())
	server.POST("/debug/pprof/symbol", handler.PprofSymbolHandler())
	server.GET("/debug/pprof/threadcreate", handler.PprofThreadCreateHandler())
	server.GET("/debug/pprof/trace", handler.PprofTraceHandler())
}
```

This will expose pprof index dashboard on `[GET] /debug/pprof/`, from where you'll be able to retrieve all pprof
profiles types.

##### Healthcheck handlers

This module provides a [HealthCheckHandler](handler/health_check.go), compatible with
the [healthcheck module](https://github.com/ankorstore/yokai/tree/main/healthcheck):

```go
package main

import (
	"probes"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/handler"
)

func main() {
	checker, _ := healthcheck.NewDefaultCheckerFactory().Create(
		healthcheck.WithProbe(probes.SomeProbe()),                            // register for startup, liveness and readiness checks           
		healthcheck.WithProbe(probes.SomeOtherProbe(), healthcheck.Liveness), // register liveness checks only
	)

	server, _ := httpserver.NewDefaultHttpServerFactory().Create()

	server.GET("/healthz", handler.HealthCheckHandler(checker, healthcheck.Startup))
	server.GET("/livez", handler.HealthCheckHandler(checker, healthcheck.Liveness))
	server.GET("/readyz", handler.HealthCheckHandler(checker, healthcheck.Readiness))
}
```

This will expose endpoints
for [k8s startup, readiness and liveness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/):

- `[GET] /healthz`: startup probes checks
- `[GET] /livez`: liveness probes checks
- `[GET] /readyz`: readiness probes checks

#### Middlewares

##### Request id middleware

This module provides a [RequestIdMiddleware](middleware/request_id.go), ensuring the request and response will always
have a request id (coming by default from the `X-Request-Id` header or generated if missing) for correlation needs.

```go
package main

import (
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/middleware"
)

func main() {
	server, _ := httpserver.NewDefaultHttpServerFactory().Create()

	server.Use(middleware.RequestIdMiddleware())
}
```

If you need, you can configure the request header name it fetches the id from, or the generator used for missing id
generation:

```go
import (
	"github.com/ankorstore/yokai/generate/generatetest/uuid"
)

server.Use(middleware.RequestIdMiddlewareWithConfig(middleware.RequestIdMiddlewareConfig{
	RequestIdHeader: "custom-header",
	Generator: uuid.NewTestUuidGenerator("some-value"),
}))
```

##### Request logger middleware

This module provides a [RequestLoggerMiddleware](middleware/request_logger.go):

- compatible with the [log module](https://github.com/ankorstore/yokai/tree/main/log)
- ensuring all log entries will contain the requests `x-request-id` header value by default, in the field `requestID`,
  for correlation
- ensuring a recap log entry will be emitted at request completion

You can then use the [CtxLogger](context.go) method to access the correlated logger from with your handlers:

```go
package main

import (
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
)

func main() {
	logger, _ := log.NewDefaultLoggerFactory().Create()

	server, _ := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
	)

	server.Use(middleware.RequestLoggerMiddleware())

	// handler
	server.GET("/test", func(c echo.Context) error {
		// emit correlated log
		httpserver.CtxLogger(c).Info().Msg("info")

		// equivalent to
		log.CtxLogger(c.Request().Context()).Info().Msg("info")
	})
}
```

By default, the middleware logs all requests with `info` level, even if failed. If needed, you can configure it to log
with a level matching the response (or http error) code:

- `code < 400`: log level `info`
- `400 <= code < 500`: log level `warn`
- `code >= 500` or `non http error`: log level `error`

```go
server.Use(middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
	LogLevelFromResponseOrErrorCode: true,
}))
```

You can configure additional request headers to log:

- the key is the header name to fetch
- the value is the log field name to fill

```go
server.Use(middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
	RequestHeadersToLog: map[string]string{
		"x-header-foo": "foo",
		"x-header-bar": "bar",
	},
}))
```

You can also configure the request URI prefixes to exclude from logging:

```go
server.Use(middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
	RequestUriPrefixesToExclude: []string{
		"/foo",
		"/bar",
	},
}))
```

Note: if a request to an excluded URI fails (error or http code >= 500), the middleware will still log for observability
purposes.

##### Request tracer middleware

This module provides a [RequestTracerMiddleware](middleware/request_tracer.go):

- using the global tracer by default
- compatible with the [trace module](https://github.com/ankorstore/yokai/tree/main/trace)
- ensuring a recap trace span will be emitted at request completion

You can then use, from within your handlers the [CtxTracer](context.go) method to access the correlated tracer:

```go
package main

import (
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/ankorstore/yokai/trace"
	"github.com/labstack/echo/v4"
)

func main() {
	server, _ := httpserver.NewDefaultHttpServerFactory().Create()

	server.Use(middleware.RequestTracerMiddleware("my-service"))

	// handler
	server.GET("/test", func(c echo.Context) error {
		// emit correlated span
		_, span := httpserver.CtxTracer(c).Start(c.Request().Context(), "my-span")
		defer span.End()

		// equivalent to
		_, span = trace.CtxTracerProvider(c.Request().Context()).Tracer("my-tracer").Start(c.Request().Context(), "my-span")
		defer span.End()
	})
}
```

If you need, you can configure the tracer provider and propagators:

```go
import (
	"github.com/ankorstore/yokai/trace"
	"go.opentelemetry.io/otel/propagation"
)

tracerProvider, _ := trace.NewDefaultTracerProviderFactory().Create()

server.Use(middleware.RequestTracerMiddlewareWithConfig("my-service", middleware.RequestTracerMiddlewareConfig{
	TracerProvider: tracerProvider,
	TextMapPropagator: propagation.TraceContext{},
}))
```

And you can also configure the request URI prefixes to exclude from tracing:

```go
server.Use(middleware.RequestTracerMiddlewareWithConfig("my-service", middleware.RequestTracerMiddlewareConfig{
	RequestUriPrefixesToExclude: []string{"/test"},
}))
```

##### Request metrics middleware

This module provides a [RequestMetricsMiddleware](middleware/request_metrics.go):

- ensuring requests processing count and duration are collected
- using the global `promauto` metrics registry by default

```go
package main

import (
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	server, _ := httpserver.NewDefaultHttpServerFactory().Create()

	server.Use(middleware.RequestMetricsMiddleware())

	// handler
	server.GET("/test", func(c echo.Context) error {
		// ...
	})
}
```

If you need, you can configure the metrics registry, namespace, subsystem, buckets and status code normalization:

```go
import (
	"github.com/prometheus/client_golang/prometheus"
)

registry := prometheus.NewPedanticRegistry()

server.Use(middleware.RequestMetricsMiddlewareWithConfig(middleware.RequestMetricsMiddlewareConfig{
	Registry:            registry,
	Namespace:           "foo",
	Subsystem:           "bar",
	Buckets:             []float64{0.01, 1, 10},
	NormalizeHTTPStatus: true,
}))
```

#### HTML Templates

This module provides a [HtmlTemplateRenderer](renderer.go) for rendering HTML templates.

Considering the following template:

```html
<!-- path/to/templates/welcome.html -->
<html>
	<body>
		<h1>Welcome {{index . "name"}}!</h1>
	</body>
</html>
```

To render it:

```go
package main

import (
	"net/http"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
)

func main() {
	server, _ := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithRenderer(httpserver.NewHtmlTemplateRenderer("path/to/templates/*.html")), // templates lookup pattern
	)

	// handler
	server.GET("/welcome", func(c echo.Context) error {
		return c.Render(http.StatusOK, "welcome.html", map[string]interface{}{
			"name": "some name",
		})
	})
}
```

See [Echo templates documentation](https://echo.labstack.com/docs/templates) for more details.
