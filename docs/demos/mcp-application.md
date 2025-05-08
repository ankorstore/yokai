---
title: Demos - MCP application
icon: material/folder-eye-outline
---

# :material-folder-eye-outline: Demo - MCP application

> Yokai's [showroom](https://github.com/ankorstore/yokai-showroom) provides an [MCP server demo application](https://github.com/ankorstore/yokai-showroom/tree/main/mcp-demo).

## Overview

This [MCP server demo application](https://github.com/ankorstore/yokai-showroom/tree/main/mcp-demo) is a simple [MCP server](https://modelcontextprotocol.io/introduction) to manage [gophers](https://go.dev/blog/gopher).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [MCP server](../modules/fxmcpserver.md) and [SQL](../modules/fxsql.md) modules to offer the gophers MCP server
- a [MySQL](https://www.mysql.com/) container to store the gophers
- a [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) container to interact with the MCP server
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

### Layout

This demo application is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `db/`:
	- `migrations/`: database migrations
	- `seeds/`: database seeds
- `internal/`:
	- `domain/`: domain
		- `model.go`: gophers model
		- `repository.go`: gophers repository
		- `service.go`: gophers service
	- `mcp/`: MCP registrations
		- `prompt/`: MCP prompts
		- `resource/`: MCP resources
		- `tool/`: MCP tools
	- `bootstrap.go`: bootstrap
	- `register.go`: dependencies registration

### Makefile

This demo application provides a `Makefile`:

```
make up      # start the docker compose stack
make down    # stop the docker compose stack
make logs    # stream the docker compose stack logs
make fresh   # refresh the docker compose stack
make migrate # run database migrations
make test    # run tests
make lint    # run linter
```

## Usage

### Start the application

To start the application, simply run:

```shell
make fresh
```

After a short moment, the application will offer:

- [http://localhost:8080/sse](http://localhost:8080/sse): application MCP server (SSE)
- [http://localhost:8081](http://localhost:8081): application core dashboard
- [http://localhost:6274](http://localhost:6274): MCP inspector
- [http://localhost:16686](http://localhost:16686): jaeger UI

### Interact with the application

#### MCP inspector

You can use the provided [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector), available on [http://localhost:6274](http://localhost:6274).

To connect to the MCP server, use:

- `SSE` as transport type
- `http://mcp-demo-app:8080/sse` as URL

Then simply click `Connect`: from there, you will be able to interact with the resources, prompts and tools of the application.

#### MCP hosts

If you use MCP compatible applications like [Cursor](https://www.cursor.com/), or [Claude desktop](https://claude.ai/download), you can register this application as MCP server:

```json
{
  "mcpServers": {
    "mcp-demo-app": {
      "url": "http://localhost:8080/sse"
    }
  }
}
```

Note, if you client does not support remote MCP servers, you can use a [local proxy](https://developers.cloudflare.com/agents/guides/test-remote-mcp-server/#connect-your-remote-mcp-server-to-claude-desktop-via-a-local-proxy):

```json
{
  "mcpServers": {
    "mcp-demo-app": {
      "command": "npx",
      "args": ["mcp-remote", "http://localhost:8080/sse"]
    }
  }
}
```
