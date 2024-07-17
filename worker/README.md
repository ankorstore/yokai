# Worker Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/worker-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/worker-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/worker)](https://goreportcard.com/report/github.com/ankorstore/yokai/worker)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=worker)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/worker)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Fworker)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/worker)](https://pkg.go.dev/github.com/ankorstore/yokai/worker)

> Worker module based on [sync](https://pkg.go.dev/sync).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
  * [Workers](#workers)
  * [WorkerPool](#workerpool)
  * [Logging](#logging)
  * [Tracing](#tracing)
  * [Metrics](#metrics)
  * [Healthcheck](#healthcheck)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/worker
```

## Documentation

This module provides a [WorkerPool](pool.go), that:

- can register any [Worker](worker.go) implementation,
- execute them in a [sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup),
- and give at any time [WorkerExecution](execution.go) reports to check the workers status and events.

The [WorkerPool](pool.go) can be configured to:

- defer all workers start with a threshold in seconds: `0` by default (start immediately)
- attempt a maximum amount of runs in case of failures: `1` by default (no restarts)

The [Worker](worker.go) executions:

- have a unique identifier
- have automatic panic recovery
- are automatically logged
- are automatically generating metrics

### Workers

This module provides a `Worker` interface to implement to provide your own workers, for example:

```go
package workers

import (
	"context"

	"github.com/ankorstore/yokai/worker"
)

// classic worker
type ClassicWorker struct{}

func NewClassicWorker() *ClassicWorker {
	return &ClassicWorker{}
}

func (w *ClassicWorker) Name() string {
	return "classic-worker"
}

func (w *ClassicWorker) Run(ctx context.Context) error {
	worker.CtxLogger(ctx).Info().Msg("run")

	return nil
}

// cancellable worker
type CancellableWorker struct{}

func NewCancellableWorker() *CancellableWorker {
	return &CancellableWorker{}
}

func (w *CancellableWorker) Name() string {
	return "cancellable-worker"
}

func (w *CancellableWorker) Run(ctx context.Context) error {
	logger := worker.CtxLogger(ctx)

	for {
		select {
		// when the WorkerPool stops, the ctx cancellation is forwarded to the workers
		case <-ctx.Done():
			logger.Info().Msg("cancel")

			return w.cancel()
		default:
			logger.Info().Msg("run")

			return w.run(ctx)
		}
	}
}

func (w *CancellableWorker) run(ctx context.Context) error {
	// your worker logic
}

func (w *CancellableWorker) cancel() error {
	// your worker cancel logic, for example graceful shutdown
}
```

Notes:

- to perform more complex tasks, you can inject dependencies to your workers implementation (ex: database, cache, etc.)
- it is recommended to design your workers with a single responsibility

### WorkerPool

You can create a [WorkerPool](pool.go) instance with the [DefaultWorkerPoolFactory](factory.go), register
your [Worker](worker.go) implementations, and start them:

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/worker"
	"path/to/workers"
)

func main() {
	// create the pool
	pool, _ := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithGlobalDeferredStartThreshold(1),                                             // will defer all workers start of 1 second
		worker.WithGlobalMaxExecutionsAttempts(2),                                              // will run 2 times max failing workers
		worker.WithWorker(workers.NewClassicWorker(), worker.WithDeferredStartThreshold(3)),    // registers the ClassicWorker, with a deferred start of 3 second
		worker.WithWorker(workers.NewCancellableWorker(), worker.WithMaxExecutionsAttempts(4)), // registers the CancellableWorker, with 4 runs max
	)

	// start the pool
	pool.Start(context.Background())

	// get all workers execution reports, in real time
	executions := pool.Executions()

	// stop the pool (will forward context cancellation to each worker)
	pool.Stop()

	// get a specific worker execution report, after pool stop
	execution, _ := pool.Execution("cancellable-worker")
}
```

### Logging

You can use the [CtxLogger()](context.go) function to retrieve the
contextual [log.Logger](https://github.com/ankorstore/yokai/tree/main/log) from your workers, and emit correlated logs.

The workers executions are logged, with the following fields added automatically to each log records:

- `worker`: worker name
- `workerExecutionID`: worker execution id

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/worker"
)

type LoggingWorker struct{}

func NewLoggingWorker() *LoggingWorker {
	return &LoggingWorker{}
}

func (w *LoggingWorker) Name() string {
	return "logging-worker"
}

func (w *LoggingWorker) Run(ctx context.Context) error {
	// log the current worker name and execution id
	worker.CtxLogger(ctx).Info().Msgf(
		"execution %s for worker %s",
		worker.CtxWorkerName(ctx),        // contextual worker name
		worker.CtxWorkerExecutionId(ctx), // contextual worker execution id
	)

	return nil
}

func main() {
	// create the pool
	pool, _ := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(NewLoggingWorker()), // registers the LoggingWorker
	)

	// start the pool
	pool.Start(context.Background())
}
```

### Tracing

You can use the [CtxTracer()](context.go) function to retrieve the contextual tracer from your workers, and emit
correlated spans: they will have the `Worker` and `WorkerExecutionID` attributes added with respectively the worker name
and execution id.

This module provides the [AnnotateTracerProvider](trace.go) function, to extend
a [TracerProvider](https://github.com/open-telemetry/opentelemetry-go/blob/main/sdk/trace/provider.go) to add
automatically current worker information id to the spans emitted during a worker execution:

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/worker"
	"go.opentelemetry.io/otel/trace"
)

// tracing worker
type TracingWorker struct{}

func NewTracingWorker() *TracingWorker {
	return &TracingWorker{}
}

func (w *TracingWorker) Name() string {
	return "tracing-worker"
}

func (w *TracingWorker) Run(ctx context.Context) error {
	// emit some trace span
	_, span := worker.CtxTracer(ctx).Start(ctx, "some span")
	defer span.End()

	return nil
}

func main() {
	// tracer provider
	tp := trace.GetTracerProvider()

	// annotate the tracer provider
	worker.AnnotateTracerProvider(tp)

	// create the pool
	pool, _ := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(NewTracingWorker()),
	)

	// start the pool
	pool.Start(context.Background())
}
```

### Metrics

The [WorkerPool](pool.go) automatically generate metrics about:

- started workers
- re started workers
- workers stopped with success
- workers stopped with error

To enable those metrics in a [registry](https://github.com/prometheus/client_golang/blob/main/prometheus/registry.go),
simply call `Register` on the [WorkerMetrics](metrics.go) of the [WorkerPool](pool.go):

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/worker"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	// metrics registry
	registry := prometheus.NewRegistry()

	// create the pool
	pool, _ := worker.NewDefaultWorkerPoolFactory().Create()

	// register the pool metrics
	pool.Metrics().Register(registry)

	// start the pool
	pool.Start(context.Background())
}
```

### Healthcheck

This module provides an [WorkerProbe](healthcheck/probe.go), compatible with
the [healthcheck module](https://github.com/ankorstore/yokai/tree/main/healthcheck):

```go
package main

import (
	"context"

	yokaihc "github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/healthcheck"
)

func main() {
	// create the pool
	pool, _ := worker.NewDefaultWorkerPoolFactory().Create()

	// create the checker with the worker probe
	checker, _ := yokaihc.NewDefaultCheckerFactory().Create(
		yokaihc.WithProbe(healthcheck.NewWorkerProbe(pool)),
	)

	// start the pool
	pool.Start(context.Background())

	// run the checker
	res, _ := checker.Check(context.Background(), yokaihc.Readiness)
}
```

This probe is successful if all the executions statuses of the [WorkerPool](pool.go) are healthy.
