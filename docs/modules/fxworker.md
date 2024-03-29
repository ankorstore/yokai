---
icon: material/cube-outline
---

# :material-cube-outline: Worker Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxworker-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxworker-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxworker)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxworker)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxworker)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxworker)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxworker)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxworker)](https://pkg.go.dev/github.com/ankorstore/yokai/fxworker)

## Overview

Yokai provides a [fxworker](https://github.com/ankorstore/yokai/tree/main/fxworker) module, providing a workers pool to your application.

It wraps the [worker](https://github.com/ankorstore/yokai/tree/main/worker) module, based on [sync](https://pkg.go.dev/sync).

It comes with:

- automatic panic recovery
- automatic logging
- automatic metrics
- possibility to defer workers
- possibility to limit workers max execution attempts

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxworker
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxworker"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxworker module
	fxworker.FxWorkerModule,
	// ...
)
```

## Configuration

```yaml title="configs/config.yaml"
modules:
  worker:
    defer: 0.1             # threshold in seconds to wait before starting all workers, immediate start by default
    attempts: 3            # max execution attempts in case of failures for all workers, no restart by default
    metrics:
      collect:
        enabled: true      # to collect metrics about workers executions
        namespace: foo     # workers metrics namespace (empty by default)
        subsystem: bar     # workers metrics subsystem (empty by default)
```

## Usage

This module provides the possibility to register several [Worker](https://github.com/ankorstore/yokai/blob/main/worker/worker.go) implementations, with an
optional list of [WorkerExecutionOption](https://github.com/ankorstore/yokai/blob/main/worker/option.go).

They will be collected and given by Yokai to the [WorkerPool](https://github.com/ankorstore/yokai/blob/main/worker/pool.go) in its dependency injection system.

### Workers creation

You can create your workers by implementing the [Worker](https://github.com/ankorstore/yokai/blob/main/worker/worker.go) interface.

For example:

```go title="internal/worker/example.go"
package worker

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/worker"
)

type ExampleWorker struct {
	config *config.Config
}

func NewExampleWorker(config *config.Config) *ExampleWorker {
	return &ExampleWorker{
		config: config,
	}
}

func (w *ExampleWorker) Name() string {
	return "example-worker"
}

func (w *ExampleWorker) Run(ctx context.Context) error {
	logger := worker.CtxLogger(ctx)

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("stopping")

			return nil
		default:
			logger.Info().Msg("running")

			// The sleep interval can be configured in the application config files.
			time.Sleep(time.Duration(w.config.GetFloat64("config.example-worker.interval")) * time.Second)
		}
	}
}
```

And the corresponding example configuration:

```yaml title="configs/config.yaml"
config:
  example-worker:
	interval: 3
```

### Workers registration

You can register your workers with the `AsWorker()` function:

```go title="internal/services.go"
package internal

import (
	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/worker"
	w "github.com/foo/bar/worker"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		fxworker.AsWorker(
			w.NewExampleWorker,                   // register the ExampleWorker
			worker.WithDeferredStartThreshold(1), // with a deferred start of 1 second
			worker.WithMaxExecutionsAttempts(2),  // and 2 max execution attempts 
		),
		// ...
	)
}
```

### Workers execution

Yokai will automatically start the [WorkerPool](https://github.com/ankorstore/yokai/blob/main/worker/pool.go) containing the registered workers.

You can get, in real time, the status of your workers executions on the [core](fxcore.md#dashboard) dashboard:

![](../../assets/images/dash-workers-light.png#only-light)
![](../../assets/images/dash-workers-dark.png#only-dark)

## Logging

To get logs correlation in your workers, you need to retrieve the logger from the context with `log.CtxLogger()`:

```go
log.CtxLogger(ctx).Info().Msg("example message")
```

You can also use the shortcut function `worker.CtxLogger()`:

```go
worker.CtxLogger(ctx)
```

As a result, log records will have the `worker` name and `workerExecutionID` fields added automatically:

```
INF example message module=worker service=app worker=example-worker workerExecutionID=b57be88f-163f-4a81-bf24-a389c93d804b
```

The workers logging will be based on the [log](fxlog.md) module configuration.

## Tracing

To get traces correlation in your workers, you need to retrieve the tracer provider from the context with `trace.CtxTracerProvider()`:

```go
ctx, span := trace.CtxTracerProvider(ctx).Tracer("example tracer").Start(ctx, "example span")
defer span.End()
```

You can also use the shortcut function `worker.CtxTracer()`:

```go
ctx, span := worker.CtxTracer(ctx).Start(ctx, "example span")
defer span.End()
```

As a result, in your application trace spans attributes:

```
service.name: app
Worker: example-worker
WorkerExecutionID: b57be88f-163f-4a81-bf24-a389c93d804b
...
```

The workers tracing will be based on the [trace](fxtrace.md) module configuration.

## Metrics

You can enable workers executions automatic metrics with `modules.worker.metrics.collect.enable=true`:

```yaml title="configs/config.yaml"
modules:
  worker:
    metrics:
      collect:
        enabled: true      # to collect metrics about workers executions
        namespace: foo     # workers metrics namespace (empty by default)
        subsystem: bar     # workers metrics subsystem (empty by default)
```

This will collect metrics about:

- workers `start` and `restart`
- workers `successes`
- workers `failures`

For example, after starting Yokai's workers pool, the [core](fxcore.md) HTTP server will expose in the configured metrics endpoint:

```makefile title="[GET] /metrics"
# ...
# HELP worker_executions_total Total number of workers executions
# TYPE worker_executions_total counter
worker_executions_total{status="started",worker="example-worker"} 1
```
