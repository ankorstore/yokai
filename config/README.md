# Config Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/config-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/config-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/config)](https://goreportcard.com/report/github.com/ankorstore/yokai/config)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=5s0g5WyseS&flag=config)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/config)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/config)](https://pkg.go.dev/github.com/ankorstore/yokai/config)

> Configuration module based on [Viper](https://github.com/spf13/viper).

<!-- TOC -->

* [Installation](#installation)
* [Documentation](#documentation)
	* [Configuration files](#configuration-files)
	* [Configuration usage](#configuration-usage)
		* [Configuration access](#configuration-access)
		* [Configuration dynamic env overrides](#configuration-dynamic-env-overrides)
		* [Configuration env var placeholders](#configuration-env-var-placeholders)
		* [Configuration env var substitution](#configuration-env-var-substitution)

<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/config
```

## Documentation

### Configuration files

By default, the module expects configuration files:

- to be present in `.` (root) or `./configs` directories of your project
- to be named `config.{format}` (ex: `config.yaml`, `config.json`, etc.)
- to offer env overrides files named `config.{env}.{format}` based on the env var `APP_ENV` (ex: `config.test.yaml` if
  env var `APP_ENV=test`)

Also:

- the config file name and lookup paths can be configured
- the following configuration files format are supported: JSON, TOML, YAML, HCL, INI, and env file.

### Configuration usage

For the following examples, we will be considering those configuration files:

```yaml
# ./configs/config.yaml
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

```yaml
# ./configs/config.test.yaml
app:
  env: test
  debug: true
config:
  values:
    string_value: test
```

and the following `Config` instance:

```go
package main

import "github.com/ankorstore/yokai/config"

var cfg, _ = config.NewDefaultConfigFactory().Create()

// equivalent to:
var cfg, _ = config.NewDefaultConfigFactory().Create(
	config.WithFileName("config"),          // config files base name
	config.WithFilePaths(".", "./configs"), // config files lookup paths
)
```

#### Configuration access

This module offers [helper methods](./config.go), as well as
all [Viper methods](https://github.com/spf13/viper/blob/master/viper.go) to access configuration values:

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/config"
)

func main() {
	// config
	cfg, _ := config.NewDefaultConfigFactory().Create()

	// helpers
	fmt.Printf("name: %s", cfg.AppName())       // name: app
	fmt.Printf("env: %s", cfg.AppEnv())         // env: dev
	fmt.Printf("version: %s", cfg.AppVersion()) // version: 0.1.0
	fmt.Printf("debug: %v", cfg.AppDebug())     // debug: false

	// others
	fmt.Printf("string_value: %s", cfg.GetString("config.values.string_value")) // string_value: default
	fmt.Printf("int_value: %s", cfg.GetInt("config.values.int_value"))          // int_value: 0
}
```

#### Configuration dynamic env overrides

This module offers the possibility to override dynamically (by merging) configuration files depending on the env
var `APP_ENV` value.

For example, if `APP_ENV=test`, the module will use `config.yaml` values and merge / override them
with `config.test.yaml` values.

```go
package main

import (
	"fmt"
	"os"

	"github.com/ankorstore/yokai/config"
)

func main() {
	// env vars
	os.Setenv("APP_ENV", "test")

	// config
	cfg, _ := config.NewDefaultConfigFactory().Create()

	// helpers
	fmt.Printf("name: %s", cfg.AppName())       // name: app
	fmt.Printf("env: %s", cfg.AppEnv())         // env: test
	fmt.Printf("version: %s", cfg.AppVersion()) // version: 0.1.0
	fmt.Printf("debug: %v", cfg.AppDebug())     // debug: true

	// others
	fmt.Printf("string_value: %s", cfg.GetString("config.values.string_value")) // string_value: test
	fmt.Printf("int_value: %s", cfg.GetInt("config.values.int_value"))          // int_value: 0
}
```

You can use any value for `APP_ENV` (to allow you to reflect your own envs): for example if `APP_ENV=custom`, the module
will use `config.yaml` values and override them with `config.custom.yaml` values (you just need to ensure
that `config.custom.yaml` exists).

#### Configuration env var placeholders

This module offers the possibility to use placeholders in the config files to reference an env var value, that will be
resolved at runtime.

Placeholder pattern: `${ENV_VAR_NAME}`.

```go
package main

import (
	"fmt"
	"os"

	"github.com/ankorstore/yokai/config"
)

func main() {
	// env vars
	os.Setenv("BAR", "bar")

	// config
	cfg, _ := config.NewDefaultConfigFactory().Create()

	// env var placeholder value
	fmt.Printf("placeholder: %s", cfg.GetString("config.placeholder")) // placeholder: foo-bar-baz
}
```

#### Configuration env var substitution

This module offers the possibility to perform configuration files values substitution from env var values.

For example, if you have a configuration key `config.substitution=foo`, providing the `CONFIG_SUBSTITUTION=bar` env var
will override the value from `foo` to `bar`.

```go
package main

import (
	"fmt"
	"os"

	"github.com/ankorstore/yokai/config"
)

func main() {
	// env vars
	os.Setenv("CONFIG_SUBSTITUTION", "bar")

	// config
	cfg, _ := config.NewDefaultConfigFactory().Create()

	// env var substitution value
	fmt.Printf("substitution: %s", cfg.GetString("config.substitution")) // substitution: bar
}
```
