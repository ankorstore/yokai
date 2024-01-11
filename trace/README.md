# Trace Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/trace-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/trace-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/trace)](https://goreportcard.com/report/github.com/ankorstore/yokai/trace)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=trace)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/trace)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ftrace)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/trace)](https://pkg.go.dev/github.com/ankorstore/yokai/trace)

> Tracing module based on [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go).

<!-- TOC -->

* [Installation](#installation)
* [Documentation](#documentation)
	* [Configuration](#configuration)
	* [Usage](#usage)
		* [Context](#context)
		* [Span processors](#span-processors)
			* [Noop span processor](#noop-span-processor)
			* [Stdout span processor](#stdout-span-processor)
			* [OTLP gRPC span processor](#otlp-grpc-span-processor)
			* [Test span processor](#test-span-processor)
		* [Samplers](#samplers)
			* [Parent based always on](#parent-based-always-on)
			* [Parent based always off](#parent-based-always-off)
			* [Parent based trace id ratio](#parent-based-trace-id-ratio)
			* [Always on](#always-on)
			* [Always off](#always-off)
			* [Trace id ratio](#trace-id-ratio)

<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/trace
```

## Documentation

### Configuration

Since The `TracerProviderFactory` use
the [resource.Default()](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/resource#Default) by default, you can use:

- `OTEL_SERVICE_NAME` env variable to configure your tracing service name
- `OTEL_RESOURCE_ATTRIBUTES` env variable to configure your resource attributes

### Usage

This module provides a [TracerProviderFactory](factory.go), allowing to set up a `TracerProvider` easily.

```go
package main

import (
	"github.com/ankorstore/yokai/trace"
	"go.opentelemetry.io/otel/sdk/resource"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create()

	// equivalent to
	tp, _ = trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),                                       // set the tracer provider as global
		trace.WithResource(resource.Default()),                   // use the default resource
		trace.WithSampler(trace.NewParentBasedAlwaysOnSampler()), // use parent based always on sampling
		trace.WithSpanProcessor(trace.NewNoopSpanProcessor()),    // use noop processor (void trace spans)
	)
}
```

See available [factory options](option.go).

#### Context

This module provides the `CtxTracerProvider()` function that allow to extract the tracer provider from
a `context.Context`.

If no tracer provider is found in context,
the [global tracer](https://github.com/open-telemetry/opentelemetry-go/blob/main/trace.go) will be used.

#### Span processors

This modules comes with 4 `SpanProcessor` ready to use:

- `Noop`: to async void traces (default)
- `Stdout`: to async print traces to the standard output
- `OtlpGrpc`: to async send traces to [OTLP/gRPC](https://opentelemetry.io/docs/specs/otlp/#otlpgrpc) collectors (
  ex: [Jaeger](https://www.jaegertracing.io/), [Grafana](https://grafana.com/docs/tempo/latest/configuration/grafana-agent/#grafana-agent),
  etc.)
- `Test`: to sync store traces in memory (for testing assertions)

##### Noop span processor

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/trace"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewNoopSpanProcessor()),
	)

	// voids trace span 
	_, span := tp.Tracer("default").Start(context.Background(), "my span")
	defer span.End()
}
```

##### Stdout span processor

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/trace"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewStdoutSpanProcessor(stdouttrace.WithPrettyPrint())),
	)

	// pretty prints trace span to stdout
	_, span := tp.Tracer("default").Start(context.Background(), "my span")
	defer span.End()
}
```

##### OTLP gRPC span processor

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/trace"
)

func main() {
	ctx := context.Background()

	conn, _ := trace.NewOtlpGrpcClientConnection(ctx, "jaeger:4317")
	proc, _ := trace.NewOtlpGrpcSpanProcessor(ctx, conn)

	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(proc),
	)

	// sends trace span to jaeger:4317
	_, span := tp.Tracer("default").Start(ctx, "my span")
	defer span.End()
}
```

##### Test span processor

```go
package main

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
)

func main() {
	ex := tracetest.NewDefaultTestTraceExporter()

	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(ex)),
	)

	// sends trace span to test exporter
	_, span := tp.Tracer("default").Start(context.Background(), "my span")
	defer span.End()

	// check
	fmt.Printf("has span: %v", ex.HasSpan("my span")) // has span: true
}
```

You can use the provided [test assertion helpers](tracetest/assert.go) in your tests:

- `AssertHasTraceSpan`: to assert on exact name and exact attributes match
- `AssertHasNotTraceSpan`: to assert on exact name and exact attributes non match
- `AssertContainTraceSpan`: to assert on exact name and partial attributes match
- `AssertContainNotTraceSpan`: to assert on exact name and partial attributes non match

```go
package main_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"go.opentelemetry.io/otel/attribute"
)

func TestTracer(t *testing.T) {
	ex := tracetest.NewDefaultTestTraceExporter()

	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(ex)),
	)

	// sends trace span to test exporter
	_, span := tp.Tracer("default").Start(
		context.Background(),
		"my span",
		attribute.String("string attr name", "string attr value"),
		attribute.Int("int attr name", 42),
	)
	span.End()

	// assertion success
	tracetest.AssertHasTraceSpan(
		t,
		ex,
		"my span",
		attribute.String("string attr name", "string attr value"),
		attribute.Int("int attr name", 42),
	)

	// assertion success
	tracetest.AssertHasNotTraceSpan(
		t,
		ex,
		"my span",
		attribute.String("string attr name", "string attr value"),
		attribute.Int("int attr name", 24),
	)

	// assertion success
	tracetest.AssertContainTraceSpan(
		t,
		ex,
		"my span",
		attribute.String("string attr name", "attr value"),
		attribute.Int("int attr name", 42),
	)

	// assertion success
	tracetest.AssertContainNotTraceSpan(
		t,
		ex,
		"my span",
		attribute.String("string attr name", "attr value"),
		attribute.Int("int attr name", 24),
	)
}
```

#### Samplers

This modules comes with 6 `Samplers` ready to use:

- `ParentBasedAlwaysOn`: always on depending on parent (default)
- `ParentBasedAlwaysOff`: always off depending on parent
- `ParentBasedTraceIdRatio`: trace id ratio based depending on parent
- `AlwaysOn`: always on
- `AlwaysOff`: always off
- `TraceIdRatio`: trace id ratio based

Note: parent based samplers returns a composite sampler which behaves differently, based on the parent of the span:

- if the span has no parent, the embedded sampler is used to make sampling decision
- if the span has a parent, it depends on whether the parent is remote and whether it is sampled

##### Parent based always on

```go
package main

import (
	"github.com/ankorstore/yokai/trace"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSampler(trace.NewParentBasedAlwaysOnSampler()),
	)
}
```

##### Parent based always off

```go
package main

import (
	"github.com/ankorstore/yokai/trace"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSampler(trace.NewParentBasedAlwaysOffSampler()),
	)
}
```

##### Parent based trace id ratio

```go
package main

import (
	"github.com/ankorstore/yokai/trace"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSampler(trace.NewParentBasedTraceIdRatioSampler(0.5)),
	)
}
```

##### Always on

```go
package main

import (
	"github.com/ankorstore/yokai/trace"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSampler(trace.NewAlwaysOnSampler()),
	)
}
```

##### Always off

```go
package main

import (
	"github.com/ankorstore/yokai/trace"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSampler(trace.NewAlwaysOffSampler()),
	)
}
```

##### Trace id ratio

```go
package main

import (
	"github.com/ankorstore/yokai/trace"
)

func main() {
	tp, _ := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSampler(trace.NewTraceIdRatioSampler(0.5)),
	)
}
```
