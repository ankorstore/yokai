# Generate Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxgenerate-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxgenerate-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxgenerate)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxgenerate)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxgenerate)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxgenerate)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxgenerate)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxgenerate)](https://pkg.go.dev/github.com/ankorstore/yokai/fxgenerate)

## Overview

Yokai provides a [fxgenerate](https://github.com/ankorstore/yokai/tree/main/fxgenerate) module, allowing your application to generate UUIDs.

It wraps the [generate](https://github.com/ankorstore/yokai/tree/main/generate) module, based on [Google UUID](https://github.com/google/uuid).

## Installation

The [fxgenerate](https://github.com/ankorstore/yokai/tree/main/fxgenerate) module is automatically loaded by
the [fxcore](https://github.com/ankorstore/yokai/tree/main/fxcore).

When you use a Yokai [application template](https://ankorstore.github.io/yokai/applications/templates/), you have nothing to install, it's ready to use.

## Usage

This module makes available the [UuidGenerator](https://github.com/ankorstore/yokai/blob/main/generate/uuid/generator.go) in
Yokai dependency injection system.

It is built on top of `Google UUID`, see its [documentation](https://github.com/google/uuid) for more details about available methods.

To access it, you just need to inject it where needed, for example:

```go title="internal/service/example.go"
package service

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuid"
)

type ExampleService struct {
	generator uuid.UuidGenerator
}

func NewExampleService(generator uuid.UuidGenerator) *ExampleService {
	return &ExampleService{
		generator: generator,
	}
}

func (s *ExampleService) DoSomething() {
	// uuid: dcb5d8b3-4517-4957-a42c-604d11758561
	fmt.Printf("uuid: %s", s.generator.Generate())
}
```

## Testing

This module provides the possibility to make the [UuidGenerator](https://github.com/ankorstore/yokai/blob/main/generate/uuid/generator.go) generate deterministic values (for testing purposes).

For this, you need to:

- first provide the deterministic value to be used for generation, annotated with `name:"generate-test-uuid-value"`
- then override the `UuidGeneratorFactory` with the provided [TestUuidGeneratorFactory](https://github.com/ankorstore/yokai/blob/main/fxgenerate/fxgeneratetest/uuid/factory.go)

```go title="internal/service/example_test.go"
package internal_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuid"
	"github.com/foo/bar/internal"
	"github.com/foo/bar/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestExampleServiceDoSomething(t *testing.T) {
	var exampleService *service.ExampleService

	internal.RunTest(
		t,
		fx.Populate(&exampleService),
		// provide and annotate the deterministic value
		fx.Provide(
			fx.Annotate(
				func() string {
					return "some value"
				},
				fx.ResultTags(`name:"generate-test-uuid-value"`),
			),
		),
		// override the UuidGeneratorFactory with the TestUuidGeneratorFactory
		fx.Decorate(uuid.NewFxTestUuidGeneratorFactory),
	)
	
	//assertion
	assert.Equal(t, "uuid: some value", exampleService.DoSomething())
}
```