# SQL Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/sql-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/sql-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/sql)](https://goreportcard.com/report/github.com/ankorstore/yokai/sql)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=sql)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/sql)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Fsql)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/sql)](https://pkg.go.dev/github.com/ankorstore/yokai/sql)

> SQL module based on [database/sql](https://pkg.go.dev/database/sql).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
  * [Usage](#usage)
  * [Hooks](#hooks)
    * [Log hook](#log-hook)
    * [Trace hook](#trace-hook)
    * [Custom hook](#custom-hook)
  * [Healthcheck](#healthcheck)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/sql
```

## Documentation

This module provides a [Driver](driver.go), decorating `database/sql` compatible drivers, with a [hooking mechanism](hook.go).

### Usage

The following database systems are [supported](system.go):

- `mysql` with [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- `postgres` with [lib/pq](https://github.com/lib/pq)
- `sqlite` with [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

To create a `*sql.DB` with the [tracing](hook/trace/hook.go) and [logging](hook/log/hook.go) hooks:

```go
package main

import (
	"database/sql"
	
	yokaisql "github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/ankorstore/yokai/sql/hook/trace"
)

func main() {
	// MySQL
	driver, _ := yokaisql.Register("mysql", trace.NewTraceHook(), log.NewLogHook())
	db, _ := sql.Open(driver, "user:password@tcp(localhost:3306)/db?parseTime=true")

	// Postgres
	driver, _ := yokaisql.Register("postgres", trace.NewTraceHook(), log.NewLogHook())
	db, _ := sql.Open(driver, "host=host port=5432 user=user password=password dbname=db sslmode=disable")

	// SQLite
	driver, _ := yokaisql.Register("sqlite", trace.NewTraceHook(), log.NewLogHook())
	db, _ := sql.Open(driver, ":memory:")
}
```

See [database/sql](https://pkg.go.dev/database/sql) documentation for more details.

### Hooks

This module provides a [hooking mechanism](hook.go) to add logic around the [SQL operations](operation.go).

#### Log hook

This module provides an [LogHook](hook/log/hook.go), that you can use to automatically `log` the [SQL operations](operation.go):

```go
package main

import (
	"database/sql"

	yokaisql "github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/rs/zerolog"
)

func main() {

	logHook := log.NewLogHook(
		log.WithLevel(zerolog.DebugLevel),        // SQL logs level, debug by default
		log.WithArguments(true),                  // SQL logs with SQL arguments, false by default
		log.WithExcludedOperations(               // SQL operations to exclude from logging, empty by default
			yokaisql.ConnectionPingOperation,
			yokaisql.ConnectionResetSessionOperation,
		),
	)

	driver, _ := yokaisql.Register("sqlite", logHook)
	db, _ := sql.Open(driver, ":memory:")
}
```

#### Trace hook

This module provides an [TraceHook](hook/trace/hook.go), that you can use to automatically `trace` the [SQL operations](operation.go):

```go
package main

import (
	"database/sql"

	yokaisql "github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/ankorstore/yokai/sql/hook/trace"
	"github.com/rs/zerolog"
)

func main() {

	traceHook := trace.NewTraceHook(
		trace.WithArguments(true),                  // SQL traces with SQL arguments, false by default
		trace.WithExcludedOperations(               // SQL operations to exclude from tracing, empty by default
			yokaisql.ConnectionPingOperation,
			yokaisql.ConnectionResetSessionOperation,
		),
	)

	driver, _ := yokaisql.Register("sqlite", traceHook)
	db, _ := sql.Open(driver, ":memory:")
}
```

#### Custom hook

This module provides a [Hook](hook.go) interface, that you can implement to extend the logic around [SQL operations](operation.go):

```go
package main

import (
	"context"
	"database/sql"

	yokaisql "github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/ankorstore/yokai/sql/hook/trace"
	"github.com/rs/zerolog"
)

type CustomHook struct{}

func (h *CustomHook) Before(ctx context.Context, event *yokaisql.HookEvent) context.Context {
	// your custom logic before SQL operation
	
	return ctx
}

func (h *CustomHook) After(ctx context.Context, event *yokaisql.HookEvent) {
	// your custom logic after SQL operation
}

func main() {
	driver, _ := yokaisql.Register("sqlite", &CustomHook{})
	db, _ := sql.Open(driver, ":memory:")
}
```

### Healthcheck

This module provides an [SQLProbe](healthcheck/probe.go), compatible with
the [healthcheck module](https://github.com/ankorstore/yokai/tree/main/healthcheck):

```go
package main

import (
	"context"

	yokaihc "github.com/ankorstore/yokai/healthcheck"
	yokaisql "github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/healthcheck"
)

func main() {
	driver, _ := yokaisql.Register("sqlite")
	db, _ := sql.Open(driver, ":memory:")

	checker, _ := yokaihc.NewDefaultCheckerFactory().Create(
		yokaihc.WithProbe(healthcheck.NewSQLProbe(db)),
	)

	checker.Check(context.Background(), yokaihc.Readiness)
}
```

This probe performs a `ping` to the configured database connection.
