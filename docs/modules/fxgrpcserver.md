---
title: Modules - gRPC Server
icon: material/cube-outline
---

# :material-cube-outline: gRPC Server Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxgrpcserver-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxgrpcserver-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxgrpcserver)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxgrpcserver)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxgrpcserver)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxgrpcserver)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxgrpcserver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxgrpcserver)](https://pkg.go.dev/github.com/ankorstore/yokai/fxgrpcserver)

## Overview

Yokai provides a [fxgrpcserver](https://github.com/ankorstore/yokai/tree/main/fxgrpcserver) module, offering an [gRPC](https://grpc.io/) server to your application.

It wraps the [grpcserver](https://github.com/ankorstore/yokai/tree/main/grpcserver) module, based on [gRPC-Go](https://github.com/grpc/grpc-go).

It comes with:

- automatic panic recovery
- automatic logging and tracing (method, duration, status, ...)
- automatic metrics
- automatic healthcheck
- automatic reflection
- possibility to register gRPC server options, interceptors and services

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxgrpcserver
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxgrpcserver"
)

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxgrpcserver.FxGrpcServerModule,
	// ...
)
```

## Configuration

```yaml title="configs/config.yaml"
modules:
  grpc:
    server:
      address: ":50051"             # gRPC server listener address (default :50051)
      log:
        metadata:                   # list of gRPC metadata to add to logs on top of x-request-id, empty by default
          x-foo: foo                # to log for example the metadata x-foo in the log field foo
          x-bar: bar
        exclude:                    # list of gRPC methods to exclude from logging, empty by default
          - /test.Service/Unary
      trace:
        enabled: true               # to trace gRPC calls, disabled by default
        exclude:                    # list of gRPC methods to exclude from tracing, empty by default
          - /test.Service/Bidi
      metrics:
        collect:
          enabled: true             # to collect gRPC server metrics, disabled by default
          namespace: foo            # gRPC server metrics namespace (empty by default)
          subsystem: bar            # gRPC server metrics subsystem (empty by default)
        buckets: 0.1, 1, 10         # to override default request duration buckets (default prometheus.DefBuckets)
      reflection:
        enabled: true               # to expose gRPC reflection service, disabled by default
      healthcheck:
        enabled: true               # to expose gRPC healthcheck service, disabled by default
      test:
        bufconn:
          size: 1048576             # test gRPC bufconn size, 1024*1024 by default
```

## Usage

This module offers the possibility to easily register gRPC server options, interceptors and services.

### Server options registration

You can use the `AsGrpcServerOptions()` function to register [grpc.ServerOption](https://pkg.go.dev/google.golang.org/grpc#ServerOption) on your [gRPC server](https://pkg.go.dev/google.golang.org/grpc#Server).

For example:

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxgrpcserver"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func Register() fx.Option {
	return fx.Options(
		// configure the server send and receive max message size
		fxgrpcserver.AsGrpcServerOptions(
			grpc.MaxSendMsgSize(1000),
			grpc.MaxRecvMsgSize(1000),
		),
		// ...
	)
}
```

### Server interceptors registration

