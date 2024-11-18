# Generate Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/generate-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/generate-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/generate)](https://goreportcard.com/report/github.com/ankorstore/yokai/generate)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=generate)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/generate)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Fgenerate)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/generate)](https://pkg.go.dev/github.com/ankorstore/yokai/generate)

> Generation module based on [Google UUID](https://github.com/google/uuid).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
  * [UUID V4](#uuid-v4)
  * [UUID V6](#uuid-v6)
  * [UUID V7](#uuid-v7)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/generate
```

## Documentation

### UUID V4

This module provides an [UuidGenerator](uuid/generator.go) interface, allowing to generate UUIDs V4.

The `DefaultUuidGenerator` implementing it is based on [Google UUID](https://github.com/google/uuid).

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

### UUID V6

This module provides an [UuidV6Generator](uuidv6/generator.go) interface, allowing to generate UUIDs V6.

The `DefaultUuidV6Generator` implementing it  is based on [Google UUID](https://github.com/google/uuid).

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuidv6"
	uuidv6test "github.com/ankorstore/yokai/generate/generatetest/uuidv6"
)

func main() {
	// default UUID V6 generator
	generator := uuidv6.NewDefaultUuidV6Generator()
	uuid, _ := generator.Generate()
	fmt.Printf("uuid: %s", uuid.String()) // uuid: 1efa5a1e-a679-67b7-ae79-f36f6749aa6b

	// test UUID generator (with deterministic value for testing, requires valid UUID v6)
	testGenerator, _ := uuidv6test.NewTestUuidV6Generator("1efa5a08-2883-6652-b357-5dd221ce0561")
    uuid, _ = testGenerator.Generate()
	fmt.Printf("uuid: %s", uuid.String()) // uuid: 1efa5a08-2883-6652-b357-5dd221ce0561
}
```

The module also provides a [UuidV6GeneratorFactory](uuidv6/factory.go) interface, to create
the [UuidV6Generator](uuidv6/generator.go) instances.

The `DefaultUuidV6GeneratorFactory` generates `DefaultUuidV6Generator` instances.

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuidv6"
)

func main() {
	// default UUID generator factory
	generator := uuidv6.NewDefaultUuidV6GeneratorFactory().Create()
	uuid, _ := generator.Generate()
	fmt.Printf("uuid: %s", uuid.String()) // uuid: 1efa5a1e-a679-67b7-ae79-f36f6749aa6b
}
```

### UUID V7

This module provides an [UuidV7Generator](uuidv7/generator.go) interface, allowing to generate UUIDs V7.

The `DefaultUuidV7Generator` implementing it  is based on [Google UUID](https://github.com/google/uuid).

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuidv7"
	uuidv7test "github.com/ankorstore/yokai/generate/generatetest/uuidv7"
)

func main() {
	// default UUID V7 generator
	generator := uuidv7.NewDefaultUuidV7Generator()
	uuid, _ := generator.Generate()
	fmt.Printf("uuid: %s", uuid.String()) // uuid: 018fdd68-1b41-7eb0-afad-57f45297c7c1

	// test UUID generator (with deterministic value for testing, requires valid UUID v7)
	testGenerator, _ := uuidv7test.NewTestUuidV7Generator("018fdd7d-1576-7a21-900e-1399637bd1a1")
    uuid, _ = testGenerator.Generate()
	fmt.Printf("uuid: %s", uuid.String()) // uuid: 018fdd7d-1576-7a21-900e-1399637bd1a1
}
```

The module also provides a [UuidV7GeneratorFactory](uuidv7/factory.go) interface, to create
the [UuidV7Generator](uuidv7/generator.go) instances.

The `DefaultUuidV7GeneratorFactory` generates `DefaultUuidV7Generator` instances.

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuidv7"
)

func main() {
	// default UUID generator factory
	generator := uuidv7.NewDefaultUuidV7GeneratorFactory().Create()
	uuid, _ := generator.Generate()
	fmt.Printf("uuid: %s", uuid.String()) // uuid: 018fdd68-1b41-7eb0-afad-57f45297c7c1
}
```