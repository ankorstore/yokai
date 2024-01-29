# ORM Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxorm-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxorm-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxorm)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxorm)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxorm)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxorm)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxorm)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxorm)](https://pkg.go.dev/github.com/ankorstore/yokai/fxorm)

## Overview

Yokai provides a [fxorm](https://github.com/ankorstore/yokai/tree/main/fxorm) module, allowing your application to interact with databases.

It wraps the [orm](https://github.com/ankorstore/yokai/tree/main/orm) module, based on [Gorm](https://gorm.io/).

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxorm
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxorm"
)


var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxorm module
	fxorm.FxOrmModule,
	// ...
)
```

## Usage

You can declare your [models](https://gorm.io/docs/models.html), for example:

```go title="internal/model/example.go"
package service

import (
	"gorm.io/gorm"
)

type ExampleModel struct {
	gorm.Model
	Name string
}
```

This module makes available the [DB](https://github.com/go-gorm/gorm/blob/master/gorm.go) in
Yokai dependency injection system.

To access it, you just need to inject it where needed, for example in a repository to manage your `ExampleModel`:

```go title="internal/respository/example.go"
package repository

import (
	"context"
	"sync"
	
	"github.com/foo/bar/internal/model"
	"gorm.io/gorm"
)

type ExampleRepository struct {
	mutex sync.Mutex
	db    *gorm.DB
}

func NewExampleRepository(db *gorm.DB) *ExampleRepository {
	return &ExampleRepository{
		db: db,
	}
}

func (r *ExampleRepository) Find(ctx context.Context, id int) (*model.ExampleModel, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	var e model.ExampleModel

	res := r.db.WithContext(ctx).Take(&e, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return &e, nil
}

func (r *ExampleRepository) Create(ctx context.Context, example *model.ExampleModel) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	res := r.db.WithContext(ctx).Create(example)

	return res.Error
}
```

Like any other services, the `ExampleRepository` needs to be registered to have its dependencies autowired:

```go title="internal/services.go"
package internal

import (
	"github.com/foo/bar/internal/repository"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// register the ExampleRepository
		fx.Provide(repository.NewExampleRepository),
		// ...
	)
}
```

## Configuration

This module provides the possibility to configure the database `driver`:

- `sqlite` for SQLite databases
- `mysql` for MySQL databases
- `postgres` for PostgreSQL databases
- `sqlserver` for SQL Server databases

You can also provide to the ORM the database`dsn`, some `config`, and configure SQL queries automatic `logging` and `tracing`.

```yaml title="configs/config.yaml"
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  orm:
    driver: mysql                                               # driver to use
    dsn: user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=True"   # database DSN to connect to
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

See [Gorm Config](https://github.com/go-gorm/gorm/blob/master/gorm.go) for more details about the ORM configuration.

## Migrations

This module provides the possibility to run automatically your [schemas migrations](https://gorm.io/docs/migration.html).

You just need to pass the list of models you want to auto migrate to `RunFxOrmAutoMigrate()`:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxorm"
	"github.com/foo/bar/internal/model"
)

//...

func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp(
		// run ORM migrations for the ExampleModel model
		fxorm.RunFxOrmAutoMigrate(&model.ExampleModel{}),
		// ...
	)
}

func RunTest(tb testing.TB, options ...fx.Option) {
	//...

	Bootstrapper.RunTestApp(
		tb,
		// run ORM migrations for the ExampleModel model for tests
		fxorm.RunFxOrmAutoMigrate(&model.ExampleModel{}),
		// ...
	)
}
```

## Performance

See general [Gorm performance recommendations](https://gorm.io/docs/performance.html).

### Disable Default Transaction

Gorm performs write (create/update/delete) operations by default inside a transaction to ensure data consistency, which
is not optimized for performance.

You can disable it in the configuration:

```yaml title="configs/config.yaml"
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  orm:
    driver: mysql                                               # driver to use
    dsn: user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=True"   # database DSN to connect to
    config:
      skip_default_transaction: true                            # disable default transaction
```

### Cache Prepared Statement

To create a prepared statement when executing any SQL (and cache them to speed up future calls):

```yaml title="configs/config.yaml"
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  orm:
    driver: mysql                                               # driver to use
    dsn: user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=True"   # database DSN to connect to
    config:
      prepare_stmt: true                                        # enable prepared statements
```

## Testing

This module provide support for the `sqlite` databases, making your tests portable (in memory, no database required):

```yaml title="configs/config.test.yaml"
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  orm:
    driver: sqlite   # use sqlite driver
    dsn: ":memory:"  # in memory
```

You can then retrieve your components using the [DB](https://github.com/go-gorm/gorm/blob/master/gorm.go), and make actual database operations:

```go title="internal/repository/example_test.go"
package repository_test

import (
	"testing"
	
	"github.com/foo/bar/internal/model"
	"github.com/foo/bar/internal/repository"
	"go.uber.org/fx"
)

func TestExampleRepository(t *testing.T) {
	var exampleRepository repository.ExampleRepository

	internal.RunTest(t, fx.Populate(&exampleRepository))

	// prepare your test data in the sqlite database
	exampleRepository.Create(
		context.Background(),
		&model.ExampleModel{
			Name: "test",
		},
	)
	
	// ...
}
```
