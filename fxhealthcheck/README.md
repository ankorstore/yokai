# Fx Health Check Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxhealthcheck-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxhealthcheck-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxhealthcheck)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxhealthcheck)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxhealthcheck)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxhealthcheck)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxhealthcheck)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxhealthcheck)](https://pkg.go.dev/github.com/ankorstore/yokai/fxhealthcheck)

> [Fx](https://uber-go.github.io/fx/) module for [healthcheck](https://github.com/ankorstore/yokai/tree/main/healthcheck).

<!-- TOC -->

* [Installation](#installation)
* [Documentation](#documentation)
	* [Loading](#loading)
	* [Registration](#registration)
	* [Override](#override)

<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxhealthcheck
```

## Documentation

### Loading

To load the module in your Fx application:

```go
package main

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/healthcheck"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxhealthcheck.FxHealthcheckModule,             // load the module
		fx.Invoke(func(checker *healthcheck.Checker) { // invoke the checker for liveness checks
			fmt.Printf("checker result: %v", checker.Check(context.Background(), healthcheck.Liveness))
		}),
	).Run()
}
```

### Registration

This module provides the possibility to register
several [CheckerProbe](https://github.com/ankorstore/yokai/blob/main/healthcheck/probe.go) implementations, and organise
them for `startup`, `liveness` and /
or `readiness` [checks](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).

They will be then collected and given by Fx to
the [Checker](https://github.com/ankorstore/yokai/blob/main/healthcheck/checker.go), made available in the Fx container.

This is done via the `AsCheckerProbe()` function:

```go
package main

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/healthcheck"
	"go.uber.org/fx"
)

// example success probe
type SuccessProbe struct{}

func NewSuccessProbe() *SuccessProbe {
	return &SuccessProbe{}
}

func (p *SuccessProbe) Name() string {
	return "successProbe"
}

func (p *SuccessProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	return healthcheck.NewCheckerProbeResult(true, "some success")
}

// example failure probe
type FailureProbe struct{}

func NewFailureProbe() *FailureProbe {
	return &FailureProbe{}
}

func (p *FailureProbe) Name() string {
	return "someProbe"
}

func (p *FailureProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	return healthcheck.NewCheckerProbeResult(false, "some failure")
}

// usage
func main() {
	fx.New(
		fxhealthcheck.FxHealthcheckModule, // load the module
		fx.Provide(
			fxhealthcheck.AsCheckerProbe(NewSuccessProbe),                       // register the SuccessProbe probe for startup, liveness and readiness checks
			fxhealthcheck.AsCheckerProbe(NewFailureProbe, healthcheck.Liveness), // register the FailureProbe probe for liveness checks only
		),
		fx.Invoke(func(checker *healthcheck.Checker) { // invoke the checker
			ctx := context.Background()

			fmt.Printf("startup: %v", checker.Check(ctx, healthcheck.Startup).Success)     // startup: true
			fmt.Printf("liveness: %v", checker.Check(ctx, healthcheck.Liveness).Success)   // liveness: false
			fmt.Printf("readiness: %v", checker.Check(ctx, healthcheck.Readiness).Success) // readiness: true
		}),
	).Run()
}
```

### Override

By default, the `healthcheck.Checker` is created by
the [DefaultCheckerFactory](https://github.com/ankorstore/yokai/blob/main/healthcheck/factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/healthcheck"
	"go.uber.org/fx"
)

type CustomCheckerFactory struct{}

func NewCustomCheckerFactory() healthcheck.CheckerFactory {
	return &CustomCheckerFactory{}
}

func (f *CustomCheckerFactory) Create(options ...healthcheck.CheckerOption) (*healthcheck.Checker, error) {
	return &healthcheck.Checker{...}, nil
}

func main() {
	fx.New(
		fxhealthcheck.FxHealthcheckModule,             // load the module
		fx.Decorate(NewCustomCheckerFactory),          // override the module with a custom factory
		fx.Invoke(func(checker *healthcheck.Checker) { // invoke the custom checker for readiness checks
			checker.Check(context.Background(), healthcheck.Readiness)
		}),
	).Run()
}
```
