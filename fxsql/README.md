# Fx SQL Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxsql-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxsql-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxsql)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxsql)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxsql)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxsql)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxsql)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxsql)](https://pkg.go.dev/github.com/ankorstore/yokai/fxsql)

> [Fx](https://uber-go.github.io/fx/) module for [sql](https://github.com/ankorstore/yokai/tree/main/sql).

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
go get github.com/ankorstore/yokai/fxsql
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
	"database/sql"
	
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Model struct {
	Name string
}

func main() {
	fx.New(
		fxconfig.FxConfigModule, // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,       // load the module
		fx.Invoke(func(db *sql.DB) {
			// use the DB
		}),
	).Run()
}
```

TODO