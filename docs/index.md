---
icon: material/magnify-expand
---

# Yokai: Go framework for backend applications

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-≥1.20-blue)](https://go.dev/)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR)](https://codecov.io/gh/ankorstore/yokai)
[![Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)

![Logo](assets/images/yokai.png){: #overview-logo .skip-glightbox width="350" height="350"}

## :material-target: Goals

Building backend applications with Go is amazing.

But to build production ready applications, you need to put in place a bunch of boilerplate code and efforts, introducing complexity not even related to the logic of your application (like components instantiations, configuration management, observability instrumentation, etc.).

To solve this, Yokai was created with the following goals in mind:

- `Simple`: it is easy to use, configure and test, enabling you to iterate fast and deliver quickly maintainable applications.
- `Modular`: it can be extended with the available Yokai modules, or with your own, to build evolvable applications.
- `Observable`: it comes with built-in logging, tracing and metrics instrumentation, to build reliable applications.

In other words, Yokai let you focus on your application logic, and takes care of the rest.

## :material-map-legend: Architecture

![Architecture](assets/images/architecture.jpg){: #overview-architecture}

- `Yokai Core modules` preloads logging, tracing, metrics and health check instrumentation, and expose a private HTTP server for infrastructure and debugging needs.
- `Yokai extensions modules` can be added to enrich your application features, like public HTTP / gRPC servers, workers, ORM, etc. You can also add your own.
- All this is made available in `Yokai Dependency Injection system` (based on [Fx](https://github.com/uber-go/fx)), on which you can rely to build your application logic.

## :material-rocket-launch-outline: Getting started

Yokai provides ready to use application templates to start your projects:

- gRPC application template
- HTTP application template
- worker application template