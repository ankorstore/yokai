# Fx Clock Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxclock-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxclock-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxclock)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxclock)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxclock)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxclock)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxclock)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxclock)](https://pkg.go.dev/github.com/ankorstore/yokai/fxclock)

> [Fx](https://uber-go.github.io/fx/) module for [clockwork](https://github.com/jonboulle/clockwork).

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
  * [Dependencies](#dependencies)
  * [Loading](#loading)
  * [Usage](#usage)
  * [Testing](#testing)
    * [Global time](#global-time)
    * [Time control](#time-control)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxclock
```

## Features

This module provides a [clockwork.Clock](https://github.com/jonboulle/clockwork) instance for your application, that you
can use to control time.

## Documentation

### Dependencies

This module is intended to be used alongside the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig)
module.

### Loading

To load the module in your Fx application:

```go
package main

import (
	"time"

	"github.com/ankorstore/yokai/fxclock"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule, // load the module dependencies
		fxclock.FxClockModule,   // load the module
		fx.Invoke(func(clock clockwork.Clock) { // invoke the clock
			clock.Sleep(3 * time.Second)
		}),
	).Run()
}
```

### Usage

This module provides a [clockwork.Clock](https://github.com/jonboulle/clockwork) instance, ready to inject in your code.

This is particularly useful if you need to control time (set time, fast-forward, ...).

For example:

```go
package service

import (
	"time"

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

### Testing

This module provides a [*clockwork.FakeClock](https://github.com/jonboulle/clockwork) instance, that will be automatically injected as `clockwork.Clock` in your constructors in `test` mode.

#### Global time

By default, the fake clock is set to `time.Now()` (your test execution time).

You can configure the global time in your test in your testing configuration file (for all your tests), in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339) format:

```yaml
# ./configs/config_test.yaml
modules:
  clock:
    test:
      time: "2006-01-02T15:04:05Z07:00" # time in RFC3339 format
```

You can also [override this value](https://ankorstore.github.io/yokai/modules/fxconfig/#env-var-substitution), per test, by setting the `MODULES_CLOCK_TEST_TIME` env var.

#### Time control

You can `populate` the [*clockwork.FakeClock](https://github.com/jonboulle/clockwork) from your test to control time:

```go
package service_test

import (
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxsql"
	"github.com/foo/bar/internal/service"
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"
)

func TestExampleService(t *testing.T) {
	testTime := "2025-03-30T12:00:00Z"
	expectedTime, _ := time.Parse(time.RFC3339, testTime)

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

See [tests example](module_test.go) for more details.
