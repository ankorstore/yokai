---
icon: material/rocket-launch-outline
---

# :material-rocket-launch-outline: Getting started - HTTP application

> Yokai provides a ready to use [HTTP application template](https://github.com/ankorstore/yokai-http-template) to start your HTTP projects.

## Overview

The [HTTP application template](https://github.com/ankorstore/yokai-http-template) provides:

- a ready to extend [Yokai](https://github.com/ankorstore/yokai) application, with the [HTTP server](../modules/fxhttpserver.md) module installed
- a ready to use [dev environment](https://github.com/ankorstore/yokai-http-template/blob/main/docker-compose.yaml), based on [Air](https://github.com/cosmtrek/air) (for live reloading)
- some examples of [handler](https://github.com/ankorstore/yokai-http-template/blob/main/internal/handler/example.go) and [test](https://github.com/ankorstore/yokai-http-template/blob/main/internal/handler/example_test.go) to get started

### Layout

This template is following the [standard Go project layout](https://github.com/golang-standards/project-layout):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `handler/`: handler and test examples
	- `bootstrap.go`: bootstrap (modules, lifecycles, etc)
	- `routing.go`: routing
	- `services.go`: dependency injection

### Makefile

This template provides a [Makefile](https://github.com/ankorstore/yokai-http-template/blob/main/Makefile):

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make test   # run tests
make lint   # run linter
```

## Installation

### With GitHub

You can create your repository [using the GitHub template](https://github.com/new?template_name=yokai-http-template&template_owner=ankorstore).

It will automatically rename your project resources, this operation can take a few minutes.

Once ready, after cloning and going into your repository, simply run:

```shell
make fresh
```

### With gonew

You can install [gonew](https://go.dev/blog/gonew), and simply run:

```shell
gonew github.com/ankorstore/yokai-http-template github.com/foo/bar
cd bar
make fresh
```

## Usage

Once ready, the application will be available on:

- [http://localhost:8080](http://localhost:8080) for the application HTTP server
- [http://localhost:8081](http://localhost:8081) for the application core dashboard

## Going further

To go further, you can:

- check the [HTTP server](../modules/fxhttpserver.md) module documentation to learn more about its features
- follow the [HTTP application tutorial](../tutorials/http-application.md) to create, step by step, an HTTP application
- test the [HTTP demo application](../demos/http-application.md) to see all this in action
