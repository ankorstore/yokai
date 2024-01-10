# Yokai

[![Go version](https://img.shields.io/badge/go-%3E%3D1.20-blue)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Simple, modular, and observable Go framework.

<p align="center">
  <img src="docs/images/yokai.png" width="350" height="350" />
</p>

## Documentation

Yokai's documentation will be available soon.

## Modules

| Module                     | Description                                                                                                                                             |
|----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| [config](config)           | Config module based on [Viper](https://github.com/spf13/viper)                                                                                          |
| [generate](generate)       | Generation module based on [Google UUID](https://github.com/google/uuid)                                                                                |
| [healthcheck](healthcheck) | Health check module compatible with [K8s probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) |
| [httpclient](httpclient)   | Http client module based on [net/http](https://pkg.go.dev/net/http)                                                                                     |
| [log](log)                 | Logging module based on [Zerolog](https://github.com/rs/zerolog)                                                                                        |
| [trace](trace)             | Tracing module based on [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go)                                                             |

## Contributing

This repository uses [release-please](https://github.com/googleapis/release-please) to automate Yokai's modules release process.

> [!IMPORTANT]
> You must provide [atomic](https://en.wikipedia.org/wiki/Atomic_commit#Revision_control) and [conventional](https://www.conventionalcommits.org/en/v1.0.0/) commits, since the release process uses them to determinate the releases version and notes to perform.
