---
icon: material/cube-outline
---

# :material-cube-outline: Health Check Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxhealthcheck-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxhealthcheck-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxhealthcheck)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxhealthcheck)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxhealthcheck)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxhealthcheck)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxhealthcheck)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxhealthcheck)](https://pkg.go.dev/github.com/ankorstore/yokai/fxhealthcheck)

## Overview

Yokai provides a [fxhealthcheck](https://github.com/ankorstore/yokai/tree/main/fxhealthcheck) module, allowing your application to provide [K8s probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).

It wraps the [healthcheck](https://github.com/ankorstore/yokai/tree/main/healthcheck) module.

## Installation

The [fxhealthcheck](https://github.com/ankorstore/yokai/tree/main/fxhealthcheck) module is automatically loaded by
the [fxcore](https://github.com/ankorstore/yokai/tree/main/fxcore).

When you use a Yokai [application template](https://ankorstore.github.io/yokai/applications/templates/), you have nothing to install, it's ready to use.

## Usage

This module will enable Yokai to collect registered [CheckerProbe](https://github.com/ankorstore/yokai/blob/main/healthcheck/probe.go) implementations, and make them available to the [Checker](https://github.com/ankorstore/yokai/blob/main/healthcheck/checker.go) in
its dependency injection system.

You can register probes for `startup`, `liveness` and / or `readiness` [checks](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).

The check result will be considered as success if ALL registered probes checks are successful.

Notes:

- to perform complex checks, you can inject dependencies to your probes implementation (ex: database, cache, etc)
- it is recommended to design your probes with a single responsibility (ex: one for database, one for cache, etc)


### Probes creation

You can create your probes by implementing the [CheckerProbe](https://github.com/ankorstore/yokai/blob/main/healthcheck/probe.go) interface.

For example:

```go title="internal/probe/success.go"
package probe

import (
	"context"
	
	"github.com/ankorstore/yokai/healthcheck"
)

type SuccessProbe struct{}

func NewSuccessProbe() *SuccessProbe {
	return &SuccessProbe{}
}

func (p *SuccessProbe) Name() string {
	return "successProbe"
}

func (p *SuccessProbe) Check(context.Context) *healthcheck.CheckerProbeResult {
	return healthcheck.NewCheckerProbeResult(true, "success example message")
}
```

and

```go title="internal/probe/failure.go"
package probe

import (
	"context"
	
	"github.com/ankorstore/yokai/healthcheck"
)

type FailureProbe struct{}

func NewFailureProbe() *FailureProbe {
	return &FailureProbe{}
}

func (p *FailureProbe) Name() string {
	return "failureProbe"
}

func (p *FailureProbe) Check(context.Context) *healthcheck.CheckerProbeResult {
	return healthcheck.NewCheckerProbeResult(false, "failure example message")
}
```

### Probes registration

You can register your probes for `startup`, `liveness` and / or `readiness` checks with the `AsCheckerProbe()` function:

```go title="internal/services.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/foo/bar/probe"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// register the SuccessProbe probe for startup, liveness and readiness checks
		fxhealthcheck.AsCheckerProbe(probe.NewSuccessProbe),
		// register the FailureProbe probe for liveness checks only
		fxhealthcheck.AsCheckerProbe(probe.NewFailureProbe, healthcheck.Liveness), 
		// ...
	)
}
```

### Probes execution

The [fxcore](https://github.com/ankorstore/yokai/tree/main/fxcore) HTTP server will automatically:

- expose the configured health check endpoints
- use the [Checker](https://github.com/ankorstore/yokai/blob/main/healthcheck/checker.go) to run the registered probes

Following previous example:

- calling the `startup` endpoint will return a `200` response:

```json title="[GET] /healthz"
{
	"success": true, 
    "probes": {
		"successProbe": {
			"success": true,
			"message": "success example message"
		}
	}
}
```

- calling the `liveness` endpoint will return a `500` response:

```json title="[GET] /livez"
{
	"success": false, 
    "probes": {
		"successProbe": {
			"success": true,
			"message": "success example message"
		},
		"failureProbe": {
			"success": false,
			"message": "failure example message"
		}
	}
}
```

- calling the `readiness` endpoint will return a `200` response:

```json title="[GET] /readyz"
{
	"success": true, 
    "probes": {
		"successProbe": {
			"success": true,
			"message": "success example message"
		}
	}
}
```