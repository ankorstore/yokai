# gRPC Server Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/grpcserver-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/grpcserver-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/grpcserver)](https://goreportcard.com/report/github.com/ankorstore/yokai/grpcserver)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=grpcserver)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/grpcserver)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Fgrpcserver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/grpcserver)](https://pkg.go.dev/github.com/ankorstore/yokai/grpcserver)

> gRPC server module based on [gRPC-Go](https://github.com/grpc/grpc-go).

<!-- TOC -->

* [Installation](#installation)
* [Documentation](#documentation)
	* [Usage](#usage)
	* [Add-ons](#add-ons)
		* [Reflection](#reflection)
		* [Panic recovery](#panic-recovery)
		* [Logger interceptor](#logger-interceptor)
		* [Healthcheck service](#healthcheck-service)

<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/grpcserver
```

## Documentation

### Usage

This module provides a [GrpcServerFactory](factory.go), allowing to build an `grpc.Server` instance.

```go
package main

import (
	"github.com/ankorstore/yokai/grpcserver"
	"google.golang.org/grpc"
)

var server, _ = grpcserver.NewDefaultGrpcServerFactory().Create()

// equivalent to:
var server, _ = grpcserver.NewDefaultGrpcServerFactory().Create(
	grpcserver.WithServerOptions([]grpc.ServerOption{}), // no specific server options by default 
	grpcserver.WithReflection(false),                    // reflection disabled by default
)
```

See [gRPC-Go documentation](https://github.com/grpc/grpc-go) for more details.

### Add-ons

This module provides several add-ons ready to use to enrich your gRPC server.

#### Reflection

This module provides the possibility to easily
enable [gRPC server reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md), disabled by default:

```go
package main

import (
	"github.com/ankorstore/yokai/grpcserver"
)

func main() {
	server, _ := grpcserver.NewDefaultGrpcServerFactory().Create(
		grpcserver.WithReflection(true),
	)
}
```

Reflection usage is helpful for developing or testing your gRPC services, but it is not recommended for production
usage.

#### Panic recovery

This module provides an [GrpcPanicRecoveryHandler](panic.go), compatible with
the [recovery interceptor](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/main/interceptors/recovery), to
automatically recover from panics in your gRPC services:

```go
package main

import (
	"github.com/ankorstore/yokai/grpcserver"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

func main() {
	handler := grpcserver.NewGrpcPanicRecoveryHandler()

	server, _ := grpcserver.NewDefaultGrpcServerFactory().Create(
		grpcserver.WithServerOptions(
			grpc.UnaryInterceptor(recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(handler.Handle(false)))),
			grpc.StreamInterceptor(recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(handler.Handle(false)))),
		),
	)
}
```

You can also use `Handle(true)` to append on the handler gRPC response and logs more information about the panic and the debug stack (
not suitable for production).

#### Logger interceptor

This module provides a [GrpcLoggerInterceptor](logger.go) to automatically log unary and streaming RPCs calls (status,
type, duration, etc.), compatible with the [log module](https://github.com/ankorstore/yokai/tree/main/log).

Using this interceptor will also provide a logger instance in the context, that you can retrieve with
the [CtxLogger](context.go) method to produce correlated logs from your gRPC services.

```go
package main

import (
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/grpcserver"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/grpc"
)

func main() {
	logger, _ := log.NewDefaultLoggerFactory().Create()

	loggerInterceptor := grpcserver.NewGrpcLoggerInterceptor(uuid.NewDefaultUuidGenerator, logger)

	server, _ := grpcserver.NewDefaultGrpcServerFactory().Create(
		grpcserver.WithServerOptions(
			grpc.UnaryInterceptor(loggerInterceptor.UnaryInterceptor()),
			grpc.StreamInterceptor(loggerInterceptor.StreamInterceptor()),
		),
	)
}
```

The interceptor will automatically enrich each log records with the `x-request-id` fetch from the context metadata in
the field `requestID`.

You can specify additional metadata to add to logs records:

- the key is the metadata name to fetch
- the value is the log field to fill

```go
loggerInterceptor.Metadata(
    map[string]string{
        "x-some-metadata": "someMetadata",
        "x-other-metadata": "otherMetadata",
    },
)
```

You can also specify a list of gRPC methods to exclude from logging:

```go
loggerInterceptor.Exlude("/test.Service/Unary", "/test.Service/Bidi")
```

Note: even if excluded, failing gRPC methods calls will still be logged for observability purposes.

#### Healthcheck service

This module provides a [GrpcHealthCheckService](healthcheck.go), compatible with
the [healthcheck module](https://github.com/ankorstore/yokai/tree/main/healthcheck):

```go
package main

import (
	"probes"

	"github.com/ankorstore/yokai/grpcserver"
	"github.com/ankorstore/yokai/healthcheck"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	checker, _ := healthcheck.NewDefaultCheckerFactory().Create(
		healthcheck.WithProbe(probes.SomeProbe()),                            // register for startup, liveness and readiness
		healthcheck.WithProbe(probes.SomeOtherProbe(), healthcheck.Liveness), // register for liveness only
	)

	server, _ := grpcserver.NewDefaultGrpcServerFactory().Create()

	grpc_health_v1.RegisterHealthServer(server, grpcserver.NewGrpcHealthCheckService(checker))
}
```

This will expose the [Check and Watch](https://github.com/grpc/grpc-proto/blob/master/grpc/health/v1/health.proto) RPCs, suitable for [k8s startup, readiness or liveness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).

The checker will:

- run the `liveness` probes checks if the request service name contains `liveness` (like `kubernetes::liveness`)
- or run the `readiness` probes checks if the request service name contains `readiness` (like `kubernetes::readiness`)
- or run the `startup` probes checks otherwise
