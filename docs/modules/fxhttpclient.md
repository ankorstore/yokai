# HTTP Client Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/fxhttpclient-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/fxhttpclient-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/fxhttpclient)](https://goreportcard.com/report/github.com/ankorstore/yokai/fxhttpclient)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=ghUBlFsjhR&flag=fxhttpclient)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/fxhttpclient)
[![Deps](https://img.shields.io/badge/osi-deps-blue)](https://deps.dev/go/github.com%2Fankorstore%2Fyokai%2Ffxhttpclient)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/fxhttpclient)](https://pkg.go.dev/github.com/ankorstore/yokai/fxhttpclient)

## Overview

Yokai provides a [fxhttpclient](https://github.com/ankorstore/yokai/tree/main/fxhttpclient) module, providing a ready to use [Client](https://pkg.go.dev/net/http#Client) to your application.

It wraps the [httpclient](https://github.com/ankorstore/yokai/tree/main/httpclient) module, based on [net/http](https://pkg.go.dev/net/http).

## Installation

First install the module:

```shell
go get github.com/ankorstore/yokai/fxhttpclient
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpclient"
)


var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxhttpclient module
	fxhttpclient.FxHttpClientModule,
	// ...
)
```

## Usage

This module makes available the [Client](https://pkg.go.dev/net/http#Client) in
Yokai dependency injection system.

To access it, you just need to inject it where needed, for example:

```go title="internal/service/example.go"
package service

import (
	"context"
	"net/http"
)

type ExampleService struct {
	client *http.Client
}

func ExampleService(client *http.Client) *ExampleService {
	return &ExampleService{
		client: client,
	}
}

func (s *ExampleService) Call(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req.WithContext(ctx))
}
```

## Configuration

You can configure the [Client](https://pkg.go.dev/net/http#Client) `timeout`, `transport`, `logging` and `tracing`:

```yaml title="configs/config.yaml"
modules:
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

## Logging

This module enables to log automatically the HTTP requests made by the [Client](https://pkg.go.dev/net/http#Client) and their responses:

```yaml title="configs/config.yaml"
modules:
  http:
    client:
      log:
        request:
          enabled: true              # to log request details, disabled by default
          body: true                 # to add request body to request details, disabled by default
          level: info                # log level for request logging
        response:
          enabled: true              # to log response details, disabled by default
          body: true                 # to add response body to request details, disabled by default
          level: info                # log level for response logging
          level_from_response: true  # to use response code for response logging
```

If `modules.http.client.log.response.level_from_response=true`, the response code will be used to determinate the log level:

- `code < 400`: log level configured in `modules.http.client.log.response.level`
- `400 <= code < 500`: log level `warn`
- `code >= 500`: log level `error`

The HTTP client logging will be based on the [fxlog](fxlog.md) module configuration.

## Tracing

This module enables to trace automatically HTTP the requests made by the [Client](https://pkg.go.dev/net/http#Client):

```yaml title="configs/config.yaml"
modules:
  http:
    client:
      trace:
      	enabled: true # to trace http calls, disabled by default
```

The HTTP client tracing will be based on the [fxtrace](fxtrace.md) module configuration.

## Testing

See [net/http/httptest](https://pkg.go.dev/net/http/httptest) documentation.