---
title: Demos - gRPC application
icon: material/folder-eye-outline
---

# :material-folder-eye-outline: Demo - gRPC application

> Yokai's [showroom](https://github.com/ankorstore/yokai-showroom) provides a [gRPC demo application](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo).

## Overview

This [gRPC demo application](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo) is a simple gRPC API offering a [text transformation service](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/proto/transform.proto).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [gRPC server](../modules/fxgrpcserver.md) module to offer the gRPC API
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

### Layout

This demo application is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
  - `internal/`:
    - `api/`: gRPC API
      - `interceptor/`: gRPC interceptors
      - `service/`: gRPC services
    - `bootstrap.go`: bootstrap
    - `register.go`: dependencies registration
- `proto/`: protobuf definition and stubs

### Makefile

This demo application provides a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make stubs  # generate gRPC stubs with protoc (ex: make stubs from=proto/transform.proto)
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

If no `Transformer` is provided, the transformation configured in `config.transform.default` will be applied.

If you update the [proto definition](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/proto/example.proto), you can run `make stubs from=proto/transform.proto` to regenerate the stubs.

This demo application also provides [reflection](../modules/fxgrpcserver.md#reflection) and [health check ](../modules/fxgrpcserver.md#health-check) services.

### Authentication

This demo application provides example [authentication interceptors](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/internal/api/interceptor/authentication.go).

You can enable authentication in the application [configuration file](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/configs/config.yaml) with `config.authentication.enabled=true`.

If enabled, you need to provide the secret configured in `config.authentication.secret` as context `authorization` metadata.

### Examples

Usage examples with [gRPCurl](https://github.com/fullstorydev/grpcurl):

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
