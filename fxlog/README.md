# Fx Log Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxlog-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxlog-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxlog)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxlog)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxlog)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxlog)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxlog)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxlog)](https://pkg.go.dev/github.com/ankorstore/yokai/fxlog)

> [Fx](https://uber-go.github.io/fx/) module for [log](https://github.com/ankorstore/yokai/tree/main/log).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
	* [Dependencies](#dependencies)
	* [Loading](#loading)
	* [Configuration](#configuration)
	* [Override](#override)
	* [Testing](#testing)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxlog
```

## Documentation

### Dependencies

This module is intended to be used alongside the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module.

### Loading

To load the module in your Fx application:

```go
package main

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/log"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule, // load the module dependency
		fxlog.FxLogModule,       // load the module
		fx.Invoke(func(logger *log.Logger) { // invoke the logger
			logger.Info().Msg("message")
		}),
	).Run()
}
```

If needed, you can also configure [Fx](https://uber-go.github.io/fx/) to use this logger for its own event logs:

```go
package main

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,               // load the module dependency
		fxlog.FxLogModule,                     // load the module
		fx.WithLogger(fxlog.NewFxEventLogger), // configure Fx event logging
	).Run()
}
```

### Configuration

This module provides the possibility to configure:

- the `log level` (possible values: `trace`, `debug`, `info`, `warning`, `error`, `fatal`, `panic`, `no-level` or `disabled`)
- the `log output` (possible values: `noop`, `stdout` or `test`)

Regarding the output:

- `stdout`: to send the log records to `os.Stdout` (default)
- `noop`: to void the log records via `os.Discard`
- `console`: [pretty prints](https://github.com/rs/zerolog#pretty-logging) logs record to `os.Stdout`
- `test`: to send the log records to the [TestLogBuffer](https://github.com/ankorstore/yokai/blob/main/log/logtest/buffer.go) made available in the Fx container, for further assertions

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  log:
    level: info    # by default
    output: stdout # by default
```

Notes:

- the config `app.name` (or env var `APP_NAME`) will be used in each log record `service` field: `{"service":"app"}`
- if the config `app.debug=true` (or env var `APP_DEBUG=true`), the `debug` level will be used, no matter given configuration
- if the config `app.env=test` (or env var `APP_ENV=test`), the `test` output will be used, no matter given configuration

### Override

By default, the `log.Logger` is created by the [DefaultLoggerFactory](https://github.com/ankorstore/yokai/blob/main/log/factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/log"
	"go.uber.org/fx"
)

type CustomLoggerFactory struct{}

func NewCustomLoggerFactory() log.LoggerFactory {
	return &CustomLoggerFactory{}
}

func (f *CustomLoggerFactory) Create(options ...log.LoggerOption) (*log.Logger, error) {
	return &log.Logger{...}, nil
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,             // load the module dependency
		fxlog.FxLogModule,                   // load the module
		fx.Decorate(NewCustomLoggerFactory), // override the module with a custom factory
		fx.Invoke(func(logger *log.Logger) { // invoke the custom logger
			logger.Info().Msg("custom message")
		}),
	).Run()
}
```

### Testing

This module provides the possibility to easily test your log records, using the [TestLogBuffer](https://github.com/ankorstore/yokai/blob/main/log/logtest/buffer.go) with `modules.log.output=test`.

```yaml
# ./configs/config.test.yaml
app:
  name: test
modules:
  log:
    output: test # to send logs to test buffer
```

You can then test:

```go
package main_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestLogger(t *testing.T) {
	t.Setenv("APP_NAME", "test")
	t.Setenv("APP_ENV", "test")

	var buffer logtest.TestLogBuffer

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fx.Invoke(func(logger *log.Logger) {
			logger.Debug().Msg("test message")
		}),
		fx.Populate(&buffer), // extracts the TestLogBuffer from the Fx container
	).RequireStart().RequireStop()

	// assertion success
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "debug",
		"service": "test",
		"message": "test message",
	})
}
```

See the `log` module testing [documentation](https://github.com/ankorstore/yokai/tree/main/log#testing) for more details.
