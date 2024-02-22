# Log Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/log-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/log-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/log)](https://goreportcard.com/report/github.com/ankorstore/yokai/log)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=log)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/log)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Flog)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/log)](https://pkg.go.dev/github.com/ankorstore/yokai/log)

> Logging module based on [Zerolog](https://github.com/rs/zerolog).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
  * [Usage](#usage)
  * [Context](#context)
  * [Testing](#testing)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/log
```

## Documentation

This module provides a [Logger](logger.go), offering all [Zerolog](https://github.com/rs/zerolog) methods.

### Usage

To create a `Logger`:

```go
package main

import (
	"os"

	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
)

var logger, _ = log.NewDefaultLoggerFactory().Create()

// equivalent to:
var logger, _ = log.NewDefaultLoggerFactory().Create(
	log.WithServiceName("default"),   // adds {"service":"default"} to log records
	log.WithLevel(zerolog.InfoLevel), // logs records with level >= info
	log.WithOutputWriter(os.Stdout),  // sends logs records to stdout
)
```

To use the `Logger`:

```go
package main

import (
	"github.com/ankorstore/yokai/log"
)

func main() {
	logger, _ := log.NewDefaultLoggerFactory().Create()

	logger.Info().Msg("some message") // {"level:"info", "service":"default", "message":"some message"}
}
```

See [Zerolog](https://github.com/rs/zerolog) documentation for more details about available methods.

### Context

This module provides the `log.CtxLogger()` function that allow to extract the logger from a `context.Context`.

If no logger is found in context, a [default](https://github.com/rs/zerolog/blob/master/ctx.go) Zerolog based logger will be used.

### Testing

This module provides a [TestLogBuffer](logtest/buffer.go), recording log records to be able to assert on them after logging:

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
)

func main() {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, _ := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(buffer))

	logger.Info().Msg("some message example")

	// test on attributes exact matching
	hasRecord, _ := buffer.HasRecord(map[string]interface{}{
		"level":   "info",
		"message": "some message example",
	})

	fmt.Printf("has record: %v", hasRecord) // has record: true

	// test on attributes partial matching
	containRecord, _ := buffer.ContainRecord(map[string]interface{}{
		"level":   "info",
		"message": "message",
	})

	fmt.Printf("contain record: %v", containRecord) // contain record: true
}
```

You can also use the provided [test assertion helpers](logtest/assert.go) in your tests:
- `AssertHasLogRecord`: to assert on exact attributes match
- `AssertHasNotLogRecord`: to assert on exact attributes non match
- `AssertContainLogRecord`: to assert on partial attributes match
- `AssertContainNotLogRecord`: to assert on partial attributes non match

and use `Dump()` to print the current content of the [TestLogBuffer](logtest/buffer.go).

For example:

```go
package main_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
)

func TestLogger(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()
	
	logger, _ := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(buffer))

	logger.Info().Msg("some message example")

	// print records
	buffer.Dump()
	
	// assertion success
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"message": "some message example",
	})

	// assertion success
	logtest.AssertHasNotLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"message": "some invalid example",
	})

	// assertion success
	logtest.AssertContainLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"message": "message",
	})

	// assertion success
	logtest.AssertContainNotLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"message": "invalid",
	})
}
```
