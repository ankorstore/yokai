# ORM Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/orm-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/orm-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/orm)](https://goreportcard.com/report/github.com/ankorstore/yokai/orm)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=orm)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/orm)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Form)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/orm)](https://pkg.go.dev/github.com/ankorstore/yokai/orm)

> ORM module based on [GORM](https://gorm.io/).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
	* [Usage](#usage)
	* [Add-ons](#add-ons)
		* [Logger](#logger)
		* [Tracer](#tracer)
		* [Healthcheck](#healthcheck)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/orm
```

## Documentation

### Usage

This module provides a [OrmFactory](factory.go), allowing to build an `gorm.DB` instance.

The following database drivers are [supported](https://gorm.io/docs/connecting_to_the_database.html):

- `SQLite`
- `MySQL`
- `PostgreSQL`
- `SQL Server`

```go
package main

import (
	"github.com/ankorstore/yokai/orm"
)

// with MySQL driver
var db, _ = orm.NewDefaultOrmFactory().Create(
	orm.WithDriver(orm.Mysql),
	orm.WithDsn("user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=True"),
)

// or with SQLite driver
var db, _ = orm.NewDefaultOrmFactory().Create(
	orm.WithDriver(orm.Sqlite),
	orm.WithDsn("file::memory:?cache=shared"),
)
```

See [GORM documentation](https://gorm.io/docs/) for more details.

### Add-ons

This module provides several add-ons ready to use to enrich your ORM.

#### Logger

This module provides an [CtxOrmLogger](logger.go), compatible with
the [log module](https://github.com/ankorstore/yokai/tree/main/log):

```go
package main

import (
	"github.com/ankorstore/yokai/orm"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	ormLogger := orm.NewCtxOrmLogger(logger.Info, false)

	db, _ := orm.NewDefaultOrmFactory().Create(
		orm.WithConfig(gorm.Config{
			Logger: ormLogger,
		}),
	)
}
```

If needed, you can set the parameter `withValues` to `true` to append SQL query parameter values in the log records:

```go
ormLogger := orm.NewCtxOrmLogger(logger.Info, true)
```

#### Tracer

This module provides an [OrmTracerPlugin](plugin/trace.go), compatible with
the [trace module](https://github.com/ankorstore/yokai/tree/main/trace):

```go
package main

import (
	"github.com/ankorstore/yokai/orm"
	"github.com/ankorstore/yokai/orm/plugin"
	"github.com/ankorstore/yokai/trace"
)

func main() {
	tracerProvider, _ := trace.NewDefaultTracerProviderFactory().Create()

	db, _ := orm.NewDefaultOrmFactory().Create()

	db.Use(plugin.NewOrmTracerPlugin(tracerProvider, false))
}
```

If needed, you can set the parameter `withValues` to `true` to append SQL query parameter values in the trace spans:

```go
db.Use(plugin.NewOrmTracerPlugin(tracerProvider, true))
```

#### Healthcheck

This module provides an [OrmProbe](healthcheck/probe.go), compatible with
the [healthcheck module](https://github.com/ankorstore/yokai/tree/main/healthcheck):

```go
package main

import (
	"context"

	hc "github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/orm"
	"github.com/ankorstore/yokai/orm/healthcheck"
)

func main() {
	db, _ := orm.NewDefaultOrmFactory().Create()

	checker, _ := hc.NewDefaultCheckerFactory().Create(
		hc.WithProbe(healthcheck.NewOrmProbe(db)),
	)

	checker.Check(context.Background(), hc.Readiness)
}
```

This probe performs a `ping` to the configured database connection.
