---
icon: material/folder-eye-outline
---

# :material-folder-eye-outline: Demo - gRPC application

> Yokai provides a [gRPC demo application](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo).

## Overview

This [gRPC demo application](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo) is a simple gRPC API offering a [text transformation service](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/proto/transform.proto).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [gRPC server](../modules/fxgrpcserver.md) module to offer the gRPC API
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

### Layout

This demo application is following the [standard go project layout](https://github.com/golang-standards/project-layout):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `interceptor/`: gRPC interceptors
	- `service/`: gRPC services
	- `bootstrap.go`: bootstrap (modules, lifecycles, etc)
	- `services.go`: dependency injection
- `proto/`: protobuf definition and stubs

### Makefile

This demo application provides a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make stubs  # generate gRPC stubs with protoc
make test   # run tests
make lint   # run linter
```

## Usage

### Start the application

To start the application, simply run:

```shell
make fresh
```

After a short moment, the application will offer:

- `localhost:50051`: application gRPC server
- [http://localhost:8081](http://localhost:8081): application core dashboard
- [http://localhost:16686](http://localhost:16686): jaeger UI

### Available services

This demo application provides a [TransformTextService](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/proto/transform.proto), with the following `RPCs`:

| RPC                     | Type      | Description                                                  |
|-------------------------|-----------|--------------------------------------------------------------|
| `TransformText`         | unary     | Transforms a given text using a given transformer            |
| `TransformAndSplitText` | streaming | Transforms and splits a given text using a given transformer |

This demo application also provides [reflection](../modules/fxgrpcserver.md#reflection) and [health check ](../modules/fxgrpcserver.md#health-check) services.

If you update the [proto definition](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/proto/transform.proto), you can run `make stubs` to regenerate the stubs.

### Authentication

This demo application provides example [authentication interceptors](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/internal/interceptor/authentication.go).

You can enable authentication in the application [configuration file](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/configs/config.yaml) with `config.authentication.enabled=true`.

If enabled, you need to provide the secret configured in `config.authentication.secret` as context `authorization` metadata.

### Examples

Usage examples with [grpcurl](https://github.com/fullstorydev/grpcurl):

- with `TransformTextService/TransformText`:

```shell
grpcurl -plaintext -d '{"text":"abc","transformer":"TRANSFORMER_UPPERCASE"}' localhost:50051 transform.TransformTextService/TransformText
{
  "text": "ABC"
}
```

- with `TransformTextService/TransformAndSplitText`:

```shell
grpcurl -plaintext -d '{"text":"ABC DEF","transformer":"TRANSFORMER_LOWERCASE"}' localhost:50051 transform.TransformTextService/TransformAndSplitText
{
  "text": "abc"
}
{
  "text": "def"
}
```

You can use any gRPC clients, for example [Postman](https://learning.postman.com/docs/sending-requests/grpc/grpc-request-interface/) or [Evans](https://github.com/ktr0731/evans).
