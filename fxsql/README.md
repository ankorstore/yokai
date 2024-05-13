# Fx SQL Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxsql-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxsql-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxsql)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxsql)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxsql)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxsql)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxsql)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxsql)](https://pkg.go.dev/github.com/ankorstore/yokai/fxsql)

> [Fx](https://uber-go.github.io/fx/) module for [sql](https://github.com/ankorstore/yokai/tree/main/sql).

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
  * [Dependencies](#dependencies)
  * [Loading](#loading)
  * [Configuration](#configuration)
  * [Migrations](#migrations)
  * [Seeds](#seeds)
  * [Hooks](#hooks)
  * [Testing](#testing)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxsql
```

## Features

This module provides a `*sql.DB` to your Fx application with:

- automatic SQL requests logging and tracing
- possibility to define and apply database migrations (based on [Goose](https://github.com/pressly/goose))
- possibility to register database hooks
- possibility to register and execute database seeds

## Documentation

### Dependencies

This module is intended to be used alongside:

- the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module
- the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module
- the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module

### Loading

To load the module in your Fx application:

```go
package main

import (
	"database/sql"
	
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule, // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,       // load the module
		fx.Invoke(func(db *sql.DB) {
			// use the DB
			res, _ := db.Exec(...)
		}),
	).Run()
}
```

### Configuration

This module provides the possibility to configure the SQL `driver`:

- `mysql` with [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- `postgres` with [lib/pq](https://github.com/lib/pq)
- `sqlite` with [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

Configuration reference:

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  sql:
    driver: mysql                                               # database driver (empty by default)
    dsn: "user:password@tcp(localhost:3306)/db?parseTime=true"  # database DSN (empty by default)
    migrations: db/migrations                                   # migrations path (empty by default)
    log:
      enabled: true             # to enable SQL queries logging (disabled by default)
      level: debug              # to configure SQL queries logs level (debug by default)
      arguments: true           # to add SQL queries arguments to logs (disabled by default)
      exclude:                  # to exclude SQL operations from logging (empty by default)
        - "connection:ping"
    trace:
      enabled: true             # to enable SQL queries tracing (disabled by default)
      arguments: true           # to add SQL queries arguments to trace spans (disabled by default)
      exclude:                  # to exclude SQL operations from tracing (empty by default)
        - "connection:ping"
```

Available SQL operations:

- `connection:begin`
- `connection:begin-tx`
- `connection:exec`
- `connection:exec-context`
- `connection:query`
- `connection:query-context`
- `connection:prepare`
- `connection:prepare-context`
- `connection:ping`
- `connection:reset-session`
- `connection:close`
- `statement:exec`
- `statement:exec-context`
- `statement:query`
- `statement:query-context`
- `transaction:commit`
- `transaction:rollback`

### Migrations

This module provides the possibility to run your DB schemas migrations, using [Goose](https://github.com/pressly/goose) under the hood.

First, configure the path for your migration files:

```yaml
# ./configs/config.yaml
modules:
  sql:
    migrations: db/migrations
```

You can then create a migration file in this path.

For example, `db/migrations/00001_create_foo_table.sql`:

```sql
-- +goose Up
CREATE TABLE foo (
    id  INTEGER NOT NULL PRIMARY KEY,
    bar VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS foo;
```

You can then run the migrations with `RunFxSQLMigration`, by specifying a migration command:

```go
package main

import (
	"database/sql"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,       // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,             // load the module
		fxsql.RunFxSQLMigration("up"), // run migration command "up"
	).Run()
}
```

Available migration commands:

- `up`: migrate the DB to the most recent version available
- `up-by-one`: migrate the DB up by 1
- `up-to VERSION`: migrate the DB to a specific VERSION
- `down`: roll back the version by 1
- `down-to VERSION`: roll back to a specific VERSION
- `redo`: re-run the latest migration
- `reset`: roll back all migrations
- `status`: dump the migration status for the current DB
- `version`: print the current version of the database
- `create NAME [sql|go]`: creates new migration file with the current timestamp
- `fix`: apply sequential ordering to migrations
- `validate`: check migration files without running them

If you want to automatically shut down your Fx application after the migrations, you can use `RunFxSQLMigrationAndShutdown`.

### Seeds

This module provides the possibility to register several [Seed](seeder.go) implementations to `seed` the database.

This is done via:

- the `AsSQLSeed()` function to register seeds
- the `RunFxSQLSeeds()` function to execute seeds

```go
package main

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
)

// example SQL seed
type ExampleSeed struct{}

func NewExampleSeed() *ExampleSeed {
	return &ExampleSeed{}
}

func (s *ExampleSeed) Name() string {
	return "example-seed"
}

func (s *ExampleSeed) Run(ctx context.Context, db *sql.DB) error {
  _, err := db.ExecContext(ctx, "INSERT INTO foo (bar) VALUES (?)", "baz")

  return err
}

// usage
func main() {
	fx.New(
		fxconfig.FxConfigModule,          // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,                // load the module
		fxsql.AsSQLHook(NewExampleSeed),  // register the ExampleSeed
		fxsql.RunFxSQLSeeds(),            // run all registered seeds
	).Run()
}
```

You can also use `AsSQLSeeds()` to register several seeds at once.

You can also call for example `RunFxSQLSeeds("example-seed", "other-seed")` to run specific seeds, in provided order.

The dependencies of your seeds constructors will be autowired.

### Hooks

This module provides the possibility to register several [Hook](https://github.com/ankorstore/yokai/blob/main/sql/hook.go) implementations to `extend` the logic around the SQL operations.

This is done via the `AsSQLHook()` function:

```go
package main

import (
	"context"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/sql"
	"go.uber.org/fx"
)

// example SQL hook
type ExampleHook struct{}

func NewExampleHook() *ExampleHook {
	return &ExampleHook{}
}

func (h *ExampleHook) Before(ctx context.Context, event *sql.HookEvent) context.Context {
	// before SQL operation logic

	return ctx
}

func (h *ExampleHook) After(ctx context.Context, event *sql.HookEvent) {
	// after SQL operation logic
}

// usage
func main() {
	fx.New(
		fxconfig.FxConfigModule,         // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,               // load the module
		fxsql.AsSQLHook(NewExampleHook), // register the ExampleHook
	).Run()
}
```

You can also use `AsSQLHooks()` to register several hooks at once.

The dependencies of your hooks constructors will be autowired.

### Testing

This module supports the `sqlite` driver, allowing the usage of an `in-memory` database, avoiding your tests to require a real database instance to run.

In your `testing` configuration:

```yaml
# ./configs/config.test.yaml
modules:
  sql:
    driver: sqlite
    dsn: ":memory:"
```

You can find tests examples in this module own [tests](module_test.go).