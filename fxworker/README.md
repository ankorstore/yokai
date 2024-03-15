# Fx Worker Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxworker-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxworker-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxworker)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxworker)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxworker)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxworker)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxworker)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxworker)](https://pkg.go.dev/github.com/ankorstore/yokai/fxworker)

> [Fx](https://uber-go.github.io/fx/) module for [worker](https://github.com/ankorstore/yokai/tree/main/worker).

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
  * [Dependencies](#dependencies)
  * [Loading](#loading)
  * [Configuration](#configuration)
  * [Registration](#registration)
  * [Override](#override)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxworker
```

## Features

This module provides a workers pool to your Fx application with:

- automatic panic recovery
- automatic logging
- automatic metrics
- possibility to defer workers
- possibility to limit workers max execution attempts

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
  "context"

  "github.com/ankorstore/yokai/fxconfig"
  "github.com/ankorstore/yokai/fxgenerate"
  "github.com/ankorstore/yokai/fxlog"
  "github.com/ankorstore/yokai/fxmetrics"
  "github.com/ankorstore/yokai/fxtrace"
  "github.com/ankorstore/yokai/fxworker"
  "github.com/ankorstore/yokai/worker"
  "go.uber.org/fx"
)

func main() {
  fx.New(
    fxconfig.FxConfigModule,                    // load the module dependencies
    fxlog.FxLogModule,
    fxtrace.FxTraceModule,
    fxtrace.FxTraceModule,
    fxmetrics.FxMetricsModule,
    fxgenerate.FxGenerateModule,
    fxworker.FxWorkerModule,                    // load the module
    fx.Invoke(func(pool *worker.WorkerPool) {
      pool.Start(context.Background())          // start the workers pool
    }),
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
  worker:
    defer: 0.1             # threshold in seconds to wait before starting all workers, immediate start by default
    attempts: 3            # max execution attempts in case of failures for all workers, no restart by default
    metrics:
      collect:
        enabled: true      # to collect metrics about workers executions
        namespace: foo     # workers metrics namespace (empty by default)
        subsystem: bar     # workers metrics subsystem (empty by default)
```

Notes:

- the workers logging will be based on the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog)
  module configuration
- the workers tracing will be based on the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace)
  module configuration

### Registration

This module provides the possibility to register
several [Worker](https://github.com/ankorstore/yokai/blob/main/worker/worker.go) implementations, with
optional [WorkerExecutionOption](https://github.com/ankorstore/yokai/blob/main/worker/option.go).

They will be then collected and given by Fx to
the [WorkerPool](https://github.com/ankorstore/yokai/blob/main/worker/pool.go), made available in the Fx container.

This is done via the `AsWorker()` function:

```go
package main

import (
  "context"

  "github.com/ankorstore/yokai/fxconfig"
  "github.com/ankorstore/yokai/fxgenerate"
  "github.com/ankorstore/yokai/fxlog"
  "github.com/ankorstore/yokai/fxmetrics"
  "github.com/ankorstore/yokai/fxtrace"
  "github.com/ankorstore/yokai/fxworker"
  "github.com/ankorstore/yokai/worker"
  "go.uber.org/fx"
)

type ExampleWorker struct{}

func NewExampleWorker() *ExampleWorker {
  return &ExampleWorker{}
}

func (w *ExampleWorker) Name() string {
  return "example-worker"
}

func (w *ExampleWorker) Run(ctx context.Context) error {
  worker.CtxLogger(ctx).Info().Msg("run")

  return nil
}

func main() {
  fx.New(
    fxconfig.FxConfigModule,                      // load the module dependencies
    fxlog.FxLogModule,
    fxtrace.FxTraceModule,
    fxtrace.FxTraceModule,
    fxmetrics.FxMetricsModule,
    fxgenerate.FxGenerateModule,
    fxworker.FxWorkerModule,                      // load the module
    fx.Provide(
      fxworker.AsWorker(
        NewExampleWorker,                         // register the ExampleWorker
        worker.WithDeferredStartThreshold(1),     // with a deferred start threshold of 1 second
        worker.WithMaxExecutionsAttempts(2),      // and 2 max execution attempts
      ),
    ),
    fx.Invoke(func(pool *worker.WorkerPool) {
      pool.Start(context.Background())            // start the workers pool
    }),
  ).Run()
}
```

To get more details about the features made available for your workers (contextual logging, tracing, etc.), check
the [worker module documentation](https://github.com/ankorstore/yokai/tree/main/worker).

### Override

By default, the `worker.WorkerPool` is created by
the [DefaultWorkerPoolFactory](https://github.com/ankorstore/yokai/blob/main/worker/factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
  "context"

  "github.com/ankorstore/yokai/fxconfig"
  "github.com/ankorstore/yokai/fxgenerate"
  "github.com/ankorstore/yokai/fxhealthcheck"
  "github.com/ankorstore/yokai/fxlog"
  "github.com/ankorstore/yokai/fxmetrics"
  "github.com/ankorstore/yokai/fxtrace"
  "github.com/ankorstore/yokai/fxworker"
  "github.com/ankorstore/yokai/healthcheck"
  "github.com/ankorstore/yokai/worker"
  "go.uber.org/fx"
)

type CustomWorkerPoolFactory struct{}

func NewCustomWorkerPoolFactory() worker.WorkerPoolFactory {
  return &CustomWorkerPoolFactory{}
}

func (f *CustomWorkerPoolFactory) Create(options ...worker.WorkerPoolOption) (*worker.WorkerPool, error) {
  return &worker.WorkerPool{...}, nil
}

func main() {
  fx.New(
    fxconfig.FxConfigModule,                     // load the module dependencies
    fxlog.FxLogModule,
    fxtrace.FxTraceModule,
    fxtrace.FxTraceModule,
    fxmetrics.FxMetricsModule,
    fxgenerate.FxGenerateModule,
    fxworker.FxWorkerModule,                     // load the module
    fx.Decorate(NewCustomWorkerPoolFactory),     // override the module with a custom factory
    fx.Invoke(func(pool *worker.WorkerPool) {
      pool.Start(context.Background())           // start the custom worker pool
    }),
  ).Run()
}
```
