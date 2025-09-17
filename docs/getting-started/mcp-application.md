---
title: Getting started - MCP application
icon: material/rocket-launch-outline
---

# :material-rocket-launch-outline: Getting started - MCP application

> Yokai provides a ready to use [MCP server application template](https://github.com/ankorstore/yokai-mcp-template) to start your MCP projects.

## Overview

The [MCP server application template](https://github.com/ankorstore/yokai-mcp-template)  provides:

- a ready to extend [Yokai](https://github.com/ankorstore/yokai) application, with the [MCP server](../modules/fxmcpserver.md) module installed
- a ready to use [dev environment](https://github.com/ankorstore/yokai-http-template/blob/main/docker-compose.yaml), based on [Air](https://github.com/air-verse/air) (for live reloading)
- a ready to use [Dockerfile](https://github.com/ankorstore/yokai-http-template/blob/main/Dockerfile) for production
- some examples of [MCP tool](https://github.com/ankorstore/yokai-mcp-template/blob/main/internal/tool/example.go) and [test](https://github.com/ankorstore/yokai-mcp-template/blob/main/internal/tool/example_test.go) to get started

### Layout

This template is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `tool/`: MCP tool and test examples
	- `bootstrap.go`: bootstrap
	- `register.go`: dependencies registration

### Makefile

This template provides a [Makefile](https://github.com/ankorstore/yokai-http-template/blob/main/Makefile):

```
make up      # start the docker compose stack
make down    # stop the docker compose stack
make logs    # stream the docker compose stack logs
make fresh   # refresh the docker compose stack
make test    # run tests
make lint    # run linter
```

## Installation

### With GitHub

You can create your repository [using the GitHub template](https://github.com/new?template_name=yokai-mcp-template&template_owner=ankorstore).

It will automatically rename your project resources, this operation can take a few minutes.

Once ready, after cloning and going into your repository, simply run:

```shell
make fresh
```

### With gonew

You can install [gonew](https://go.dev/blog/gonew), and simply run:

```shell
gonew github.com/ankorstore/yokai-mcp-template github.com/foo/bar
cd bar
make fresh
```

## Usage

Once ready, the application will be available on:

- [http://localhost:8080/sse](http://localhost:8080/sse) for the application MCP server
- [http://localhost:8081](http://localhost:8081) for the application core dashboard

## Going further

To go further, you can:

- check the [MCP server](../modules/fxmcpserver.md) module documentation to learn more about its features
- follow the [MCP application tutorial](../tutorials/mcp-application.md) to create, step by step, an MCP server application
- test the [MCP demo application](../demos/mcp-application.md) to see all this in action
