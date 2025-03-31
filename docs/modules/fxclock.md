---
title: Modules - Clock
icon: material/cube-outline
---

# :material-cube-outline: Clock Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxclock-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxclock-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxclock)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxclock)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxclock)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxclock)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxclock)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxclock)](https://pkg.go.dev/github.com/ankorstore/yokai/fxclock)

## Overview

Yokai provides a [fxclock](https://github.com/ankorstore/yokai/tree/main/fxclock) module, that you can use to control time.

It wraps the [clockwork](https://github.com/jonboulle/clockwork) module.

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxclock
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxclock"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
    fxclock.FxClockModule,
	// ...
)
```

## Usage

This module provides a [clockwork.Clock](https://github.com/jonboulle/clockwork) instance, ready to inject in your code.

This is particularly useful if you need to control time (set time, fast-forward, ...).

For example:

```go title="internal/service/example.go"
package service

import (
	"github.com/jonboulle/clockwork"
)

type ExampleService struct {
	clock clockwork.Clock
}

func NewExampleService(clock clockwork.Clock) *ExampleService {
	return &ExampleService{
		clock: clock,
	}
}

func (s *ExampleService) Now() string {
	return s.clock.Now().String()
}
```

See the underlying vendor [documentation](https://github.com/jonboulle/clockwork) for more details.

## Testing

This module provides a [*clockwork.FakeClock](https://github.com/jonboulle/clockwork) instance, that will be automatically injected as `clockwork.Clock` in your constructors in `test` mode.

### Global time

By default, the fake clock is set to `time.Now()` (your test execution time).

You can configure the global time in your test in your testing configuration file (for all your tests), in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339) format:

```yaml title="configs/config_test.yaml"
modules:
  clock:
    test:
      time: "2006-01-02T15:04:05Z07:00" # time in RFC3339 format
```

You can also [override this value](https://ankorstore.github.io/yokai/modules/fxconfig/#env-var-substitution), per test, by setting the `MODULES_CLOCK_TEST_TIME` env var.

### Time control

You can `populate` the [*clockwork.FakeClock](https://github.com/jonboulle/clockwork) from your test to control time:

```go title="internal/service/example_test.go"
package service_test

import (
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxsql"
	"github.com/foo/bar/internal"
	"github.com/foo/bar/internal/service"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestExampleService(t *testing.T) {
	testTime := "2025-03-30T12:00:00Z"
	expectedTime, err := time.Parse(time.RFC3339, testTime)
	assert.NoError(t, err)

	t.Setenv("MODULES_CLOCK_TEST_TIME", testTime)

	var svc service.ExampleService
	var clock *clockwork.FakeClock

	internal.RunTest(t, fx.Populate(&svc, &clock))

	// current time as configured above
	assert.Equal(t, expectedTime, svc.Now()) // 2025-03-30T12:00:00Z

	clock.Advance(5 * time.Hour)

	// current time is now advanced by 5 hours
	assert.Equal(t, expectedTime.Add(5*time.Hour), svc.Now()) // 2025-03-30T17:00:00Z
}
```

See [tests example](https://github.com/ankorstore/yokai/blob/main/fxclock/module_test.go) for more details.
