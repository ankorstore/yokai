# Health Check Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/healthcheck-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/healthcheck-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/healthcheck)](https://goreportcard.com/report/github.com/ankorstore/yokai/healthcheck)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=healthcheck)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/healthcheck)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/healthcheck)](https://pkg.go.dev/github.com/ankorstore/yokai/healthcheck)

> Health check module compatible with [K8s probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).

<!-- TOC -->

* [Installation](#installation)
* [Documentation](#documentation)
	* [Probes](#probes)
	* [Checker](#checker)

<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/healthcheck
```

## Documentation

This module provides a [Checker](checker.go), that:

- can register any [CheckerProbe](probe.go) implementations and organise them for `startup`, `liveness` and /
  or `readiness` checks
- and execute them to get an overall [CheckerResult](checker.go)

The checker result will be considered as success if **ALL** registered probes checks are successful.

### Probes

This module provides a `CheckerProbe` interface to implement to provide your own probes, for example:

```go
package probes

import (
	"context"

	"github.com/ankorstore/yokai/healthcheck"
)

// success probe
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

// failure probe
type FailureProbe struct{}

func NewFailureProbe() *FailureProbe {
	return &FailureProbe{}
}

func (p *FailureProbe) Name() string {
	return "failureProbe"
}

func (p *FailureProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	return healthcheck.NewCheckerProbeResult(false, "some failure")
}
```

Notes:

- to perform more complex checks, you can inject dependencies to your probes implementation (ex: database, cache, etc)
- it is recommended to design your probes with a single responsibility (ex: one for database, one for cache, etc)

### Checker

You can create a [Checker](checker.go) instance, register your [CheckerProbe](probe.go) implementations, and launch
checks:

```go
package main

import (
	"context"
	"fmt"

	"path/to/probes"
	"github.com/ankorstore/yokai/healthcheck"
)

func main() {
	ctx := context.Background()

	checker, _ := healthcheck.NewDefaultCheckerFactory().Create(
		healthcheck.WithProbe(probes.NewSuccessProbe()),                       // registers for startup, readiness and liveness
		healthcheck.WithProbe(probes.NewFailureProbe(), healthcheck.Liveness), // registers for liveness only
	)

	// startup health check: invoke only successProbe
	startupResult := checker.Check(ctx, healthcheck.Startup)

	fmt.Printf("startup check success: %v", startupResult.Success) // startup check success: true

	for probeName, probeResult := range startupResult.ProbesResults {
		fmt.Printf("probe name: %s, probe success: %v, probe message: %s", probeName, probeResult.Success, probeResult.Message)
		// probe name: successProbe, probe success: true, probe message: some success
	}

	// liveness health check: invoke successProbe and failureProbe
	livenessResult := checker.Check(ctx, healthcheck.Liveness)

	fmt.Printf("liveness check success: %v", livenessResult.Success) // liveness check success: false

	for probeName, probeResult := range livenessResult.ProbesResults {
		fmt.Printf("probe name: %s, probe success: %v, probe message: %s", probeName, probeResult.Success, probeResult.Message)
		// probe name: successProbe, probe success: true, probe message: some success
		// probe name: failureProbe, probe success: false, probe message: some failure
	}

	// readiness health check: invoke successProbe and failureProbe
	readinessResult := checker.Check(ctx, healthcheck.Readiness)

	fmt.Printf("readiness check success: %v", readinessResult.Success) // readiness check success: false

	for probeName, probeResult := range readinessResult.ProbesResults {
		fmt.Printf("probe name: %s, probe success: %v, probe message: %s", probeName, probeResult.Success, probeResult.Message)
		// probe name: successProbe, probe success: true, probe message: some success
		// probe name: failureProbe, probe success: false, probe message: some failure
	}
}
```
