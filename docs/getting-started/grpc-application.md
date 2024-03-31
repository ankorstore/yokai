---
title: Getting started - gRPC application
icon: material/rocket-launch-outline
---

# :material-rocket-launch-outline: Getting started - gRPC application

> Yokai provides a ready to use [gRPC application template](https://github.com/ankorstore/yokai-grpc-template) to start your gRPC projects.

## Overview

The [gRPC application template](https://github.com/ankorstore/yokai-grpc-template) provides:

- a ready to extend [Yokai](https://github.com/ankorstore/yokai) application, with the [gRPC server](../modules/fxgrpcserver.md) module installed
- a ready to use [dev environment](https://github.com/ankorstore/yokai-grpc-template/blob/main/docker-compose.yaml), based on [Air](https://github.com/cosmtrek/air) (for live reloading)
- a ready to use [Dockerfile](https://github.com/ankorstore/yokai-grpc-template/blob/main/Dockerfile) for production
- some examples of [service](https://github.com/ankorstore/yokai-grpc-template/blob/main/internal/service/example.go) and [test](https://github.com/ankorstore/yokai-grpc-template/blob/main/internal/service/example_test.go) to get started

### Layout

This template is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `service/`: gRPC service and test examples
	- `bootstrap.go`: bootstrap
	- `register.go`: dependencies registration
- `proto/`: protobuf definition and stubs

### Makefile

This template provides a [Makefile](https://github.com/ankorstore/yokai-grpc-template/blob/main/Makefile):

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make stubs  # generate gRPC stubs with protoc (ex: make stubs from=proto/example.proto)
make test   # run tests
make lint   # run linter
```

## Installation

### With GitHub

You can create your repository [using the GitHub template](https://github.com/new?template_name=yokai-grpc-template&template_owner=ankorstore).

It will automatically rename your project resources, this operation can take a few minutes.

Once ready, after cloning and going into your repository, simply run:

```shell
make fresh
```

### With gonew

You can install [gonew](https://go.dev/blog/gonew), and simply run:

```shell
gonew github.com/ankorstore/yokai-grpc-template github.com/foo/bar
cd bar
make fresh
```

## Usage

Once ready, the application will be available on:

- `localhost:50051` for the application gRPC server
- [http://localhost:8081](http://localhost:8081) for the application core dashboard

If you update the [proto definition](https://github.com/ankorstore/yokai-grpc-template/blob/main/proto/example.proto), you can run `make stubs from=proto/example.proto` to regenerate the stubs.

Usage examples with [gRPCurl](https://github.com/fullstorydev/grpcurl):

- with `ExampleService/ExampleUnary`:

```shell
grpcurl -plaintext -d '{"text":"hello"}' localhost:50051 example.ExampleService/ExampleUnary
{
  "text": "response from grpc-app: you sent hello"
}
```

- with `ExampleService/ExampleStreaming`:

```shell
grpcurl -plaintext -d '@' localhost:50051 example.ExampleService/ExampleStreaming <<EOF
{"text":"hello"}
{"text":"world"}
EOF
{
  "text": "response from grpc-app: you sent hello"
}
{
  "text": "response from grpc-app: you sent world"
}
```

You can use any gRPC clients, for example [Postman](https://learning.postman.com/docs/sending-requests/grpc/grpc-request-interface/) or [Evans](https://github.com/ktr0731/evans).

## Going further

To go further, you can:

- check the [gRPC server](../modules/fxgrpcserver.md) module documentation to learn more about its features
- follow the [gPRC application tutorial](../tutorials/grpc-application.md) to create, step by step, a gRPC application
- test the [gPRC demo application](../demos/grpc-application.md) to see all this in action
