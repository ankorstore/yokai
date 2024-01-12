# Fx Http Client Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxhttpclient-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxhttpclient-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxhttpclient)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxhttpclient)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxhttpclient)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxhttpclient)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxhttpclient)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxhttpclient)](https://pkg.go.dev/github.com/ankorstore/yokai/fxhttpclient)

> [Fx](https://uber-go.github.io/fx/) module for [httpclient](https://github.com/ankorstore/yokai/tree/main/httpclient).

<!-- TOC -->
* [Installation](#installation)
* [Features](#features)
* [Documentation](#documentation)
	* [Dependencies](#dependencies)
	* [Loading](#loading)
	* [Configuration](#configuration)
	* [Override](#override)
<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/fxhttpclient
```

## Features

This module provides the possibility to provide to your Fx application a `http.Client` with:

- configurable transport
- automatic and configurable request / response logging
- configurable request / response tracing

## Documentation

### Dependencies

This module is intended to be used alongside:

- the [fxconfig](https://github.com/ankorstore/yokai/tree/main/fxconfig) module
- the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module
- the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module

### Loading

To load the module in your Fx application:

```go
package main

import (
	"net/http"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxhttpclient"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxtrace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxconfig.FxConfigModule,                  // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxhttpclient.FxHttpClientModule,          // load the module
		fx.Invoke(func(httpClient *http.Client) { // invoke the client
			resp, err := httpClient.Get("https://example.com")
		}),
	).Run()
}
```

### Configuration

Configuration reference:

```yaml
# ./configs/config.yaml
app:
  name: app
  env: dev
  version: 0.1.0
  debug: true
modules:
  log:
    level: info
    output: stdout
  trace:
    processor:
      type: stdout
  http:
    client:
      timeout: 30                            # in seconds, 30 by default
      transport:
        max_idle_connections: 100            # 100 by default
        max_connections_per_host: 100        # 100 by default
        max_idle_connections_per_host: 100   # 100 by default
      log:
        request:
          enabled: true                      # to log request details, disabled by default
          body: true                         # to add request body to request details, disabled by default
          level: info                        # log level for request logging
        response:
          enabled: true                      # to log response details, disabled by default
          body: true                         # to add response body to request details, disabled by default
          level: info                        # log level for response logging
          level_from_response: true          # to use response code for response logging
      trace:
        enabled: true                        # to trace http calls, disabled by default
```

If `http.client.log.response.level_from_response=true`, the response code will be used to determinate the log level:

- `code < 400`: log level configured in `http.client.log.response.level`
- `400 <= code < 500`: log level `warn`
- `code >= 500`: log level `error`

Notes:

- the http client logging will be based on the [fxlog](https://github.com/ankorstore/yokai/tree/main/fxlog) module
  configuration
- the http client tracing will be based on the [fxtrace](https://github.com/ankorstore/yokai/tree/main/fxtrace) module
  configuration

### Override

By default, the `http.Client` is created by
the [DefaultHttpClientFactory](https://github.com/ankorstore/yokai/blob/main/httpclient/factory.go).

If needed, you can provide your own factory and override the module:

```go
package main

import (
	"net/http"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxhttpclient"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/httpclient"
	"go.uber.org/fx"
)

type CustomHttpClientFactory struct{}

func NewCustomHttpClientFactory() httpclient.HttpClientFactory {
	return &CustomHttpClientFactory{}
}

func (f *CustomHttpClientFactory) Create(options ...httpclient.HttpClientOption) (*http.Client, error) {
	return http.DefaultClient, nil
}

func main() {
	fx.New(
		fxconfig.FxConfigModule, // load the module dependencies
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxhttpclient.FxHttpClientModule,          // load the module
		fx.Decorate(NewCustomHttpClientFactory),  // override the module with a custom factory
		fx.Invoke(func(httpClient *http.Client) { // invoke the custom client
			// ...
		}),
	).Run()
}
```
