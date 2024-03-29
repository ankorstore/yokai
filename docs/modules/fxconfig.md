---
icon: material/cube-outline
---

# :material-cube-outline: Config Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxconfig-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxconfig-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxconfig)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxconfig)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxconfig)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxconfig)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxconfig)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxconfig)](https://pkg.go.dev/github.com/ankorstore/yokai/fxconfig)

## Overview

Yokai provides a [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module, allowing to define and retrieve configurations for your application.

It wraps the [config](https://github.com/ankorstore/yokai/tree/main/config) module, based on [Viper](https://github.com/spf13/viper).

## Installation

The [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module is automatically loaded by Yokai's [core](fxcore.md).

When you use a Yokai `application template`, you have nothing to install, it's ready to use.

## Configuration files

By default, the module expects the configuration files:

- to be present in the `./configs` directory of your project
- to be named `config.{format}` (ex: `config.yaml`, `config.json`, etc.)
- to offer env overrides files named `config.{env}.{format}` based on the env var `APP_ENV` (ex: `config.test.yaml` if
  env var `APP_ENV=test`)

Supported configuration files formats: `.json`, `.toml`, `.yaml`, `.hcl`, `.ini`, and `.env`.

## Usage

For the following examples, we will be considering those configuration files:

```yaml title="configs/config.yaml"
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
config:
  values:
    string_value: default
    int_value: 0
  placeholder: foo-${BAR}-baz
  substitution: foo
```

and

```yaml title="configs/config.test.yaml"
app:
  env: test
  debug: true
config:
  values:
    string_value: test
```

### Configuration access

This module makes available the [Config](https://github.com/ankorstore/yokai/blob/main/config/config.go) in
Yokai dependency injection system.

It is built on top of `Viper`, see its [documentation](https://github.com/spf13/viper) for more details about available methods.

To access it, you just need to inject it where needed, for example:

```go title="internal/service/example.go"
package service

import (
	"fmt"

	"github.com/ankorstore/yokai/config"
)

type ExampleService struct {
	config *config.Config
}

func NewExampleService(config *config.Config) *ExampleService {
	return &ExampleService{
		config: config,
	}
}

func (s *ExampleService) PrintConfig() {
	// helpers
	fmt.Printf("name: %s", s.config.AppName())       // name: app
	fmt.Printf("env: %s", s.config.AppEnv())         // env: dev
	fmt.Printf("version: %s", s.config.AppVersion()) // version: 0.1.0
	fmt.Printf("debug: %v", s.config.AppDebug())     // debug: false

	// others
	fmt.Printf("string_value: %s", s.config.GetString("config.values.string_value")) // string_value: default
	fmt.Printf("int_value: %s", s.config.GetInt("config.values.int_value"))          // int_value: 0
}
```

### Dynamic env overrides

This module offers the possibility to override dynamically (by merging) configuration files depending on the env
var `APP_ENV` value.

For example, if `APP_ENV=test`, the module will use `config.yaml` values and merge / override them
with `config.test.yaml` values.

If you run your application in `test` mode:

```go title="internal/service/example.go"
// helpers
fmt.Printf("var: %s", s.config.GetEnvVar("APP_ENV")) // var: test

fmt.Printf("name: %s", s.config.AppName())           // name: app
fmt.Printf("env: %s", s.config.AppEnv())             // env: test
fmt.Printf("version: %s", s.config.AppVersion())     // version: 0.1.0
fmt.Printf("debug: %v", s.config.AppDebug())         // debug: true

// others
fmt.Printf("string_value: %s", s.config.GetString("config.values.string_value")) // string_value: test
fmt.Printf("int_value: %s", s.config.GetInt("config.values.int_value"))          // int_value: 0
```

You can use any value for `APP_ENV` (to allow you to reflect your own envs): for example if `APP_ENV=custom`, the module
will use `config.yaml` values and override them with `config.custom.yaml` values (you just need to ensure
that `config.custom.yaml` exists).

### Env var placeholders

This module offers the possibility to use placeholders in the config files to reference an env var value, that will be
resolved at runtime.

Placeholder pattern: `${ENV_VAR_NAME}`.

For example, with the env var `BAR=bar`:

```go title="internal/service/example.go"
// placeholder: foo-bar-baz
fmt.Printf("placeholder: %s", s.config.GetString("config.placeholder"))
```

### Env var substitution

This module offers the possibility to perform configuration files values substitution from env var values.

For example, if you have a configuration key `config.substitution=foo`, providing the env var `CONFIG_SUBSTITUTION=bar`
will override the value from `foo` to `bar`.

```go title="internal/service/example.go"
// substitution: bar
fmt.Printf("substitution: %s", cfg.GetString("config.substitution")) 
```
