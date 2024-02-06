---
icon: material/school-outline
---

# :material-school-outline: HTTP application tutorial

> How to build, step by step, an HTTP application with Yokai.

## Overview

In this tutorial, we will create an HTTP REST API to manage [gophers](https://go.dev/blog/gopher).

We will create it in the `github.com/foo/bar` example repository.

You can find a complete implementation in the [HTTP application demo](../../applications/demos#http-application-demo).

## Create your repository

To create your `github.com/foo/bar` repository, you can use the [HTTP application template](../../applications/templates#http-application-template).

It provides:

- a ready to extend Yokai application, with the [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module installed
- a ready to use [dev environment](https://github.com/ankorstore/yokai-http-template/blob/main/docker-compose.yaml), based on [Air](https://github.com/cosmtrek/air) (for live reloading)

Once your repository is created, you should have the following the content:

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `handler/`: handler and test examples
	- `bootstrap.go`: bootstrap (modules, lifecycles, etc)
	- `routing.go`: routing
	- `services.go`: dependency injection

And a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make test   # run tests
make lint   # run linter
```