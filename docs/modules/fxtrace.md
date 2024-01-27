# Trace Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxtrace-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxtrace-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxtrace)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxtrace)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxtrace)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxtrace)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxtrace)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxtrace)](https://pkg.go.dev/github.com/ankorstore/yokai/fxtrace)

## Overview

Yokai provides a [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module, allowing your application to produce traces.

It wraps the [trace](https://github.com/ankorstore/yokai/tree/main/trace) module, based on [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go).

## Installation

The [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module is automatically loaded by
the [fxcore](https://github.com/ankorstore/yokai/tree/main/fxcore).

When you use a Yokai [application template](https://ankorstore.github.io/yokai/applications/templates/), you have nothing to install, it's ready to use.

## Usage

This module makes available the [TracerProvider](https://github.com/open-telemetry/opentelemetry-go) in
Yokai dependency injection system.

It is built on top of `OpenTelemetry`, see its [documentation](https://github.com/open-telemetry/opentelemetry-go) for more details about available methods

You can inject the tracer provider where needed, but it's recommended to use the one carried by the `context.Context` when possible (for automatic traces correlation).

## Configuration

This module provides the possibility to configure a `processor`:

- `noop`: to async void traces (default and fallback)
- `stdout`: to async print traces to stdout
- `otlp-grpc`: to async send traces to [OTLP/gRPC](https://opentelemetry.io/docs/specs/otlp/#otlpgrpc) collectors (ex: [Jaeger](https://www.jaegertracing.io/), [Grafana](https://grafana.com/docs/tempo/latest/configuration/grafana-agent/#grafana-agent), etc.)
- `test`: to sync store traces in memory (for testing assertions)

If an error occurs while creating the processor (for example failing OTLP/gRPC connection), the `noop` processor will be
used as safety fallback (to prevent outages).

This module also provides possibility to configure a `sampler`:

- `parent-based-always-on`: always on depending on parent (default)
- `parent-based-always-off`: always off depending on parent
- `parent-based-trace-id-ratio`: trace id ratio based depending on parent
- `always-on`: always on
- `always-off`: always off
- `trace-id-ratio`: trace id ratio based

Example with `stdout` processor (with pretty print) and `parent-based-trace-id-ratio` sampler (ratio=0.5):

```yaml title="configs/config.yaml"
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

Another example with `otlp-grpc` processor (sending on jaeger:4317) and `always-on` sampler:

```yaml title="configs/config.yaml"
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

## Testing

This module provides the possibility to easily test your trace spans, using the [TestTraceExporter](https://github.com/ankorstore/yokai/blob/main/trace/tracetest/exporter.go) with `modules.trace.processor.type=test`.

```yaml title="configs/config.test.yaml"
modules:
  trace:
    processor:
      type: test # to send traces to test buffer
```

You can use the provided [test assertion helpers](https://github.com/ankorstore/yokai/blob/main/trace/tracetest/assert.goo) in your tests:

- `AssertHasTraceSpan`: to assert on exact name and exact attributes match
- `AssertHasNotTraceSpan`: to assert on exact name and exact attributes non match
- `AssertContainTraceSpan`: to assert on exact name and partial attributes match
- `AssertContainNotTraceSpan`: to assert on exact name and partial attributes non match

For example:

```go title="internal/example_test.go"
package internal_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/foo/bar/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func TestExample(t *testing.T) {
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(
		t,
		fx.Populate(&traceExporter),
		fx.Invoke(func(tracerProvider trace.TracerProvider) {
			_, span := tracerProvider.Tracer("example tracer").Start(
				context.Background(),
				"example span",
				trace.WithAttributes(attribute.String("example name", "example value")),
			)
			defer span.End()
		}),
	)

	// trace assertion success
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"example span",
		attribute.String("example name", "example value"),
	)
}
```