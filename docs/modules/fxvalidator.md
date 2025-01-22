---
title: Modules - Validator
icon: material/cube-outline
---

# :material-cube-outline: Validator Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxvalidator-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxvalidator-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxvalidator)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxvalidator)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxvalidator)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxvalidator)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxvalidator)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxvalidator)](https://pkg.go.dev/github.com/ankorstore/yokai/fxvalidator)
## Overview

Yokai provides a [fxvalidator](https://github.com/ankorstore/yokai/tree/main/fxvalidator) module, allowing you to inject a validator anywhere needed.

It wraps the [go-playground/validator](https://github.com/go-playground/validator) module.

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxvalidator
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxvalidator"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxvalidator.FXValidatorModule,
	// ...
)
```

## Configuration

Configuration reference:

```yaml title="configs/config.yaml"
modules:
  validator:
    tag_name: validate    # struct tag to define validation rules, default = validate
    private_fields: false # to enable validation on private fields, disabled by default
```

## Usage

This module makes available a `*validator.Validate` instance in
Yokai dependency injection system, that you can inject anywhere.

For example:

```go title="internal/service/example.go"
package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ExampleStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

type ExampleService struct {
	validate *validator.Validate
}

func NewExampleService(validate *validator.Validate) *ExampleService {
	return &ExampleService{
		validate: validate,
	}
}

func (s *ExampleService) DoSomething() error {
	es := ExampleStruct{
		Name:  "name",
		Email: "name@example.com",
	}
	
	err := s.validate.Struct(es)
	if valErrs, ok := err.(*validator.ValidationErrors); ok {
		fmt.Println(valErrs)
	}
	
	return valErrs
}
```

See [go-playground/validator](https://github.com/go-playground/validator) documentation for more details about available validation features.

## Customization

This module provides the possibility to easily customize your validator.

### Custom aliases

You can register custom validation aliases with `AsAlias()`:

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxvalidator"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the rla alias: provided value must be required, alpha only and lowercase
		fxvalidator.AsAlias("rla", "required,alpha,lowercase"),
		// ...
	)
}
```

Then:

```go title="internal/service/example.go"
package service

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type ExampleService struct {
	validate *validator.Validate
}

func NewExampleService(validate *validator.Validate) *ExampleService {
	return &ExampleService{
		validate: validate,
	}
}

func (s *ExampleService) DoSomething(ctx context.Context) error {
	// valid
	err := s.validate.VarCtx(ctx, "valid", "rla")

	// invalid
	err = s.validate.VarCtx(ctx, "1234", "rla")

	return err
}
```

### Custom validations

You can register custom validations functions with `AsValidation()`:

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxvalidator"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register the foobar validation: provided lowercased value must be "foo" or "bar"
		fxvalidator.AsValidation(
			"foobar",
			func(ctx context.Context, fl validator.FieldLevel) bool {
				s := strings.ToLower(fl.Field().String())
	
				return s == "foo" || s == "bar"
			},
			true,
		),
		// ...
	)
}
```

Then:

```go title="internal/service/example.go"
package service

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type ExampleService struct {
	validate *validator.Validate
}

func NewExampleService(validate *validator.Validate) *ExampleService {
	return &ExampleService{
		validate: validate,
	}
}

func (s *ExampleService) DoSomething(ctx context.Context) error {
	// valid
	err := s.validate.VarCtx(ctx, "FoO", "foobar")

	// invalid
	err = s.validate.VarCtx(ctx, "invalid", "foobar")

	return err
}
```

### Custom struct validations

You can register custom struct validations functions with `AsStructValidation()`:

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/foo/bar/internal/service"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register a struct validation for the struct service.ExampleStruct{}
		fxvalidator.AsStructValidation(
			func(ctx context.Context, sl validator.StructLevel) {
				fb, ok := sl.Current().Interface().(service.ExampleStruct)
				if ok {
					if fb.Foo != "foo" {
						sl.ReportError(fb.Foo, "Foo", "Foo", "invalid-foo", "invalid foo")
					}
					if fb.Bar != "bar" {
						sl.ReportError(fb.Bar, "Bar", "Bar", "invalid-bar", "invalid bar")
					}
				}
			},
			service.ExampleStruct{},
		),
		// ...
	)
}
```

Then:

```go title="internal/service/example.go"
package service

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type ExampleStruct struct {
	Foo string
	Bar string
}

type ExampleService struct {
	validate *validator.Validate
}

func NewExampleService(validate *validator.Validate) *ExampleService {
	return &ExampleService{
		validate: validate,
	}
}

func (s *ExampleService) DoSomething(ctx context.Context) error {
	// valid
	err := s.validate.StructCtx(ctx, ExampleStruct{Foo: "foo", Bar: "bar"})

	// invalid
	err = s.validate.StructCtx(ctx, ExampleStruct{Foo: "invalid", Bar: ""})

	return err
}
```

### Custom types

You can register custom types functions with `AsCustomType()`:

```go title="internal/register.go"
package internal

import (
	"reflect"
	
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/foo/bar/internal/service"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		// register a custom type function for service.Foo{} and service.Bar{}
		fxvalidator.AsCustomType(
			func(field reflect.Value) interface{} {
				if f, ok := field.Interface().(service.Foo); ok {
					if f.Value != "foo" {
						return ""
					}
	
					return f.Value
				}
	
				if b, ok := field.Interface().(service.Bar); ok {
					if b.Value != "bar" {
						return ""
					}
	
					return b.Value
				}
	
				return ""
			},
			service.Foo{},
			service.Bar{},
		),
		// ...
	)
}
```

Then:

```go title="internal/service/example.go"
package service

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type Foo struct {
	Value string `validate:"required"`
}

type Bar struct {
	Value string `validate:"required"`
}

type ExampleStruct struct {
	Foo Foo
	Bar Bar
}

type ExampleService struct {
	validate *validator.Validate
}

func NewExampleService(validate *validator.Validate) *ExampleService {
	return &ExampleService{
		validate: validate,
	}
}

func (s *ExampleService) DoSomething(ctx context.Context) error {
	// valid
	err := s.validate.StructCtx(ctx, ExampleStruct{Foo: Foo{Value: "foo"}, Bar: Bar{Value: "bar"}})

	// invalid
	err = s.validate.StructCtx(ctx, ExampleStruct{Foo: Foo{Value: "invalid"}, Bar: Bar{Value: ""}})

	return err
}
```
