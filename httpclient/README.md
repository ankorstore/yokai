# Http Client Module

[![ci](https://github.com/ankorstore/yokai/actions/workflows/httpclient-ci.yml/badge.svg)](https://github.com/ankorstore/yokai/actions/workflows/httpclient-ci.yml)
[![go report](https://goreportcard.com/badge/github.com/ankorstore/yokai/httpclient)](https://goreportcard.com/report/github.com/ankorstore/yokai/httpclient)
[![codecov](https://codecov.io/gh/ankorstore/yokai/graph/badge.svg?token=5s0g5WyseS&flag=httpclient)](https://app.codecov.io/gh/ankorstore/yokai/tree/main/httpclient)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ankorstore/yokai/httpclient)](https://pkg.go.dev/github.com/ankorstore/yokai/httpclient)

> Http client module based on [net/http](https://pkg.go.dev/net/http).

<!-- TOC -->

* [Installation](#installation)
* [Documentation](#documentation)
	* [Requests](#requests)
	* [Transports](#transports)
		* [BaseTransport](#basetransport)
		* [LoggerTransport](#loggertransport)

<!-- TOC -->

## Installation

```shell
go get github.com/ankorstore/yokai/httpclient
```

## Documentation

To create a `http.Client`:

```go
package main

import (
	"time"

	"github.com/ankorstore/yokai/httpclient"
	"github.com/ankorstore/yokai/httpclient/transport"
)

var client, _ = httpclient.NewDefaultHttpClientFactory().Create()

// equivalent to:
var client, _ = httpclient.NewDefaultHttpClientFactory().Create(
	httpclient.WithTransport(transport.NewBaseTransport()), // base http transport (optimized)
	httpclient.WithTimeout(30*time.Second),                 // 30 seconds timeout
	httpclient.WithCheckRedirect(nil),                      // default redirection checks
	httpclient.WithCookieJar(nil),                          // default cookie jar
)
```

### Requests

This module provide some [request helpers](request.go) to ease client requests headers propagation from an incoming
request:

- `CopyObservabilityRequestHeaders` to copy `x-request-id` and `traceparent` headers
- `CopyRequestHeaders` to choose a list of headers to copy

For example:

```go
package main

import (
	"net/http"

	"github.com/ankorstore/yokai/httpclient"
)

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	// create http client
	client, _ := httpclient.NewDefaultHttpClientFactory().Create()

	// build a request to send with the client
	rc, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)

	// propagate observability headers: x-request-id and traceparent
	httpclient.CopyObservabilityRequestHeaders(r, rc)

	// client call
	resp, _ := client.Do(rc)

	// propagate response code
	w.WriteHeader(resp.StatusCode)
}

func main() {
	http.HandleFunc("/", exampleHandler)
	http.ListenAndServe(":8080", nil)
}
```

### Transports

#### BaseTransport

This module provide a [BaseTransport](transport/base.go), optimized regarding max connections handling.

To use it:

```go
package main

import (
	"github.com/ankorstore/yokai/httpclient"
	"github.com/ankorstore/yokai/httpclient/transport"
)

var client, _ = httpclient.NewDefaultHttpClientFactory().Create(
	httpclient.WithTransport(transport.NewBaseTransport()),
)

// equivalent to:
var client, _ = httpclient.NewDefaultHttpClientFactory().Create(
	httpclient.WithTransport(
		transport.NewBaseTransportWithConfig(&transport.BaseTransportConfig{
			MaxIdleConnections:        100,
			MaxConnectionsPerHost:     100,
			MaxIdleConnectionsPerHost: 100,
		}),
	),
)
```

#### LoggerTransport

This module provide a [LoggerTransport](transport/logger.go), able to decorate any `http.RoundTripper` to add logging:

- with requests and response details (and optionally body)
- with configurable log level for each

To use it:

```go
package main

import (
	"github.com/ankorstore/yokai/httpclient"
	"github.com/ankorstore/yokai/httpclient/transport"
	"github.com/rs/zerolog"
)

var client, _ = httpclient.NewDefaultHttpClientFactory().Create(
	httpclient.WithTransport(transport.NewLoggerTransport(nil)),
)

// equivalent to:
var client, _ = httpclient.NewDefaultHttpClientFactory().Create(
	httpclient.WithTransport(
		transport.NewLoggerTransportWithConfig(
			transport.NewBaseTransport(),
			&transport.LoggerTransportConfig{
				LogRequest:                       false,             // to log request details
				LogResponse:                      false,             // to log response details
				LogRequestBody:                   false,             // to log request body (if request details logging enabled)
				LogResponseBody:                  false,             // to log response body (if response details logging enabled)
				LogRequestLevel:                  zerolog.InfoLevel, // log level for request log
				LogResponseLevel:                 zerolog.InfoLevel, // log level for response log
				LogResponseLevelFromResponseCode: false,             // to use response code for response log level
			},
		),
	),
)
```

Note: if no transport is provided for decoration in `transport.NewLoggerTransport(nil)`, the [BaseTransport](transport/base.go) will be used as base transport.
