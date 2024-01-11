# Fx Config Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxconfig-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxconfig-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxconfig)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxconfig)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxconfig)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxconfig)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxconfig)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxconfig)](https://pkg.go.dev/github.com/ankorstore/yokai/fxconfig)

> [Fx](https://uber-go.github.io/fx/) module for [config](https://github.com/ankorstore/yokai/tree/main/config).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
  * [Loading](#loading)
  * [Configuration files](#configuration-files)
  * [Configuration usage](#configuration-usage)
  * [Override](#override)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxconfig
```

## Documentation

### Loading

To load the module in your Fx application:

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,                          // load the module
		fx.Invoke(func(cfg *config.Config) {              // invoke the config
			fmt.Printf("app name: %s", cfg.AppName())
		}),
	).Run()
}
```

### Configuration files

The module expects configuration files to be present:
- in `.` (project root)
- or in the`./configs` directory
- or any directory referenced in the `APP_CONFIG_PATH` env var

Check the [configuration files documentation](https://github.com/ankorstore/yokai/tree/main/config#configuration-files) for more details.

### Configuration usage

This module offers several features, such as:
- config helpers and typed accessors
- config dynamic environment overrides
- config env vars placeholders and runtime substitution

Check the [configuration usage documentation](https://github.com/ankorstore/yokai/tree/main/config#configuration-usage) for more details.

### Override

By default, the `config.Config` is created by the [DefaultConfigFactory](https://github.com/ankorstore/yokai/blob/main/config/factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"go.uber.org/fx"
)

type CustomConfigFactory struct{}

func NewCustomConfigFactory() config.ConfigFactory {
	return &CustomConfigFactory{}
}

func (f *CustomConfigFactory) Create(options ...config.ConfigOption) (*config.Config, error) {
	return &config.Config{...}, nil
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,                                 // load the module
		fx.Decorate(NewCustomConfigFactory),                     // decorate the module with a custom factory
		fx.Invoke(func(cfg *config.Config) {                     // invoke the custom config
			fmt.Printf("custom app name: %s", cfg.AppName())
		}),
	).Run()
}
```
