# Metrics Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxmetrics-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxmetrics-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxmetrics)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxmetrics)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxmetrics)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxmetrics)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxmetrics)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxmetrics)](https://pkg.go.dev/github.com/ankorstore/yokai/fxmetrics)

## Overview

Yokai provides a [fxmetrics](https://github.com/ankorstore/yokai/tree/main/fxmetrics) module, allowing your application to provide [metrics](https://prometheus.io/docs/concepts/metric_types).

It wraps the [Prometheus](https://github.com/prometheus/client_golang) module.

## Installation

The [fxmetrics](https://github.com/ankorstore/yokai/tree/main/fxmetrics) module is automatically loaded by
the [fxcore](https://github.com/ankorstore/yokai/tree/main/fxcore).

When you use a Yokai [application template](https://ankorstore.github.io/yokai/applications/templates/), you have nothing to install, it's ready to use.

## Usage

This module will enable Yokai to collect registered metrics [collectors](https://github.com/prometheus/client_golang/blob/main/prometheus/collector.go), and make available to a metrics [registry](https://github.com/prometheus/client_golang/blob/main/prometheus/registry.go) in
its dependency injection system.

### Metrics creation

You can add metrics anywhere in your application.

For example:

```go title="internal/service/example.go"
package service

import (
	"fmt"
	
	"github.com/prometheus/client_golang/prometheus"
)

var ExampleCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "example_total",
	Help: "Example counter",
})

type ExampleService struct {}

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (s *ExampleService) DoSomething() {
	// service logic
	fmt.Println("do something")
	
	// increment counter
	ExampleCounter.Inc()
}
```

### Metrics registration

Even if convenient, it's recommended to NOT use the [promauto](https://github.com/prometheus/client_golang/tree/main/prometheus/promauto) way of registering metrics, as promauto uses a global registry that leads to data race conditions (especially while testing).

You can instead register your metrics collector with the `AsMetricsCollector()` function:

```go title="internal/services.go"
package internal

import (
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/foo/bar/service"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// register the ExampleCounter metrics collector
		fxmetrics.AsMetricsCollector(service.ExampleCounter),
	)
}
```

You can also register several metrics collectors at once with `AsMetricsCollectors()`.

### Metrics execution

The [fxcore](https://github.com/ankorstore/yokai/tree/main/fxcore) HTTP server will automatically:

- expose the configured metrics endpoints
- use the [registry](https://github.com/prometheus/client_golang/blob/main/prometheus/registry.go) to expose the registered metrics collectors

Following previous example, after invoking the `ExampleService`, the metrics endpoint will return:

```yaml title="[GET] /metrics"
# HELP example_total Example counter
# TYPE example_total counter
example_total 1
```