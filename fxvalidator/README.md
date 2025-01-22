# Fx Validator Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxvalidator-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxvalidator-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxvalidator)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxvalidator)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxvalidator)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxvalidator)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxvalidator)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxvalidator)](https://pkg.go.dev/github.com/ankorstore/yokai/fxvalidator)

> [Fx](https://uber-go.github.io/fx/) module for [go-playground/validator](https://github.com/go-playground/validator).

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
  * [Dependencies](#dependencies)
  * [Loading](#loading)
  * [Configuration](#configuration)
  * [Customization](#customization)
    * [Custom aliases](#custom-aliases)
    * [Custom validations](#custom-validations)
    * [Custom struct validations](#custom-struct-validations)
    * [Custom types](#custom-types)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxvalidator
```

## Features

This module provides a [*validator.Validate](https://github.com/go-playground/validator) to your Fx application, that:

- you can inject anywhere
- you can customize depending on your needs

## Documentation

### Dependencies

This module is intended to be used alongside the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module.

### Loading

To load the module in your Fx application:

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,                       // load the module dependency
		fxvalidator.FXValidatorModule,                 // load the module
		fx.Invoke(func(validate *validator.Validate) { // invoke the validator
			err := validate.Var("foo", "alpha")
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
		}),
	).Run()
}
```

### Configuration

Configuration reference:

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: false
modules:
  validator:
    tag_name: validate    # struct tag to define validation rules, default = validate
    private_fields: false # to enable validation on private fields, disabled by default
```


### Customization

#### Custom aliases

This module provides the possibility to register custom aliases with `AsAlias()`:

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,                                 // load the module dependency
		fxvalidator.FXValidatorModule,                           // load the module
		fxvalidator.AsAlias("foobar", "required,oneof=foo bar"), // register the foobar alias: required and either foo or bar
		fx.Invoke(func(validate *validator.Validate) {           // invoke the validator
			// validation success
			err := validate.Var("foo", "foobar")
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
			
			// validation error
			err = validate.Var("invalid", "foobar")
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
		}),
	).Run()
}
```

#### Custom validations

This module provides the possibility to register custom validations functions with `AsValidation()`:

```go
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var fn = func(ctx context.Context, fl validator.FieldLevel) bool {
	s := strings.ToLower(fl.Field().String())
	
	return s == "foo" || s == "bar"
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,                         // load the module dependency
		fxvalidator.FXValidatorModule,                   // load the module
		fxvalidator.AsValidation("foobar-ci", fn, true), // register the foobar-ci validation: either foo or bar (case insensitive)
		fx.Invoke(func(validate *validator.Validate) {   // invoke the validator
			// validation success
			err := validate.Var("Foo", "foobar-ci")
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
			
			// validation success
			err = validate.Var("baR", "foobar-ci")
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
			
			// validation error
			err = validate.Var("invalid", "foobar-ci")
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
		}),
	).Run()
}
```

#### Custom struct validations

This module provides the possibility to register custom struct validations functions with `AsStructValidation()`:

```go
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

type FooBar struct {
	Foo string
	Bar string
}

var fn = func(ctx context.Context, sl validator.StructLevel) {
	fb, ok := sl.Current().Interface().(FooBar)
	if ok {
		if fb.Foo != "foo" {
			sl.ReportError(fb.Foo, "Foo", "Foo", "invalid-foo", "invalid foo")
		}
		if fb.Bar != "bar" {
			sl.ReportError(fb.Bar, "Bar", "Bar", "invalid-bar", "invalid bar")
		}
	}
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,                         // load the module dependency
		fxvalidator.FXValidatorModule,                   // load the module
		fxvalidator.AsStructValidation(fn, FooBar{}),    // register a struct validation for FooBar{}
		fx.Invoke(func(validate *validator.Validate) {   // invoke the validator
			// validation success
			err := validate.Struct(FooBar{Foo: "foo", Bar: "bar"})
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
			
			// validation error
			err = validate.Struct(FooBar{Foo: "invalid", Bar: "bar"})
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
		}),
	).Run()
}
```

#### Custom types

This module provides the possibility to register custom type functions with `AsCustomType()`:

```go
package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

type Foo struct {
	Value string `validate:"required"`
}

type Bar struct {
	Value string `validate:"required"`
}

type FooBar struct {
	Foo Foo
	Bar Bar
}

var fn = func(field reflect.Value) interface{} {
	if f, ok := field.Interface().(Foo); ok {
		if f.Value != "foo" {
			return ""
		}

		return f.Value
	}

	if b, ok := field.Interface().(Bar); ok {
		if b.Value != "bar" {
			return ""
		}

		return b.Value
	}

	return ""
}

func main() {
	fx.New(
		fxconfig.FxConfigModule,                       // load the module dependency
		fxvalidator.FXValidatorModule,                 // load the module
		fxvalidator.AsCustomType(fn, Foo{}, Bar{}),    // register a custom types Foo{} and Bar{}
		fx.Invoke(func(validate *validator.Validate) { // invoke the validator
			// validation success
			err := validate.Struct(FooBar{Foo: Foo{Value: "foo"}, Bar: Bar{Value: "bar"}})
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}

			// validation error
			err = validate.Struct(FooBar{Foo: Foo{Value: "invalid"}, Bar: Bar{Value: ""}})
			if valErrs, ok := err.(*validator.ValidationErrors); ok {
				fmt.Println(valErrs)
			}
		}),
	).Run()
}
```