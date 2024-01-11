# Yokai

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.20-blue)](https://go.dev/)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR)](https://codecov.io/gh/ankorstore/yokai)

> Simple, modular, and observable Go framework.

<p align="center">
  <img src="docs/images/yokai.png" width="350" height="350" />
</p>

## Documentation

Yokai's documentation will be available soon.

## Fx Modules

Yokai is using [Fx](https://github.com/uber-go/fx) for its plugin system.

Yokai's `Fx modules` are the plugins for your Yokai application.

| Fx Module                | Description                        |
|--------------------------|------------------------------------|
| [fxconfig](fxconfig)     | Fx module for [config](config)     |
| [fxgenerate](fxgenerate) | Fx module for [generate](generate) |
| [fxlog](fxlog)           | Fx module for [log](log)           |
| [fxtrace](fxtrace)       | Fx module for [trace](trace)       |

They can also be used in any [Fx](https://github.com/uber-go/fx) based Go application.

## Modules

Yokai's `modules` are the foundation of the framework.

| Module                     | Description                                                                                                                                             |
|----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| [config](config)           | Config module based on [Viper](https://github.com/spf13/viper)                                                                                          |
| [generate](generate)       | Generation module based on [Google UUID](https://github.com/google/uuid)                                                                                |
| [healthcheck](healthcheck) | Health check module compatible with [K8s probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) |
| [httpclient](httpclient)   | Http client module based on [net/http](https://pkg.go.dev/net/http)                                                                                     |
| [httpserver](httpserver)   | Http server module based on [Echo](https://echo.labstack.com/)                                                                                          |
| [log](log)                 | Logging module based on [Zerolog](https://github.com/rs/zerolog)                                                                                        |
| [orm](orm)                 | ORM module based on [Gorm](https://gorm.io/)                                                                                                            |
| [trace](trace)             | Tracing module based on [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go)                                                             |

They can also be used in any Go application (no Yokai or [Fx](https://github.com/uber-go/fx) dependencies).

## Contributing

This repository uses [release-please](https://github.com/googleapis/release-please) to automate Yokai's modules release process.

> [!IMPORTANT]
> You must provide [atomic](https://en.wikipedia.org/wiki/Atomic_commit#Revision_control) and [conventional](https://www.conventionalcommits.org/en/v1.0.0/) commits, since the release process uses them to determinate the releases version and notes to perform.
