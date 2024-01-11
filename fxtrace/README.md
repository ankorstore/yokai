# Fx Trace Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxtrace-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxtrace-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxtrace)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxtrace)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxtrace)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxtrace)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxtrace)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxtrace)](https://pkg.go.dev/github.com/ankorstore/yokai/fxtrace)

> [Fx](https://uber-go.github.io/fx/) module for [trace](https://github.com/ankorstore/yokai/tree/main/trace).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
	* [Dependencies](#dependencies)
	* [Loading](#loading)
	* [Configuration](#configuration)
	* [Override](#override)
	* [Testing](#testing)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxtrace
```

## Documentation

### Dependencies

This module is intended to be used alongside the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig)
module.

### Loading

To load the module in your Fx application:

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxtrace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule, // load the module dependency
		fxtrace.FxTraceModule,   // load the module
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			// invoke the tracer provider to create a span
			_, span := tracerProvider.Tracer("some tracer").Start(context.Background(), "some span")
			defer span.End()
		}),
	).Run()
}
```

### Configuration

This module provides the possibility to configure the `processor`:

- `noop`: to async void traces (default and fallback)
- `stdout`: to async print traces to stdout
- `otlp-grpc`: to async send traces to [OTLP/gRPC](https://opentelemetry.io/docs/specs/otlp/#otlpgrpc) collectors (ex: [Jaeger](https://www.jaegertracing.io/), [Grafana](https://grafana.com/docs/tempo/latest/configuration/grafana-agent/#grafana-agent), etc.)
- `test`: to sync store traces in memory (for testing assertions)

If an error occurs while creating the processor (for example failing OTLP/gRPC connection), the `noop` processor will be
used as safety fallback (to prevent outages).

This module also provides possibility to configure the `sampler`:

- `parent-based-always-on`: always on depending on parent (default)
- `parent-based-always-off`: always off depending on parent
- `parent-based-trace-id-ratio`: trace id ratio based depending on parent
- `always-on`: always on
- `always-off`: always off
- `trace-id-ratio`: trace id ratio based

Example with `stdout` processor (with pretty print) and `parent-based-trace-id-ratio` sampler (ratio=0.5):

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  trace:
    processor:
      type: stdout
      options:
        pretty: true
    sampler:
      type: parent-based-trace-id-ratio
      options:
        ratio: 0.5
```

Another example with `otlp-grpc` processor (on jaeger:4317 host) and `always-on` sampler:

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  trace:
    processor:
      type: otlp-grpc
      options:
        host: jaeger:4317
    sampler:
      type: always-on
```

### Override

By default, the `oteltrace.TracerProvider` is created by the [DefaultTracerProviderFactory](https://github.com/ankorstore/yokai/blob/main/trace/factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/trace"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type CustomTracerProviderFactory struct{}

func NewCustomTracerProviderFactory() trace.TracerProviderFactory {
	return &CustomTracerProviderFactory{}
}

func (f *CustomTracerProviderFactory) Create(options ...trace.TracerProviderOption) (*otelsdktrace.TracerProvider, error) {
	return &otelsdktrace.TracerProvider{...}, nil
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,                     // load the module dependency
		fxtrace.FxTraceModule,                       // load the module
		fx.Decorate(NewCustomTracerProviderFactory), // override the module with a custom factory
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) { // invoke the custom tracer provider
			_, span := tracerProvider.Tracer("custom tracer").Start(context.Background(), "custom span")
			defer span.End()
		}),
	).Run()
}
```

### Testing

This module provides the possibility to easily test your trace spans, using the [TestTraceExporter](https://github.com/ankorstore/yokai/blob/main/trace/tracetest/exporter.go) with `modules.trace.processor.type=test`.

```yaml
# ./configs/config.test.yaml
modules:
  trace:
    processor:
      type: test # to send traces to test buffer
```

You can then test:

```go
package main_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestTracerProvider(t *testing.T) {
	t.Setenv("APP_NAME", "test")
	t.Setenv("APP_ENV", "test")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("some tracer").Start(
				context.Background(),
				"some span",
				oteltrace.WithAttributes(attribute.String("some attribute name", "some attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter), // extracts the TestTraceExporter from the Fx container
	).RequireStart().RequireStop()

	// assertion success
	tracetest.AssertHasTraceSpan(t, exporter, "some span", attribute.String("some attribute name", "some attribute value"))
}
```

See the `trace` module testing [documentation](https://github.com/ankorstore/yokai/tree/main/trace#test-span-processor) for more details.
