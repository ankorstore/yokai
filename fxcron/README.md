# Fx Cron Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxcron-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxcron-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxcron)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxcron)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxcron)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxcron)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxcron)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxcron)](https://pkg.go.dev/github.com/ankorstore/yokai/fxcron)

> [Fx](https://uber-go.github.io/fx/) module for [gocron](https://github.com/go-co-op/gocron).

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
  * [Dependencies](#dependencies)
  * [Loading](#loading)
  * [Configuration](#configuration)
  * [Cron jobs](#cron-jobs)
    * [Definition](#definition)
    * [Registration](#registration)
  * [Override](#override)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxcron
```

## Features

This module provides the possibility to run **internal** cron jobs in your application with:

- automatic panic recovery
- configurable cron jobs scheduling and execution options
- configurable logging, tracing and metrics for cron jobs executions

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
	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,                      // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxcron.FxCronModule,                          // load the module
		fx.Invoke(func(scheduler gocron.Scheduler) {
			scheduler.Start()                         // start the cron jobs scheduler
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
  log:
    level: info
    output: stdout
  trace:
    processor:
      type: stdout
  cron:
    scheduler:
      seconds: true                   # to allow seconds based cron jobs expressions (impact all jobs), disabled by default
      concurrency:
        limit:
          enabled: true               # to limit concurrent cron jobs executions, disabled by default
          max: 3                      # concurrency limit
          mode: wait                  # "wait" or "reschedule"
      stop:
        timeout: 5s                   # scheduler shutdown timeout for graceful cron jobs termination, 10 seconds by default
    jobs:                             # common cron jobs options
      execution:
        start:
          immediately: true           # to start cron jobs executions immediately (by default)
          at: "2023-01-01T14:00:00Z"  # or a given date time (RFC3339)
        limit:
          enabled: true               # to limit the number of per cron jobs executions, disabled by default
          max: 3                      # executions limit
      singleton:
        enabled: true                 # to execute the cron jobs in singleton mode, disabled by default
        mode: wait                    # "wait" or "reschedule"
    log:
      enabled: true                   # to log cron jobs executions, disabled by default (errors will always be logged).
      exclude:                        # to exclude by name cron jobs from logging
        - foo
        - bar
    metrics:
      collect:
        enabled: true                 # to collect cron jobs executions metrics (executions count and duration), disabled by default
        namespace: app                # cron jobs metrics namespace (default app.name value)
        subsystem: cron               # cron jobs metrics subsystem (default cron)
      buckets: 1, 1.5, 10, 15, 100    # to define custom cron jobs executions durations metric buckets (in seconds)
    trace:
      enabled: true                   # to trace cron jobs executions, disabled by default
      exclude:                        # to exclude by name cron jobs from tracing
        - foo
        - bar
```

Notes:

- the cron jobs executions logging will be based on the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog)
  module configuration
- the cron jobs executions tracing will be based on the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace)
  module configuration

Check the [configuration files documentation](https://github.com/ankorstore/yokai/tree/main/config#configuration-files)
for more details.

### Cron jobs

#### Definition

This module provides a simple [CronJob](registry.go) interface to implement for your cron jobs:

```go
package cron

import (
	"context"

	"github.com/ankorstore/yokai/fxcron"
	"path/to/service"
)

type SomeCron struct {
	service *service.SomeService
}

func NewSomeCron(service *service.SomeService) *SomeCron {
	return &SomeCron{
		service: service,
	}
}

func (c *SomeCron) Name() string {
	return "some cron job"
}

func (c *SomeCron) Run(ctx context.Context) error {
	// contextual job name and execution id
	name, id := fxcron.CtxCronJobName(ctx), fxcron.CtxCronJobExecutionId(ctx)

	// contextual tracing
	ctx, span := fxcron.CtxTracer(ctx).Start(ctx, "some span")
	defer span.End()

	// contextual logging
	fxcron.CtxLogger(ctx).Info().Msg("some log")

	// invoke autowired dependency
	err := c.service.DoSomething(ctx, name, id)

	// returned errors will automatically be logged
	return err
}
```

Notes:

- your cron job dependencies will be autowired
- you can access from the provided context:
	- the cron job name with `CtxCronJobName()`
	- the cron job execution id with `CtxCronJobExecutionId()`
	- the tracer with `CtxTracer()`, which will automatically add to your spans the `CronJob` name
	  and `CronJobExecutionID` attributes
	- the logger with `CtxLogger()`, which will automatically add to your log records the `cronJob` name
	  and `cronJobExecutionID` fields

#### Registration

Once ready, you can register and schedule your cron job with `AsCronJob()`:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/fx"
	"path/to/cron"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,      // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxcron.FxCronModule,          // load the module
		fx.Options(
			// register, autowire and schedule SomeCron to run every 2 minutes
			fxcron.AsCronJob(cron.NewSomeCron, `*/2 * * * *`),
		),
	).Run()
}
```

You can override, per job, the common job execution options by providing your own list
of [gocron.JobOption](https://github.com/go-co-op/gocron/blob/v2/job.go), for example:

```go
fxcron.AsCronJob(cron.NewSomeCron, `*/2 * * * *`, gocron.WithLimitedRuns(10)),
```

If you need cron jobs to be scheduled on the seconds level, configure the scheduler
with `modules.cron.scheduler.seconds=true`.

It will add `seconds` field to the beginning of the scheduling expression, for example to run every 5 seconds:

```go
fxcron.AsCronJob(cron.NewSomeCron, `*/5 * * * * *`),
```

Note: you can use [https://crontab.guru](https://crontab.guru) to help you with your scheduling definitions.

### Override

By default, the `gocron.Scheduler` is created by the [DefaultCronSchedulerFactory](factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/fx"
)

type CustomCronSchedulerFactory struct{}

func NewCustomCronSchedulerFactory() fxcron.CronSchedulerFactory {
	return &CustomCronSchedulerFactory{}
}

func (f *CustomCronSchedulerFactory) Create(options ...gocron.SchedulerOption) (gocron.Scheduler, error) {
	return gocron.NewScheduler(options...)
}

func main() {
	fx.New(
		fxconfig.FxConfigModule, // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxcron.FxCronModule,                         // load the module
		fx.Decorate(NewCustomCronSchedulerFactory),  // override the module with a custom factory
		fx.Invoke(func(scheduler gocron.Scheduler) { // invoke the cron scheduler
			// ...
		}),
	).Run()
}
```
