---
title: Modules - SQL
icon: material/cube-outline
---

# :material-cube-outline: SQL Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxsql-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxsql-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxsql)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxsql)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxsql)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxsql)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxsql)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxsql)](https://pkg.go.dev/github.com/ankorstore/yokai/fxsql)

## Overview

Yokai provides a [fxsql](https://github.com/ankorstore/yokai/tree/main/fxsql) module, allowing your application to interact with databases.

It wraps the [sql](https://github.com/ankorstore/yokai/tree/main/sql) module, based on [database/sql](https://pkg.go.dev/database/sql).

It comes with:

- automatic SQL operations logging and tracing
- possibility to define and execute database migrations (based on [Goose](https://github.com/pressly/goose))
- possibility to register and execute database seeds
- possibility to register database hooks around the SQL operations

Since this module enables you to work with `sql.DB`, you keep full control on your database interactions with `SQL`, and you can enhance your developer experience with tools like [SQLC](https://sqlc.dev/).

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxsql
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxsql"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxsql module
	fxsql.FxSQLModule,
	// ...
)
```

## Configuration

This module provides the possibility to configure the database `driver`:

- `mysql` for MySQL databases (based on [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql))
- `postgres` for PostgreSQL databases (based on [lib/pq](https://github.com/lib/pq))
- `sqlite` for SQLite databases (based on [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))

```yaml title="configs/config.yaml"
modules:
  sql:
    driver: mysql                                               # database driver
    dsn: "user:password@tcp(localhost:3306)/db?parseTime=true"  # database DSN
    migrations:
      path: db/migrations  # migrations path (empty by default)
      stdout: true         # to print in stdout the migration logs (disabled by default)
    log:
      enabled: true        # to enable SQL queries logging (disabled by default)
      level: debug         # to configure SQL queries logs level (debug by default)
      arguments: true      # to add SQL queries arguments to logs (disabled by default)
      exclude:             # to exclude SQL operations from logging (empty by default)
        - "connection:ping"
        - "connection:reset-session"
    trace:
      enabled: true        # to enable SQL queries tracing (disabled by default)
      arguments: true      # to add SQL queries arguments to trace spans (disabled by default)
      exclude:             # to exclude SQL operations from tracing (empty by default)
        - "connection:ping"
        - "connection:reset-session"
```

You can find below the list of supported `SQL operations`:

=== "Connection"
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

=== "Statement"
    - `statement:exec`
    - `statement:exec-context`
    - `statement:query`
    - `statement:query-context`

=== "Transaction"
    - `transaction:commit`
    - `transaction:rollback`

## Usage

Installing this module will automatically make a configured `sql.DB` instance available in Yokai dependency injection system.

To access it, you just need to inject it where needed, for example in a repository:

```go title="internal/repository/foo.go"
package repository

import (
	"context"
	"database/sql"
)

type FooRepository struct {
	db *sql.DB
}

func NewFooRepository(db *sql.DB) *FooRepository {
	return &FooRepository{
		db: db,
	}
}

func (r *FooRepository) Insert(ctx context.Context, bar string) (sql.Result, error) {
	return r.db.ExecContext(ctx, "INSERT INTO foo (bar) VALUES ?", bar)
}
```

Like any other services, the `FooRepository` needs to be registered to have its dependencies autowired:

```go title="internal/register.go"
package internal

import (
	"github.com/foo/bar/internal/repository"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the FooRepository
		fx.Provide(repository.NewFooRepository),
		// ...
	)
}
```

## Migrations

This module provides the possibility to run your `database migrations`, using [Goose](https://github.com/pressly/goose) under the hood.

### Migrations creation

You can configure where to find your migration files:

```yaml title="configs/config.yaml"
modules:
  sql:
    migrations: db/migrations
```

And create them following [Goose SQL migrations](https://github.com/pressly/goose?tab=readme-ov-file#sql-migrations) conventions:

```sql title="db/migrations/00001_create_foo_table.sql"
-- +goose Up
CREATE TABLE foo (
	id  INTEGER NOT NULL PRIMARY KEY,
	bar VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS foo;

```

### Migrations execution

This is done via:

- `RunFxSQLMigration(command, args)` to execute a migration command
- `RunFxSQLMigrationAndShutdown(command, args)` to execute a migration command and shut down

Available `migration commands`:

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

#### At bootstrap

To run the migrations automatically at bootstrap, you just need to call `RunFxSQLMigration()`:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxsql"
)

// ...

func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp(
		// run database migrations
		fxsql.RunFxSQLMigration("up"),
		// ...
	)
}

func RunTest(tb testing.TB, options ...fx.Option) {
	// ...

	Bootstrapper.RunTestApp(
		tb,
		// test options
		fx.Options(options...),
		// run database migrations for tests
		fxsql.RunFxSQLMigration("up"),
		// ...
	)
}
```

#### Dedicated command

A preferable way to run migrations is via a dedicated command.

You can create it in the `cmd/` directory of your application:

```go title="cmd/migrate.go"
package cmd

import (
    "github.com/ankorstore/yokai/fxcore"
    "github.com/ankorstore/yokai/fxsql"
    "github.com/spf13/cobra"
    "go.uber.org/fx"
)

func init() {
    rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Run database migrations",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        fxcore.
            NewBootstrapper().
            WithContext(cmd.Context()).
            WithOptions(
                fx.NopLogger,
                // modules
                fxsql.FxSQLModule,
                // migrate and shutdown
                fxsql.RunFxSQLMigrationAndShutdown(args[0], args[1:]...),
            ).
            RunApp()
    },
}

```

You can then execute this command when needed by running for example `app migrate up` from a dedicated step in your deployment pipeline.

## Seeds

This module provides the possibility to `seed` your database, useful for testing.

### Seeds creation

This module provides the [Seed](https://github.com/ankorstore/yokai/blob/main/fxsql/seeder.go) interface for your seeds implementations.

For example:

```go title="db/seeds/example.go"
package seeds

import (
    "context"
    "database/sql"

    "github.com/ankorstore/yokai/config"
)

type ExampleSeed struct {
    config *config.Config
}

func NewExampleSeed(config *config.Config) *ExampleSeed {
    return &ExampleSeed{
		config: config,
    }
}

func (s *ExampleSeed) Name() string {
    return "example-seed"
}

func (s *ExampleSeed) Run(ctx context.Context, db *sql.DB) error {
    _, err := db.ExecContext(
		ctx,
		"INSERT INTO foo (bar) VALUES (?)",
		s.config.GetString("config.seeds.example-seed.value"),
	)

    return err
}
```

### Seeds registration

Once your seeds are created, you can register them via:

- `AsSQLSeed()` to register a seed
- `AsSQLSeeds()` to register several seeds at once

```go title="internal/register.go"
package internal

import (
	"github.com/foo/bar/db/seeds"
    "github.com/ankorstore/yokai/fxsql"
    "go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the ExampleSeed
        fxsql.AsSQLSeed(seeds.NewExampleSeed),
		// ...
	)
}
```

The dependencies of your seeds constructors will be autowired.

### Seeds execution

Once your seeds are registered, you can execute them via `RunFxSQLSeeds()`:

```go title="internal/example_test.go"
package internal_test

import (
    "testing"

    "github.com/ankorstore/yokai/fxsql"
    "github.com/foo/bar/internal"
    "go.uber.org/fx"
)

func TestExample(t *testing.T) {
    internal.RunTest(
        t,
		// apply seeds
        fxsql.RunFxSQLSeeds(),
    )

    // ...
}
```

You can also call for example `RunFxSQLSeeds("example-seed", "other-seed")` to run only specific seeds, in provided order.

## Hooks

This module provides the possibility to `extend` the logic around the `SQL operations` via a hooking mechanism.

### Hooks creation

This module provides the [Hook](https://github.com/ankorstore/yokai/blob/main/sql/hook.go) interface for your hooks implementations.

For example:

```go title="db/hooks/example.go"
package hooks

import (
    "context"

    "github.com/ankorstore/yokai/config"
    "github.com/ankorstore/yokai/sql"
)

type ExampleHook struct {
    config *config.Config
}

func NewExampleHook(config *config.Config) *ExampleHook {
    return &ExampleHook{
		config: config,
    }
}

func (h *ExampleHook) Before(ctx context.Context, event *sql.HookEvent) context.Context {
    // before SQL operation logic
	if h.config.GetBool("config.hooks.example-hook.enabled") {
		// ...
    }

    return ctx
}

func (h *ExampleHook) After(ctx context.Context, event *sql.HookEvent) {
    // after SQL operation logic
    if h.config.GetBool("config.hooks.example-hook.enabled") {
        // ...
    }
}
```

### Hooks registration

Once your hooks are created, you can register them via:

- `AsSQLHook()` to register a hook
- `AsSQLHooks()` to register several hooks at once

```go title="internal/register.go"
package internal

import (
	"github.com/foo/bar/db/hooks"
    "github.com/ankorstore/yokai/fxsql"
    "go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the ExampleHook
        fxsql.AsSQLHook(hooks.NewExampleHook),
		// ...
	)
}
```

The dependencies of your hooks constructors will be autowired.

### Hooks execution

Yokai collects all registered hooks and executes them `automatically` on each SQL operations.

## Health Check

This module provides a ready to use [SQLProbe](https://github.com/ankorstore/yokai/blob/main/sql/healthcheck/probe.go), to be used by the [health check](fxhealthcheck.md) module.

It will perform a `ping` to the configured database connection to ensure it is healthy.

You just need to register it:

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/sql/healthcheck"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the SQLProbe probe for startup, liveness and readiness checks
		fxhealthcheck.AsCheckerProbe(healthcheck.NewSQLProbe),
		// ...
	)
}
```

## Logging

You can enable the SQL queries automatic logging with `modules.sql.log.enabled=true`:

```yaml title="configs/config.yaml"
modules:
  sql:
    log:
      enabled: true    # to enable SQL queries logging (disabled by default)
      level: debug     # to configure SQL queries logs level (debug by default)
      arguments: true  # to add SQL queries arguments to logs (disabled by default)
      exclude:         # to exclude SQL operations from logging (empty by default)
        - "connection:ping"
        - "connection:reset-session"
```

As a result, in your application logs:

```
DBG system:"mysql" operation:"connection:exec-context" latency="54.32µs" query="INSERT INTO foo (bar) VALUES (?)" lastInsertId=0 rowsAffected=0
```

If needed, you can log the SQL queries arguments with `modules.sql.log.arguments=true`:

```
DBG system:"mysql" operation:"connection:exec-context" latency="54.32µs" query="INSERT INTO foo (bar) VALUES (?)" arguments=[map[Name: Ordinal:1 Value:baz]] lastInsertId=0 rowsAffected=0
```

## Tracing

You can enable the SQL queries automatic tracing with `modules.sql.trace.enabled=true`:

```yaml title="configs/config.yaml"
modules:
  sql:
    trace:
      enabled: true    # to enable SQL queries tracing (disabled by default)
      arguments: true  # to add SQL queries arguments to trace spans (disabled by default)
      exclude:         # to exclude SQL operations from tracing (empty by default)
        - "connection:ping"
        - "connection:reset-session"
```

As a result, in your application trace spans attributes:

```
db.system: "mysql"
db.statement: "INSERT INTO foo (bar) VALUES (?)"
db.lastInsertId: 0
db.rowsAffected: 0
...
```

If needed, you can trace the SQL queries arguments with `modules.sql.trace.arguments=true`:

```
db.system: "mysql"
db.statement: "INSERT INTO foo (bar) VALUES (?)"
db.statement.arguments: "[{Name: Ordinal:1 Value:baz}]"
db.lastInsertId: 0
db.rowsAffected: 0
...
```

## Testing

This module provide support for the `sqlite` databases, making your tests portable (in memory, no database required):

```yaml title="configs/config.test.yaml"
modules:
  sql:
    driver: sqlite   # use sqlite driver
    dsn: ":memory:"  # in memory
```

You can then retrieve your components using the `sql.DB`, and make actual database operations:

```go title="internal/example_test.go"
package internal_test

import (
    "testing"

    "github.com/ankorstore/yokai/fxsql"
    "github.com/foo/bar/internal"
    "github.com/foo/bar/internal/repository"
    "go.uber.org/fx"
)

func TestExample(t *testing.T) {
    var fooRepository repository.FooRepository

    internal.RunTest(
        t,
		// apply migrations in sqlite in-memory
        fxsql.RunFxSQLMigration("up"),
		// apply seeds in sqlite in-memory
        fxsql.RunFxSQLSeeds(),
		// retrieve your components
        fx.Populate(&fooRepository),
    )

    // ...
}
```