You can create [gRPC server interceptors](https://github.com/grpc/grpc-go/blob/master/examples/features/interceptor/README.md#server-side) for your [gRPC server](https://pkg.go.dev/google.golang.org/grpc#Server).

You need to implement:

- the [GrpcServerUnaryInterceptor](https://github.com/ankorstore/yokai/blob/main/fxgrpcserver/registry.go) interface for `unary` interceptors
- the [GrpcServerStreamInterceptor](https://github.com/ankorstore/yokai/blob/main/fxgrpcserver/registry.go) interface for `stream` interceptors

Example of `unary` interceptor:

```go title="internal/interceptor/unary.go"
package interceptor

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/service"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/grpc"
)

type UnaryInterceptor struct {
	config *config.Config
}

func NewUnaryInterceptor(cfg *config.Config) *UnaryInterceptor {
	return &UnaryInterceptor{
		config: cfg,
	}
}

func (i *UnaryInterceptor) HandleUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.CtxLogger(ctx).Info().Msgf("in unary interceptor of %s", i.config.AppName())

		return handler(ctx, req)
	}
}
```

Example of `stream` interceptor:

```go title="internal/interceptor/stream.go"
package interceptor

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/service"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/grpc"
)

type StreamInterceptor struct {
	config *config.Config
}

func NewStreamInterceptor(cfg *config.Config) *StreamInterceptor {
	return &StreamInterceptor{
		config: cfg,
	}
}

func (i *StreamInterceptor) HandleStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.CtxLogger(ss.Context()).Info().Msgf("in stream interceptor of %s", i.config.AppName())

		return handler(srv, ss)
	}
}
```

You can register your interceptors:

- with `AsGrpcServerUnaryInterceptor()` to register a `unary` interceptor
- with `AsGrpcServerStreamInterceptor()` to register a `stream` interceptor

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/foo/bar/internal/interceptor"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func Register() fx.Option {
	return fx.Options(
		// registers UnaryInterceptor as server unary interceptor
		fxgrpcserver.AsGrpcServerUnaryInterceptor(interceptor.NewUnaryInterceptor),
		// registers StreamInterceptor as server stream interceptor
		fxgrpcserver.AsGrpcServerStreamInterceptor(interceptor.NewStreamInterceptor),
		// ...
	)
}
```

The dependencies of your interceptors will be autowired.

### Server services registration

You can use the `AsGrpcServerService()` function to register your gRPC server services and their definitions.

For example, with the [TestService](https://github.com/ankorstore/yokai/blob/main/fxgrpcserver/testdata/service/service.go), server implementation for the [test.proto](https://github.com/ankorstore/yokai/blob/main/fxgrpcserver/testdata/proto/test.proto):

```go title="internal/register.go"
package internal

import (
	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/proto"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/service"
	"github.com/foo/bar/internal/interceptor"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func Register() fx.Option {
	return fx.Options(
		// register the TestServiceServer for the proto.Service_ServiceDesc
		fxgrpcserver.AsGrpcServerService(service.NewTestServiceServer, &proto.Service_ServiceDesc),
		// ...
	)
}
```

The dependencies of your services will be autowired.

## Reflection

This module provides the possibility to enable [gRPC server reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md) with `modules.grpc.server.reflection.enabled=true`.

```yaml title="configs/config.yaml"
modules:
  grpc:
    server:
      reflection:
        enabled: true # to expose gRPC reflection service, disabled by default
```

Reflection usage is helpful for developing or testing your gRPC services, but it is NOT recommended for production usage (disabled by default).

## Health Check

This module automatically expose the [GrpcHealthCheckService](https://github.com/ankorstore/yokai/blob/main/grpcserver/healthcheck.go) with `modules.grpc.server.healthcheck.enabled=true`, to offer the [Check and Watch RPCs](https://github.com/grpc/grpc-proto/blob/master/grpc/health/v1/health.proto), suitable for [k8s gRPC startup, readiness or liveness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).

```yaml title="configs/config.yaml"
modules:
  grpc:
    server:
      healthcheck:
        enabled: true # to expose gRPC healthcheck service, disabled by default
```

You can use the `fxhealthcheck.AsCheckerProbe()` function to register several CheckerProbe (more details on the [fxhealthcheck module documentation](fxhealthcheck.md#probes-registration)).

The [GrpcHealthCheckService](https://github.com/ankorstore/yokai/blob/main/grpcserver/healthcheck.go) will:

- run the `liveness` probes checks if the request service name contains `liveness` (like `kubernetes::liveness`)
- or run the `readiness` probes checks if the request service name contains `readiness` (like `kubernetes::readiness`)
- or run the `startup` probes checks otherwise

## Logging

You can configure RPC calls automatic logging:

```yaml title="configs/config.yaml"
modules:
  grpc:
    server:
      log:
        metadata:      # list of gRPC metadata to add to logs on top of x-request-id, empty by default
          x-foo: foo   # to log for example the metadata x-foo in the log field foo
          x-bar: bar
        exclude:       # list of gRPC methods to exclude from logging, empty by default
          - /test.Service/ToExclude
```

As a result, in your application logs:

```
DBG grpc call start grpcMethod=/test.Service/Unary grpcType=unary requestID=77480bd0-6d7e-42ba-bf60-9a168b9d416f service=app spanID=129a13d8f496481b system=grpcserver traceID=b016d12bdef779d793f314d476aa271f
INF grpc call success grpcCode=0 grpcDuration="126.745µs" grpcMethod=/test.Service/Unary grpcStatus=OK grpcType=unary requestID=77480bd0-6d7e-42ba-bf60-9a168b9d416f service=app spanID=129a13d8f496481b system=grpcserver traceID=b016d12bdef779d793f314d476aa271f
```

If both gRPC server logging and tracing are enabled, log records will automatically have the current `traceID` and `spanID` to be able to correlate logs and trace spans.

If a request to an excluded gRPC method fails, the gRPC server will still log for observability purposes.

To get logs correlation in your gRPC server services, you need to retrieve the logger from the context with `log.CtxLogger()`:

```go
log.CtxLogger(ctx).Info().Msg("example message")
```

You can also use the function `grpcserver.CtxLogger()`:

```go
grpcserver.CtxLogger(ctx).Info().Msg("example message")
```

The gRPC server logging will be based on the [log](fxlog.md) module configuration.

## Tracing

You can enable RPC calls automatic tracing with `modules.grpc.server.trace.enable=true`:

```yaml title="configs/config.yaml"
modules:
  grpc:
    server:
      trace:
        enabled: true   # to trace gRPC calls, disabled by default
        exclude:        # list of gRPC methods to exclude from tracing, empty by default
          - /test.Service/ToExclude
```

As a result, in your application trace spans attributes:

```
rpc.service: test.Service
rpc.method: Unary
rpc.grpc.status_code: 0
...
```

To get traces correlation in your grpc server services, you need to retrieve the tracer provider from the context with `trace.CtxTracerProvider()`:

```go
ctx, span := trace.CtxTracerProvider(ctx).Tracer("example tracer").Start(ctx, "example span")
defer span.End()
```

You can also use the shortcut function `grpcserver.CtxTracer()`:

```go
ctx, span := grpcserver.CtxTracer(ctx).Start(ctx, "example span")
defer span.End()
```

The gRPC server tracing will be based on the [trace](fxtrace.md) module configuration.

## Metrics

You can enable RPC calls automatic metrics with `modules.grpc.server.metrics.collect.enable=true`:

```yaml title="configs/config.yaml"
modules:
  grpc:
    server:
      metrics:
        collect:
          enabled: true          # to collect gRPC server metrics, disabled by default
          namespace: foo         # gRPC server metrics namespace (empty by default)
          subsystem: bar         # gRPC server metrics subsystem (empty by default)
        buckets: 0.1, 1, 10      # to override default request duration buckets (default prometheus.DefBuckets)
```

For example, after calling `/test.Service/Unary`, the [core](fxcore.md) HTTP server will expose in the configured metrics endpoint:

```makefile title="[GET] /metrics"
# ...
# HELP grpc_server_started_total Total number of RPCs started on the server.
# TYPE grpc_server_started_total counter
grpc_server_started_total{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary"} 1
# HELP grpc_server_handled_total Total number of RPCs completed on the server, regardless of success or failure.
# TYPE grpc_server_handled_total counter
rpc_server_handled_total{grpc_code="OK",grpc_method="Unary",grpc_service="test.Service",grpc_type="unary"} 1
# HELP rpc_server_msg_received_total Total number of RPC stream messages received on the server.
# TYPE rpc_server_msg_received_total counter
grpc_server_msg_received_total{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary"} 1
# HELP grpc_server_msg_sent_total Total number of gRPC stream messages sent by the server.
# TYPE grpc_server_msg_sent_total counter
grpc_server_msg_sent_total{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary"} 1
# HELP grpc_server_handling_seconds Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.
# TYPE grpc_server_handling_seconds histogram
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="0.005"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="0.01"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="0.025"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="0.05"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="0.1"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="0.25"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="0.5"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="1"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="2.5"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="5"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="10"} 1
grpc_server_handling_seconds_bucket{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary",le="+Inf"} 1
grpc_server_handling_seconds_sum{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary"} 0.000103358
grpc_server_handling_seconds_count{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary"} 1
```

## Testing

This module provides a `*bufconn.Listener` that will automatically be used by the gRPC server in `test` mode.

You can then use this listener with your gRPC clients to provide `functional` tests for your gRPC services.

You can find tests examples in the [gRPC server module tests](https://github.com/ankorstore/yokai/blob/main/fxgrpcserver/module_test.go).