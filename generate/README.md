# Generate Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/generate-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/generate-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/generate)](https://goreportcard.com/report/github.com/ankorstore/yokai/generate)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=5s0g5WyseS&flag=generate)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/generate)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/generate)](https://pkg.go.dev/github.com/ankorstore/yokai/generate)

> Generation module based on [Google UUID](https://github.com/google/uuid).

<!-- TOC -->

* [Installation](#installation)
* [Documentation](#documentation)
	* [UUID](#uuid)

<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/generate
```

## Documentation

### UUID

This module provides an [UuidGenerator](uuid/generator.go) interface, allowing to generate UUIDs.

The `DefaultUuidGenerator` is based on [Google UUID](https://github.com/google/uuid).

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuid"
	uuidtest "github.com/ankorstore/yokai/generate/generatetest/uuid"
)

func main() {
	// default UUID generator
	generator := uuid.NewDefaultUuidGenerator()
	fmt.Printf("uuid: %s", generator.Generate()) // uuid: dcb5d8b3-4517-4957-a42c-604d11758561

	// test UUID generator (with deterministic value for testing)
	testGenerator := uuidtest.NewTestUuidGenerator("test")
	fmt.Printf("uuid: %s", testGenerator.Generate()) // uuid: test
}
```

The module also provides a [UuidGeneratorFactory](uuid/factory.go) interface, to create
the [UuidGenerator](uuid/generator.go) instances.

The `DefaultUuidGeneratorFactory` generates `DefaultUuidGenerator` instances.

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuid"
)

func main() {
	// default UUID generator factory
	generator := uuid.NewDefaultUuidGeneratorFactory().Create()
	fmt.Printf("uuid: %s", generator.Generate()) // uuid: dcb5d8b3-4517-4957-a42c-604d11758561
}
```