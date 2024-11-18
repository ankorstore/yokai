# Fx Generate Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxgenerate-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxgenerate-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxgenerate)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxgenerate)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxgenerate)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxgenerate)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxgenerate)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxgenerate)](https://pkg.go.dev/github.com/ankorstore/yokai/fxgenerate)

> [Fx](https://uber-go.github.io/fx/) module for [generate](https://github.com/ankorstore/yokai/tree/main/generate).

<!-- TOC -->
* [Installation](#installation)
* [Documentation](#documentation)
  * [Loading](#loading)
  * [Generators](#generators)
    * [UUID V4](#uuid-v4)
      * [Usage](#usage)
      * [Testing](#testing)
    * [UUID V6](#uuid-v6)
      * [Usage](#usage-1)
      * [Testing](#testing-1)
    * [UUID V7](#uuid-v7)
      * [Usage](#usage-2)
      * [Testing](#testing-2)
  * [Override](#override)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxgenerate
```

## Documentation

### Loading

To load the module in your Fx application:

```go
package main

import (
	"github.com/ankorstore/yokai/fxgenerate"
	"go.uber.org/fx"
)

func main() {
	fx.New(fxgenerate.FxGenerateModule).Run()
}
```

### Generators

#### UUID V4

##### Usage

This module provides a [UuidGenerator](https://github.com/ankorstore/yokai/blob/main/generate/uuid/generator.go), made available into the Fx container.

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/fxgenerate"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxgenerate.FxGenerateModule,                     // load the module
		fx.Invoke(func(generator uuid.UuidGenerator) {   // invoke the uuid generator
			fmt.Printf("uuid: %s", generator.Generate()) // uuid: dcb5d8b3-4517-4957-a42c-604d11758561
		}),
	).Run()
}
```

##### Testing

This module provides the possibility to make your [UuidGenerator](https://github.com/ankorstore/yokai/blob/main/generate/uuid/generator.go) generate deterministic values, for testing purposes.

You need to:

- first provide into the Fx container the deterministic value to be used for generation, annotated with `name:"generate-test-uuid-value"`
- then decorate into the Fx container the `UuidGeneratorFactory` with the provided [TestUuidGeneratorFactory](fxgeneratetest/uuid/factory.go)

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/fxgenerate"
	fxtestuuid "github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuid"
	"github.com/ankorstore/yokai/generate/uuid"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxgenerate.FxGenerateModule, // load the module
		fx.Provide(                  // provide and annotate the deterministic value
			fx.Annotate(
				func() string {
					return "some deterministic value"
				},
				fx.ResultTags(`name:"generate-test-uuid-value"`),
			),
		),
		fx.Decorate(fxtestuuid.NewFxTestUuidGeneratorFactory), // override the module with the test factory
		fx.Invoke(func(generator uuid.UuidGenerator) {         // invoke the generator
			fmt.Printf("uuid: %s", generator.Generate())       // uuid: some deterministic value
		}),
	).Run()
}
```

#### UUID V6

##### Usage

This module provides a [UuidV6Generator](https://github.com/ankorstore/yokai/blob/main/generate/uuidv6/generator.go), made available into the Fx container.

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuidv6"
	"github.com/ankorstore/yokai/fxgenerate"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxgenerate.FxGenerateModule,                        // load the module
		fx.Invoke(func(generator uuidv6.UuidV7Generator) {
			uuid, _ := generator.Generate()                 // invoke the uuid v6 generator
			fmt.Printf("uuid: %s", uuid.String())           // uuid: 1efa5a47-e5d0-6667-9d00-49bf4f758c68
		}),
	).Run()
}
```

##### Testing

This module provides the possibility to make your [UuidV6Generator](https://github.com/ankorstore/yokai/blob/main/generate/uuidv6/generator.go) generate deterministic values, for testing purposes.

You need to:

- first provide into the Fx container the deterministic value to be used for generation, annotated with `name:"generate-test-uuid-v6-value"`
- then decorate into the Fx container the `UuidV6GeneratorFactory` with the provided [TestUuidGeneratorV6Factory](fxgeneratetest/uuidv6/factory.go)

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/fxgenerate"
	fxtestuuidv6 "github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuidv6"
	"github.com/ankorstore/yokai/generate/uuidv6"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxgenerate.FxGenerateModule, // load the module
		fx.Provide(                  // provide and annotate the deterministic value
			fx.Annotate(
				func() string {
					return "1efa5a47-e5d0-6663-99da-f6e8045dd166"
				},
				fx.ResultTags(`name:"generate-test-uuid-v6-value"`),
			),
		),
		fx.Decorate(fxtestuuidv6.NewFxTestUuidV6GeneratorFactory), // override the module with the test factory
		fx.Invoke(func(generator uuidv6.UuidV6Generator) {         // invoke the generator
			uuid, _ := generator.Generate()
			fmt.Printf("uuid: %s", uuid.String())                  // uuid: 1efa5a47-e5d0-6663-99da-f6e8045dd166
		}),
	).Run()
}
```

#### UUID V7

##### Usage

This module provides a [UuidV7Generator](https://github.com/ankorstore/yokai/blob/main/generate/uuidv7/generator.go), made available into the Fx container.

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/generate/uuidv7"
	"github.com/ankorstore/yokai/fxgenerate"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxgenerate.FxGenerateModule,                        // load the module
		fx.Invoke(func(generator uuidv7.UuidV7Generator) {
			uuid, _ := generator.Generate()                 // invoke the uuid v7 generator
			fmt.Printf("uuid: %s", uuid.String())           // uuid: 018fdd68-1b41-7eb0-afad-57f45297c7c1
		}),
	).Run()
}
```

##### Testing

This module provides the possibility to make your [UuidV7Generator](https://github.com/ankorstore/yokai/blob/main/generate/uuidv7/generator.go) generate deterministic values, for testing purposes.

You need to:

- first provide into the Fx container the deterministic value to be used for generation, annotated with `name:"generate-test-uuid-v7-value"`
- then decorate into the Fx container the `UuidV7GeneratorFactory` with the provided [TestUuidGeneratorV7Factory](fxgeneratetest/uuidv7/factory.go)

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/fxgenerate"
	fxtestuuidv7 "github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuidv7"
	"github.com/ankorstore/yokai/generate/uuidv7"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxgenerate.FxGenerateModule, // load the module
		fx.Provide(                  // provide and annotate the deterministic value
			fx.Annotate(
				func() string {
					return "018fdd68-1b41-7eb0-afad-57f45297c7c1"
				},
				fx.ResultTags(`name:"generate-test-uuid-v7-value"`),
			),
		),
		fx.Decorate(fxtestuuidv7.NewFxTestUuidV7GeneratorFactory), // override the module with the test factory
		fx.Invoke(func(generator uuidv7.UuidV7Generator) {         // invoke the generator
			uuid, _ := generator.Generate()
			fmt.Printf("uuid: %s", uuid.String())                  // uuid: 018fdd68-1b41-7eb0-afad-57f45297c7c1
		}),
	).Run()
}
```

### Override

If needed, you can provide your own factories and override the module:

```go
package main

import (
	"fmt"

	"github.com/ankorstore/yokai/fxgenerate"
	testuuid "github.com/ankorstore/yokai/fxgenerate/testdata/uuid"
	"github.com/ankorstore/yokai/generate/uuid"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxgenerate.FxGenerateModule,                             // load the module
		fx.Decorate(testuuid.NewTestStaticUuidGeneratorFactory), // override the module with a custom factory
		fx.Invoke(func(generator uuid.UuidGenerator) {          // invoke the custom generator
			fmt.Printf("uuid: %s", generator.Generate())         // uuid: static
		}),
	).Run()
}
```
