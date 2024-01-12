# Fx Metrics Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxmetrics-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxmetrics-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxmetrics)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxmetrics)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxmetrics)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxmetrics)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxmetrics)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxmetrics)](https://pkg.go.dev/github.com/ankorstore/yokai/fxmetrics)

> [Fx](https://uber-go.github.io/fx/) module for [prometheus](https://github.com/prometheus/client_golang).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
	* [Dependencies](#dependencies)
	* [Loading](#loading)
	* [Registration](#registration)
	* [Override](#override)
	* [Testing](#testing)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxmetrics
```

## Documentation

### Dependencies

This module is intended to be used alongside:

- the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module
- the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module

### Loading

To load the module in your Fx application:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,                        // load the module dependencies
		fxlog.FxLogModule,
		fxmetrics.FxMetricsModule,                      // load the module
		fx.Invoke(func(registry *prometheus.Registry) { // invoke the metrics registry
			// ...
		}),
	).Run()
}
```

### Registration

This module provides the possibility to register your metrics [collectors](https://github.com/prometheus/client_golang/blob/main/prometheus/collector.go) in a common `*prometheus.Registry` via `AsMetricsCollector()`:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

var SomeCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "some_total",
	Help: "some help",
})

func main() {
	fx.New(
		fxconfig.FxConfigModule,                       // load the module dependencies
		fxlog.FxLogModule,
		fxmetrics.FxMetricsModule,                     // load the module
		fx.Options(
			fxmetrics.AsMetricsCollector(SomeCounter), // register the counter
		),
		fx.Invoke(func() {
			SomeCounter.Inc()                          // manipulate the counter
		}),
	).Run()
}
```

**Important**: even if convenient, it's recommended to **NOT** use the [promauto](https://github.com/prometheus/client_golang/tree/main/prometheus/promauto) way of registering metrics,
but to use instead `fxmetrics.AsMetricsCollector()`, as `promauto` uses a global registry that leads to data race
conditions in testing.

Also, if you want to register several collectors at once, you can use `fxmetrics.AsMetricsCollectors()`

### Override

By default, the `*prometheus.Registry` is created by the [DefaultMetricsRegistryFactory](factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

type CustomMetricsRegistryFactory struct{}

func NewCustomMetricsRegistryFactory() fxmetrics.MetricsRegistryFactory {
	return &CustomMetricsRegistryFactory{}
}

func (f *CustomMetricsRegistryFactory) Create() (*prometheus.Registry, error) {
	return prometheus.NewPedanticRegistry(), nil
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,                        // load the module dependencies
		fxlog.FxLogModule,
		fxmetrics.FxMetricsModule,                      // load the module
		fx.Decorate(NewCustomMetricsRegistryFactory),   // override the module with a custom factory
		fx.Invoke(func(registry *prometheus.Registry) { // invoke the custom registry
			// ...
		}),
	).Run()
}
```

### Testing

This module provides the possibility to easily test your metrics with the prometheus package [testutil](https://github.com/prometheus/client_golang/tree/main/prometheus/testutil) helpers.

```go
package main_test

import (
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var SomeCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "some_total",
	Help: "some help",
})

func TestSomeCounter(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var registry *prometheus.Registry

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxmetrics.FxMetricsModule,
		fx.Options(
			fxmetrics.AsMetricsCollector(SomeCounter),
		),
		fx.Invoke(func() {
			SomeCounter.Add(9)
		}),
		fx.Populate(&registry),
	).RequireStart().RequireStop()

	// metric assertions
	expectedHelp := `
		# HELP some_total some help
		# TYPE some_total counter
	`
	expectedMetric := `
		some_total 9
	`

	err := testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedHelp+expectedMetric),
		"some_total",
	)
	assert.NoError(t, err)
}
```
