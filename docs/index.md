---
icon: material/book-open-page-variant-outline
---

# Yokai

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-≥1.20-blue)](https://go.dev/)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR)](https://codecov.io/gh/ankorstore/yokai)
[![Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)

> A `simple`, `modular` and `observable` Go framework for `backend applications`.

![Logo](assets/images/yokai.png){: #overview-logo .skip-glightbox width="300" height="300"}

## Goals

Building backend applications with [Go](https://go.dev/) is amazing.

But to build production-grade applications, you need to put in place a bunch of efforts and boilerplate code, introducing complexity not even related to the logic of your application (like dependencies wiring, configuration management, observability instrumentation, etc.).

To solve this, Yokai was created with the following goals in mind:

- `Simple`: it is easy to use, configure and test, enabling you to iterate fast and deliver quickly maintainable applications.
- `Modular`: it can be extended with the available Yokai modules, or with your own, to build evolvable applications.
- `Observable`: it comes with built-in logging, tracing and metrics instrumentation, to build reliable applications.

In other words, Yokai lets you focus on your application logic, while taking care of the rest.

## Overview

### Architecture

![Architecture](assets/images/architecture.jpg){: #overview-architecture}

- Yokai `core modules` preloads logging, tracing, metrics and health check instrumentation, and expose a private HTTP server for infrastructure and debugging needs.
- Yokai `extensions modules` can be added to enrich your application features, like public HTTP / gRPC servers, workers, ORM, etc. You can also add the [contrib modules](https://github.com/ankorstore/yokai-contrib) or your own.
- Those are made available in Yokai `dependency injection system`, on which you can rely to build your application logic.

### Foundations

Yokai was built using `robust` and `well known` Go libraries, like:

- [Echo](https://github.com/labstack/echo) for HTTP servers
- [gRPC-go](https://github.com/grpc/grpc-go) for gRPC servers
- [Viper](https://github.com/spf13/viper) for configuration management
- [OTEL](https://github.com/open-telemetry/opentelemetry-go) for observability instrumentation
- [Fx](https://github.com/uber-go/fx) for dependency injection system
- and more...


### Extensions

Yokai's `extension system` enables you to `enrich` your application features with:

- the Yokai `built-in` modules
- the Yokai [contrib modules](https://github.com/ankorstore/yokai-contrib)
- your own modules

## Getting started

Yokai provides ready to use `application templates` to start your projects:

- for [gRPC applications](getting-started/grpc-application.md)
- for [HTTP applications](getting-started/http-application.md)
- for [worker applications](getting-started/worker-application.md)