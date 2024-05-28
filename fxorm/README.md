# Fx ORM Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxorm-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxorm-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxorm)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxorm)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxorm)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxorm)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxorm)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxorm)](https://pkg.go.dev/github.com/ankorstore/yokai/fxorm)

> [Fx](https://uber-go.github.io/fx/) module for [orm](https://github.com/ankorstore/yokai/tree/main/orm).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
	* [Dependencies](#dependencies)
	* [Loading](#loading)
	* [Configuration](#configuration)
	* [Auto migrations](#auto-migrations)
	* [Performance](#performance)
		* [Disable Default Transaction](#disable-default-transaction)
		* [Cache Prepared Statement](#cache-prepared-statement)
	* [Override](#override)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxorm
```

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
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxorm"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Model struct {
	Name string
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,          // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxorm.FxOrmModule,                // load the module
		fx.Invoke(func(gormDB *gorm.DB) { // invoke the orm
			gormDB.Create(&Model{Name: "some name"})
		}),
	).Run()
}
```

### Configuration

This module provides the possibility to configure the ORM `driver`:

- `sqlite` for SQLite databases
- `mysql` for MySQL databases
- `postgres` for PostgreSQL databases
- `sqlserver` for SQL Server databases

You can also provide to the ORM the database`dsn`, some `config`, and configure SQL queries `logging` and `tracing`.

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  orm:
    driver: mysql                                               # driver to use
    dsn: "user:password@tcp(localhost:3306)/db?parseTime=True"  # database DSN to use
    config:
      dry_run: false                                            # disabled by default
      skip_default_transaction: false                           # disabled by default
      full_save_associations: false                             # disabled by default
      prepare_stmt: false                                       # disabled by default
      disable_automatic_ping: false                             # disabled by default
      disable_foreign_key_constraint_when_migrating: false      # disabled by default
      ignore_relationships_when_migrating: false                # disabled by default
      disable_nested_transaction: false                         # disabled by default
      allow_global_update: false                                # disabled by default
      query_fields: false                                       # disabled by default
      translate_error: false                                    # disabled by default
    log:
      enabled: true  # to log SQL queries, disabled by default
      level: info    # with a minimal level
      values: true   # by adding or not clear SQL queries parameters values in logs, disabled by default
    trace:
      enabled: true  # to trace SQL queries, disabled by default
      values: true   # by adding or not clear SQL queries parameters values in trace spans, disabled by default
```

See [GORM Config](https://github.com/go-gorm/gorm/blob/master/gorm.go) for more details about the `modules.orm.config` configuration keys.

For security reasons, you should avoid to hardcode DSN sensible parts (like the password) in your config files, you can use the [env vars placeholders](https://github.com/ankorstore/yokai/tree/main/fxconfig#configuration-env-var-placeholders) instead:

```yaml
# ./configs/config.yaml
modules:
  orm:
    driver: mysql
    dsn: "${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}?parseTime=True"
```

### Auto migrations

This module provides the possibility to run automatically your [schemas migrations](https://gorm.io/docs/migration.html)
with `RunFxOrmAutoMigrate()`:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxorm"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Model struct {
	Name string
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,             // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxorm.FxOrmModule,                   // load the module
		fxorm.RunFxOrmAutoMigrate(&Model{}), // run auto migration for Model
		fx.Invoke(func(gormDB *gorm.DB) {        // invoke the orm
			gormDB.Create(&Model{Name: "some name"})
		}),
	).Run()
}
```

### Performance

See [GORM performance recommendations](https://gorm.io/docs/performance.html).

#### Disable Default Transaction

Gorm performs write (create/update/delete) operations by default inside a transaction to ensure data consistency, which
is not optimized for performance.

You can disable it in the configuration:

```yaml
# ./configs/config.yaml

modules:
  orm:
    config:
      skip_default_transaction: true # disable default transaction
```

#### Cache Prepared Statement

To create a prepared statement when executing any SQL (and cache them to speed up future calls):

```yaml
# ./configs/config.yaml
modules:
  orm:
    config:
      prepare_stmt: true # enable prepared statements
```

### Override

By default, the `gorm.DB` is created by
the [DefaultOrmFactory](https://github.com/ankorstore/yokai/blob/main/orm/factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxorm"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/orm"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type CustomOrmFactory struct{}

func NewCustomOrmFactory() orm.OrmFactory {
	return &CustomOrmFactory{}
}

func (f *CustomOrmFactory) Create(options ...orm.OrmOption) (*gorm.DB, error) {
	return &gorm.DB{...}, nil
}

type Model struct {
	Name string
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,            // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxorm.FxOrmModule,                  // load the module
		fx.Decorate(NewCustomOrmFactory),   // override the module with a custom factory
		fx.Invoke(func(customDb *gorm.DB) { // invoke the custom ORM
			customDb.Create(&Model{Name: "custom"})
		}),
	).Run()
}
```
