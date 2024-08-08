---
title: Getting started - Worker application
icon: material/rocket-launch-outline
---

# :material-rocket-launch-outline: Getting started - worker application

> Yokai provides a ready to use [worker application template](https://github.com/ankorstore/yokai-worker-template) to start your worker projects.

## Overview

The [worker application template](https://github.com/ankorstore/yokai-worker-template) provides:

- a ready to extend [Yokai](https://github.com/ankorstore/yokai) application, with the [worker](../modules/fxworker.md) module installed
- a ready to use [dev environment](https://github.com/ankorstore/yokai-worker-template/blob/main/docker-compose.yaml), based on [Air](https://github.com/air-verse/air) (for live reloading)
- a ready to use [Dockerfile](https://github.com/ankorstore/yokai-worker-template/blob/main/Dockerfile) for production
- some examples of [worker](https://github.com/ankorstore/yokai-worker-template/blob/main/internal/worker/example.go) and [test](https://github.com/ankorstore/yokai-worker-template/blob/main/internal/worker/example_test.go) to get started

### Layout

This template is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `worker/`: worker and test examples
	- `bootstrap.go`: bootstrap
	- `register.go`: dependencies registration

### Makefile

This template provides a [Makefile](https://github.com/ankorstore/yokai-worker-template/blob/main/Makefile):

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

You can create your repository [using the GitHub template](https://github.com/new?template_name=yokai-worker-template&template_owner=ankorstore).

It will automatically rename your project resources, this operation can take a few minutes.

Once ready, after cloning and going into your repository, simply run:

```shell
make fresh
```

### With gonew

You can install [gonew](https://go.dev/blog/gonew), and simply run:

```shell
gonew github.com/ankorstore/yokai-worker-template github.com/foo/bar
cd bar
make fresh
```

## Usage

Once ready, the application core dashboard will be available on [http://localhost:8081](http://localhost:8081).

To see the [provided example worker](https://github.com/ankorstore/yokai-worker-template/blob/main/internal/worker/example.go) in action, simply run:

```shell
make logs
```

## Going further

To go further, you can:

- check the [worker](../modules/fxworker.md) module documentation to learn more about its features
- follow the [worker application tutorial](../tutorials/worker-application.md) to create, step by step, a worker application
- test the [worker demo application](../demos/worker-application.md) to see all this in action

