---
title: Modules - Log
icon: material/cube-outline
---

# :material-cube-outline: Log Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxlog-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxlog-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxlog)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxlog)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxlog)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxlog)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxlog)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxlog)](https://pkg.go.dev/github.com/ankorstore/yokai/fxlog)

## Overview

Yokai provides a [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module, allowing your application to produce logs.

It wraps the [log](https://github.com/ankorstore/yokai/tree/main/log) module, based on [Zerolog](https://github.com/rs/zerolog).

## Installation

The [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module is automatically loaded by Yokai's [core](fxcore.md).

When you use a Yokai `application template`, you have nothing to install, it's ready to use.

## Configuration

This module provides the possibility to configure:

- the `log level` (possible values: `trace`, `debug`, `info`, `warning`, `error`, `fatal`, `panic`, `no-level` or `disabled`)
- the `log output` (possible values: `noop`, `stdout` or `test`)

Regarding the output:

- `stdout`: to send the log records to `os.Stdout` (default)
- `noop`: to void the log records via `os.Discard`
- `console`: [pretty prints](https://github.com/rs/zerolog#pretty-logging) logs record to `os.Stdout`
- `test`: to send the log records to the [TestLogBuffer](https://github.com/ankorstore/yokai/blob/main/log/logtest/buffer.go) made available in the Fx container, for further assertions

```yaml title="configs/config.yaml"
modules:
  log:
    level: info    # by default
    output: stdout # by default
```

## Usage

This module makes available the [Logger](https://github.com/ankorstore/yokai/blob/main/log/logger.go) in
Yokai dependency injection system.

It is built on top of `Zerolog`, see its [documentation](https://github.com/rs/zerolog) for more details about available methods.

You can inject the logger where needed, but it's recommended to use the one carried by the `context.Context` when possible (for automatic logs correlation).

## Testing

This module provides the possibility to easily test your application logs, using the [TestLogBuffer](https://github.com/ankorstore/yokai/blob/main/log/logtest/buffer.go) with `modules.log.output=test`.

```yaml title="configs/config.test.yaml"
modules:
  log:
    output: test # to send logs to test buffer
```

You can use the provided [test assertion helpers](https://github.com/ankorstore/yokai/blob/main/log/logtest/assert.go) in your tests:

- `AssertHasLogRecord`: to assert on exact attributes match
- `AssertHasNotLogRecord`: to assert on exact attributes non match
- `AssertContainLogRecord`: to assert on partial attributes match
- `AssertContainNotLogRecord`: to assert on partial attributes non match

and use `Dump()` to print the current content of the [TestLogBuffer](https://github.com/ankorstore/yokai/blob/main/log/logtest/buffer.go).

For example:

```go title="internal/example_test.go"
package internal_test

import (
	"testing"
	
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/foo/bar/internal"
	"go.uber.org/fx"
)

func TestExample(t *testing.T) {
	var logBuffer logtest.TestLogBuffer

	internal.RunTest(
		t,
		fx.Populate(&logBuffer),
		fx.Invoke(func(logger *log.Logger) {
			logger.Debug().Msg("test message")
		}),
	)
	
	// print logs records
	logBuffer.Dump()
	
	// log assertion example
	logtest.AssertHasLogRecord(
		t,
		logBuffer,
		map[string]interface{}{
			"level":   "debug",
			"message": "test message",
		}, 
	)
}
```
