---
icon: material/cube-outline
---

# :material-cube-outline: Cron Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxcron-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxcron-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxcron)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxcron)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxcron)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxcron)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxcron)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxcron)](https://pkg.go.dev/github.com/ankorstore/yokai/fxcron)

## Overview

Yokai provides a [fxcron](https://github.com/ankorstore/yokai/tree/main/fxcron) module, providing a cron jobs scheduler to your application.

It wraps the [gocron](https://github.com/go-co-op/gocron) module.

It comes with:

- automatic panic recovery
- configurable cron jobs scheduling and execution options
- configurable logging, tracing and metrics for cron jobs executions

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxcron
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxcron"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxcron.FxCronModule,
	// ...
)
```

## Configuration

```yaml title="configs/config.yaml"
modules:
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
        namespace: foo                # cron jobs metrics namespace (empty by default)
        subsystem: bar                # cron jobs metrics subsystem (empty by default)
      buckets: 1, 1.5, 10, 15, 100    # to define custom cron jobs executions durations metric buckets (in seconds)
    trace:
      enabled: true                   # to trace cron jobs executions, disabled by default
      exclude:                        # to exclude by name cron jobs from tracing
        - foo
        - bar
```

## Usage

This module provides the possibility to register [CronJob](https://github.com/ankorstore/yokai/blob/main/fxcron/registry.go) implementations, with:

- a [cron expression](https://crontab.guru/)
- and an optional list of [JobOption](https://github.com/go-co-op/gocron/blob/v2/job.go).

They will be collected and given by Yokai to the [Scheduler](https://github.com/go-co-op/gocron/blob/v2/scheduler.go) in its dependency injection system.

### Cron jobs creation

You can create your cron jobs by implementing the [CronJob](https://github.com/ankorstore/yokai/blob/main/fxcron/registry.go) interface.

For example:

```go title="internal/cron/example.go"
package cron

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcron"
)

type ExampleCronJob struct {
	config *config.Config
}

func NewExampleCronJob(config *config.Config) *ExampleCronJob {
	return &ExampleCronJob{
		config: config,
	}
}

func (c *ExampleCronJob) Name() string {
	return "example-cron-job"
}

func (c *ExampleCronJob) Run(ctx context.Context) error {
	// contextual job name and execution id
	name, id := fxcron.CtxCronJobName(ctx), fxcron.CtxCronJobExecutionId(ctx)

	// contextual tracing
	ctx, span := fxcron.CtxTracer(ctx).Start(ctx, "example span")
	defer span.End()

	// contextual logging
	fxcron.CtxLogger(ctx).Info().Msg("example log from app:%s, job:%s, id:%s", c.config.AppName(), name, id)

	// returned errors will automatically be logged
	return nil
}
```

You can access from the provided context:

- the cron job name with `CtxCronJobName()`
- the cron job execution id with `CtxCronJobExecutionId()`
- the tracer with `CtxTracer()`, which will automatically add to your spans the `CronJob` name
  and `CronJobExecutionID` attributes
- the logger with `CtxLogger()`, which will automatically add to your log records the `cronJob` name
  and `cronJobExecutionID` fields

### Cron jobs registration

You can register your cron jobs with the `AsCronJob()` function in `internal/register.go`:

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcron"
	"github.com/foo/bar/cron"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		fxcron.AsCronJob(
			cron.NewExampleCronJob,        // register the ExampleCronJob
			`*/2 * * * *`,                 // to run every 2 minutes
			gocron.WithLimitedRuns(10),    // and with 10 max runs 
		),
		// ...
	)
}
```

This module also supports the definition of cron expression on `seconds` level with `modules.cron.scheduler.seconds=true`.

It will add `seconds` field to the beginning of the scheduling expression, for example, to run every 3 seconds:

```go
fxcron.AsCronJob(cron.NewExampleCronJob, `*/3 * * * * *`),
```

You can use [https://crontab.guru](https://crontab.guru/) for building you cron expressions.

### Cron jobs execution

Yokai will automatically start the [Scheduler](https://github.com/go-co-op/gocron/blob/v2/scheduler.go) with the registered cron jobs.

You can get, in real time, the status of your cron jobs on the [core](fxcore.md#dashboard) dashboard:

![](../../assets/images/cron-tutorial-core-jobs-light.png#only-light)
![](../../assets/images/cron-tutorial-core-jobs-dark.png#only-dark)

## Logging

To get logs correlation in your cron jobs, you need to retrieve the logger from the context with `log.CtxLogger()`:

```go
log.CtxLogger(ctx).Info().Msg("example message")
```

You can also use the shortcut function `fxcron.CtxLogger()`:

```go
fxcron.CtxLogger(ctx)
```

As a result, log records will have the `cronJob` name and `cronJobExecutionID` fields added automatically:

```
INF job execution success cronJob=example-cron-job cronJobExecutionID=507a78d2-b466-445c-a113-9a3a89f6fbc7 service=app system=cron
```

The cron jobs logging will be based on the [log](fxlog.md) module configuration.

## Tracing

To get traces correlation in your cron jobs, you need to retrieve the tracer provider from the context with `trace.CtxTracerProvider()`:

```go
ctx, span := trace.CtxTracerProvider(ctx).Tracer("example tracer").Start(ctx, "example span")
defer span.End()
```

You can also use the shortcut function `fxcron.CtxTracer()`:

```go
ctx, span := fxcron.CtxTracer(ctx).Start(ctx, "example span")
defer span.End()
```

As a result, in your application trace spans attributes:

```
service.name: app
CronJob: example-cron-job
CronJobExecutionID: 507a78d2-b466-445c-a113-9a3a89f6fbc7
...
```

The cron jobs tracing will be based on the [trace](fxtrace.md) module configuration.

## Metrics

You can enable cron jobs automatic metrics with `modules.cron.metrics.collect.enable=true`:

```yaml title="configs/config.yaml"
modules:
  cron:
    metrics:
      collect:
        enabled: true                 # to collect cron jobs executions metrics (executions count and duration), disabled by default
        namespace: foo                # cron jobs metrics namespace (empty by default)
        subsystem: bar                # cron jobs metrics subsystem (empty by default)
      buckets: 1, 1.5, 10, 15, 100    # to define custom cron jobs executions durations metric buckets (in seconds)
```

This will collect metrics about:

- cron job successes
- cron job failures
- cron job execution durations

For example, after starting Yokai's cron jobs scheduler, the [core](fxcore.md) HTTP server will expose in the configured metrics endpoint:

```makefile title="[GET] /metrics"
# ...
# HELP cron_execution_duration_seconds Duration of cron job executions in seconds
# TYPE cron_execution_duration_seconds histogram
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.001"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.002"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.005"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.01"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.02"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.05"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.1"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.2"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="0.5"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="1"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="2"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="5"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="10"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="20"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="50"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="100"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="200"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="500"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="1000"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="2000"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="5000"} 2
cron_execution_duration_seconds_bucket{job="example_cron_job",le="+Inf"} 2
cron_execution_duration_seconds_sum{job="example_cron_job"} 0.000227993
cron_execution_duration_seconds_count{job="example_cron_job"} 2
# HELP cron_execution_total Total number of cron job executions
# TYPE cron_execution_total counter
cron_execution_total{job="example_cron_job",status="success"} 2
```
