---
icon: material/school-outline
---

# :material-school-outline: Tutorial - gRPC application

> How to build, step by step, an gRPC application with Yokai.

## Overview

In this tutorial, we will create an `gRPC` API to offer a `text transformation service`.

You can find a complete implementation in the [gRPC demo application](../demos/grpc-application.md).

## Setup

In this tutorial, we will create our application in the `github.com/foo/bar` example repository.

### Repository creation

To create your `github.com/foo/bar` repository, you can use the [gRPC application template](../getting-started/grpc-application.md).

It provides:

- a ready to extend [Yokai](https://github.com/ankorstore/yokai) application, with the [gRPC server](../modules/fxgrpcserver.md) module installed
- a ready to use [dev environment](https://github.com/ankorstore/yokai-grpc-template/blob/main/docker-compose.yaml), based on [Air](https://github.com/cosmtrek/air) (for live reloading)
- some examples of [service](https://github.com/ankorstore/yokai-grpc-template/blob/main/internal/service/example.go) and [test](https://github.com/ankorstore/yokai-grpc-template/blob/main/internal/service/example_test.go) to get started

### Repository content

Once your repository is created, you should have the following the content:

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `service/`: gRPC service and test examples
	- `bootstrap.go`: bootstrap (modules, lifecycles, etc)
	- `services.go`: dependency injection
- `proto/`: protobuf definition and stubs

And a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make stubs  # generate gRPC stubs with protoc
make test   # run tests
make lint   # run linter
```

## Discovery

You can start your application by running:

```shell
make fresh
```

After a short time, the application will expose:

- `localhost:50051`: application gRPC server
- [http://localhost:8081](http://localhost:8081): application core dashboard

### Example service

When you use the template, an [example gRPC service](https://github.com/ankorstore/yokai-grpc-template/blob/main/internal/service/example.go) is provided, implementing the following [protobuf definition](https://github.com/ankorstore/yokai-grpc-template/blob/main/proto/example.proto):

```protobuf
syntax = "proto3";

package example;

option go_package = "github.com/ankorstore/yokai-grpc-template/proto";

message ExampleRequest {
  string text = 1;
}

message ExampleResponse {
  string text = 1;
}

service ExampleService {
  // Example unary rpc
  rpc ExampleUnary(ExampleRequest) returns (ExampleResponse) {};
  // Example bidi streaming rpc
  rpc ExampleStreaming(stream ExampleRequest) returns (stream ExampleResponse) {};
}
```

You can test to call the  `ExampleService/ExampleUnary` RPC with [grpcurl](https://github.com/fullstorydev/grpcurl):

```shell
grpcurl -plaintext -d '{"text":"test"}' localhost:50051 example.ExampleService/ExampleUnary 
{
  "text": "response from grpc-app: you sent test"
}
```

To ease development, [Air](https://github.com/cosmtrek/air) is watching any changes you perform on `Go code` or `config files` to perform hot reload.

Let's rename your application to `transform-api` by updating `app.name` in the configuration:

```yaml title="config/config.yaml"
app:
	name: transform-api
	# ...
```

Calling again `ExampleService/ExampleUnary ` should now return:

```shell
grpcurl -plaintext -d '{"text":"test"}' localhost:50051 example.ExampleService/ExampleUnary 
{
  "text": "response from transform-api: you sent test"
}
```

### Core dashboard

Yokai is providing a [core](../modules/fxcore.md) dashboard on [http://localhost:8081](http://localhost:8081):

![](../../assets/images/grpc-tutorial-core-dash-light.png#only-light)
![](../../assets/images/grpc-tutorial-core-dash-dark.png#only-dark)

From there, you can get:

- an overview of your application
- information and tooling about your application: build, config, metrics, pprof, etc.
- access to the configured health check endpoints
- access to the loaded modules information (when exposed)

Here we can see for example the [gRPC server](../modules/fxgrpcserver.md) information in the `Modules` section:

- server port
- registered services

See Yokai's [core](../modules/fxcore.md) documentation for more information.

## Implementation

Let's start your application implementation, by:

- providing a [protobuf definition](https://protobuf.dev/getting-started/gotutorial/)
- generating `stubs` with the [protobuf compiler](https://grpc.io/docs/protoc-installation/) for `client` and `server` side implementations
- implementing a gRPC service using those stubs

### Protobuf

#### Definition

We want to create a service that offer 2 RPC:

- `TransformText`: this accepts a `text` and optionally a `transformer`, and returns the `transformed text`.
- `TransformAndSplitText`: this accepts a `text` and optionally a `transformer`, and streams the `transformed text` split in `words`.

Let's update your `proto/example.proto` to define the `TransformTextService`:

```protobuf title="proto/example.proto"
syntax = "proto3";

package example;

option go_package = "github.com/foo/bar/proto";

enum Transformer {
  TRANSFORMER_UNSPECIFIED = 0;
  TRANSFORMER_UPPERCASE = 1;
  TRANSFORMER_LOWERCASE = 2;
}

message TransformTextRequest {
  Transformer transformer = 1;
  string text = 2;
}

message TransformTextResponse {
  string text = 1;
}

service TransformTextService {
  // Unary rpc
  rpc TransformText(TransformTextRequest) returns (TransformTextResponse) {};
  // BiDi rpc
  rpc TransformAndSplitText(stream TransformTextRequest) returns (stream TransformTextResponse) {};
}
```

It provides:

- the `Transformer` enum for the different kind of text transformations
- the `TransformTextRequest` and `TransformTextResponse` messages
- the `TransformTextService` service, with:
    - the `TransformText` RPC (unary)
    - the `TransformAndSplitText` RPC (bidirectional stream)

#### Stubs

You first need to install [protoc](https://grpc.io/docs/languages/go/quickstart/#prerequisites) with the `Go plugins`.

Once available, you can generate the `stubs` for `client` and `server` side implementations with:

```shell
make stubs
```

They will be available in the `proto` folder.

### Service

First, since we updated the `protobuf definition` and the `stubs`, you **must delete** the files:

- `internal/service/example.go`
- `internal/service/example_test.go`

as they are now obsolete.

You also **must remove** the template example service registration from `internal/services.go`.

We can then create a `TransformTextService` in `internal/service/transform.go`:

```go title="internal/service/transform.go"
package service

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/ankorstore/yokai/config"
	"github.com/foo/bar/proto"
)

type TransformTextService struct {
	proto.UnimplementedTransformTextServiceServer
	config *config.Config
}

func NewTransformTextService(cfg *config.Config) *TransformTextService {
	return &TransformTextService{
		config: cfg,
	}
}

func (s *TransformTextService) TransformText(ctx context.Context, in *proto.TransformTextRequest) (*proto.TransformTextResponse, error) {
	return &proto.TransformTextResponse{
		Text: s.transform(in),
	}, nil
}

func (s *TransformTextService) TransformAndSplitText(stream proto.TransformTextService_TransformAndSplitTextServer) error {
	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		split := strings.Split(s.transform(req), " ")

		for _, word := range split {
			err = stream.Send(&proto.TransformTextResponse{
				Text: word,
			})

			if err != nil {
				return err
			}
		}
	}
}

func (s *TransformTextService) transform(in *proto.TransformTextRequest) string {
	switch in.Transformer {
	case proto.Transformer_TRANSFORMER_UPPERCASE:
		return strings.ToUpper(in.Text)
	case proto.Transformer_TRANSFORMER_LOWERCASE:
		return strings.ToLower(in.Text)
	default:
		if strings.ToLower(s.config.GetString("config.transform.default")) == "upper" {
			return strings.ToUpper(in.Text)
		} else {
			return strings.ToLower(in.Text)
		}
	}
}
```

We then need to register the service in `internal/services.go`:

```go title="internal/services.go"
package internal

import (
	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/foo/bar/internal/service"
	"github.com/foo/bar/proto"
	"go.uber.org/fx"
)

// ProvideServices is used to register the application services.
func ProvideServices() fx.Option {
	return fx.Options(
		// gRPC service
		fxgrpcserver.AsGrpcServerService(service.NewTransformTextService, &proto.TransformTextService_ServiceDesc),
	)
}
```

This will automatically inject the `*config.Config` in the `TransformTextService` constructor.

Lets' configure the default `transformer` to apply if none is provided:

```yaml title="configs/config.yaml"
config:
  transform:
    default: upper
```

## Observability

At this stage, we are able to create and list gophers.

To provide a better understanding of what is happening at runtime, let's instrument it with:

- logs
- traces
- metrics

### Logging

With Yokai, `logging` is `contextual`.

This means that you should [propagate the context](https://go.dev/blog/context) and retrieve the [logger](../modules/fxlog.md#usage) from it in order to produce `correlated` logs.

The [HTTP server](../modules/fxhttpserver.md#logging) module automatically injects a logger in the context provided to HTTP handlers.

Let's add logs to our `ListGophersHandler` with `log.CtxLogger()`:

```go title="internal/handler/gopher/list.go"
package gopher

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai/log"
	"github.com/foo/bar/internal/service"
	"github.com/labstack/echo/v4"
)

type ListGophersHandler struct {
	service *service.GopherService
}

func NewListGophersHandler(service *service.GopherService) *ListGophersHandler {
	return &ListGophersHandler{
		service: service,
	}
}

func (h *ListGophersHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		
		log.CtxLogger(ctx).Info().Msg("called ListGophersHandler")

		gophers, err := h.service.List(ctx)
		if err != nil {
			return fmt.Errorf("cannot list gophers: %w", err)
		}

		return c.JSON(http.StatusOK, gophers)
	}
}
```

And to our `GopherService` as well:

```go title="internal/service/gopher.go"
package service

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/foo/bar/internal/model"
	"github.com/foo/bar/internal/repository"
)

type GopherService struct {
	repository *repository.GopherRepository
}

func NewGopherService(repository *repository.GopherRepository) *GopherService {
	return &GopherService{
		repository: repository,
	}
}

// ...

func (s *GopherService) List(ctx context.Context) ([]model.Gopher, error) {
	log.CtxLogger(ctx).Info().Msg("called GopherService.List()")

	return s.repository.FindAll(ctx)
}
```

If you call `[GET] http://localhost:8080/gophers` while observing the logs with `make logs`, you should see:

```shell
INF called GopherService.List() module=httpserver requestID=1a06ab1d-9dec-4424-a3be-23d1c929597a service=gopher-api
INF called ListGophersHandler module=httpserver requestID=1a06ab1d-9dec-4424-a3be-23d1c929597a service=gopher-api
DBG latency="446.978µs" module=httpserver requestID=1a06ab1d-9dec-4424-a3be-23d1c929597a service=gopher-api sqlQuery="SELECT * FROM `gophers` WHERE `gophers`.`deleted_at` IS NULL" sqlRows=1
INF request logger latency="687.925µs" method=GET module=httpserver referer= remoteIp=172.19.0.1 requestID=1a06ab1d-9dec-4424-a3be-23d1c929597a service=gopher-api uri=/gophers
```

You can see that:

- all logs are automatically correlated by `requestID`, allowing you to understand what happened in a specific request scope
- the ORM automatically logged the SQL query, also in this request scope

You can get more information about ORM logging in the [ORM](../modules/fxorm.md#logging) documentation.

### Tracing

With Yokai, `tracing` is `contextual`.

This means that you should [propagate the context](https://go.dev/blog/context) and retrieve the [tracer provider](../modules/fxtrace.md#usage) from it in order to produce `correlated` trace spans.

The [HTTP server](../modules/fxhttpserver.md#logging) module automatically injects the tracer provider in the context provided to HTTP handlers.

First let's activate the [trace](../modules/fxtrace.md#configuration) module exporter to `stdout`:

```yaml title="configs/config.yaml"
modules:
  trace:
    processor: stdout
```

Let's then add trace spans from our `ListGophersHandler` with `trace.CtxTracerProvider()`:

```go title="internal/handler/gopher/list.go"
package gopher

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/foo/bar/internal/service"
	"github.com/labstack/echo/v4"
)

type ListGophersHandler struct {
	service *service.GopherService
}

func NewListGophersHandler(service *service.GopherService) *ListGophersHandler {
	return &ListGophersHandler{
		service: service,
	}
}

func (h *ListGophersHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		ctx, span := trace.CtxTracerProvider(ctx).Tracer("gopher-api").Start(ctx, "ListGophersHandler span")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("called ListGophersHandler")

		gophers, err := h.service.List(ctx)
		if err != nil {
			return fmt.Errorf("cannot list gophers: %w", err)
		}

		return c.JSON(http.StatusOK, gophers)
	}
}
```

If you call `[GET] http://localhost:8080/gophers` while observing with `make logs`, you should see:

```shell
// logs
INF called ListGophersHandler module=httpserver requestID=2c7f596a-e371-4640-83d7-66a3428fd024 service=gopher-api spanID=42331b45b3cfc7bc traceID=6216e1fa6691d994fd980002ede47840
INF called GopherService.List() module=httpserver requestID=2c7f596a-e371-4640-83d7-66a3428fd024 service=gopher-api spanID=42331b45b3cfc7bc traceID=6216e1fa6691d994fd980002ede47840
DBG latency="536.777µs" module=httpserver requestID=2c7f596a-e371-4640-83d7-66a3428fd024 service=gopher-api spanID=64c20e358f00238d sqlQuery="SELECT * FROM `gophers` WHERE `gophers`.`deleted_at` IS NULL" sqlRows=1 traceID=6216e1fa6691d994fd980002ede47840
INF request logger latency="863.981µs" method=GET module=httpserver referer= remoteIp=172.19.0.1 requestID=2c7f596a-e371-4640-83d7-66a3428fd024 service=gopher-api spanID=f857be99a099aa2d status=200 traceID=6216e1fa6691d994fd980002ede47840 uri=/gophers

// trace spans
{"Name":"orm.Query","SpanContext":{"TraceID":"6216e1fa6691d994fd980002ede47840","SpanID":"64c20e358f00238d","TraceFlags":"01","TraceState":"","Remote":false},"Parent":{"TraceID":"6216e1fa6691d994fd980002ede47840","SpanID":"42331b45b3cfc7bc","TraceFlags":"01","TraceState":"","Remote":false},"SpanKind":3,"StartTime":"2024-02-06T11:15:06.611334019Z","EndTime":"2024-02-06T11:15:06.611341607Z","Attributes":[{"Key":"guid:x-request-id","Value":{"Type":"STRING","Value":"2c7f596a-e371-4640-83d7-66a3428fd024"}},{"Key":"db.system","Value":{"Type":"STRING","Value":"mysql"}},{"Key":"db.statement","Value":{"Type":"STRING","Value":"SELECT * FROM `gophers` WHERE `gophers`.`deleted_at` IS NULL"}},{"Key":"db.sql.table","Value":{"Type":"STRING","Value":"gophers"}}],"Events":null,"Links":null,"Status":{"Code":"Unset","Description":""},"DroppedAttributes":0,"DroppedEvents":0,"DroppedLinks":0,"ChildSpanCount":0,"Resource":[{"Key":"service.name","Value":{"Type":"STRING","Value":"gopher-api"}}],"InstrumentationLibrary":{"Name":"orm","Version":"","SchemaURL":""}}
{"Name":"ListGophersHandler span","SpanContext":{"TraceID":"6216e1fa6691d994fd980002ede47840","SpanID":"42331b45b3cfc7bc","TraceFlags":"01","TraceState":"","Remote":false},"Parent":{"TraceID":"6216e1fa6691d994fd980002ede47840","SpanID":"f857be99a099aa2d","TraceFlags":"01","TraceState":"","Remote":false},"SpanKind":1,"StartTime":"2024-02-06T11:15:06.610681301Z","EndTime":"2024-02-06T11:15:06.611506266Z","Attributes":[{"Key":"guid:x-request-id","Value":{"Type":"STRING","Value":"2c7f596a-e371-4640-83d7-66a3428fd024"}}],"Events":null,"Links":null,"Status":{"Code":"Unset","Description":""},"DroppedAttributes":0,"DroppedEvents":0,"DroppedLinks":0,"ChildSpanCount":1,"Resource":[{"Key":"service.name","Value":{"Type":"STRING","Value":"gopher-api"}}],"InstrumentationLibrary":{"Name":"gopher-api","Version":"","SchemaURL":""}}
{"Name":"GET /gophers","SpanContext":{"TraceID":"6216e1fa6691d994fd980002ede47840","SpanID":"f857be99a099aa2d","TraceFlags":"01","TraceState":"","Remote":false},"Parent":{"TraceID":"00000000000000000000000000000000","SpanID":"0000000000000000","TraceFlags":"00","TraceState":"","Remote":false},"SpanKind":2,"StartTime":"2024-02-06T11:15:06.610638183Z","EndTime":"2024-02-06T11:15:06.611598344Z","Attributes":[{"Key":"http.route","Value":{"Type":"STRING","Value":"/gophers"}},{"Key":"http.method","Value":{"Type":"STRING","Value":"GET"}},{"Key":"http.scheme","Value":{"Type":"STRING","Value":"http"}},{"Key":"http.flavor","Value":{"Type":"STRING","Value":"1.1"}},{"Key":"net.host.name","Value":{"Type":"STRING","Value":"gopher-api"}},{"Key":"net.host.port","Value":{"Type":"INT64","Value":8080}},{"Key":"net.sock.peer.addr","Value":{"Type":"STRING","Value":"172.19.0.1"}},{"Key":"net.sock.peer.port","Value":{"Type":"INT64","Value":38054}},{"Key":"http.user_agent","Value":{"Type":"STRING","Value":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"}},{"Key":"guid:x-request-id","Value":{"Type":"STRING","Value":"2c7f596a-e371-4640-83d7-66a3428fd024"}},{"Key":"http.status_code","Value":{"Type":"INT64","Value":200}}],"Events":null,"Links":null,"Status":{"Code":"Unset","Description":""},"DroppedAttributes":0,"DroppedEvents":0,"DroppedLinks":0,"ChildSpanCount":1,"Resource":[{"Key":"service.name","Value":{"Type":"STRING","Value":"gopher-api"}}],"InstrumentationLibrary":{"Name":"gopher-api","Version":"","SchemaURL":""}}
```

Here, we can see on logs side, that:

- they are still correlated by `requestID`
- but they also have the `traceID` and `spanID` fields, correlating logs and trace spans

And on trace spans side, that:

- they are correlated by `TraceID`
- they contain the `guid:x-request-id` attribute matching the logs `requestID`
- the ORM automatically traced the SQL query

You can get more information about ORM tracing in the [ORM](../modules/fxorm.md#tracing) documentation.

### Metrics

Yokai's [metrics](../modules/fxmetrics.md) module is collecting and exposing automatically metrics.

The [core](../modules/fxcore.md) HTTP server of your application will expose them by default on [http://localhost:8081/metrics](http://localhost:8081/metrics), but you can also see them on your [core dashboard](http://localhost:8081):

![](../../assets/images/http-tutorial-core-metrics-light.png#only-light)
![](../../assets/images/http-tutorial-core-metrics-dark.png#only-dark)

You can see that, by default, the [HTTP server](../modules/fxhttpserver.md#metrics) module automatically collects HTTP requests metrics on your HTTP handlers.

Let's now add an example custom metric in our `GopherService` to count the number of times we listed the gophers:

```go title="internal/service/gopher.go"
package service

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/foo/bar/internal/model"
	"github.com/foo/bar/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
)

var GopherListCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "gophers_list_total",
	Help: "The number of times gophers were listed",
})

type GopherService struct {
	repository *repository.GopherRepository
}

func NewGopherService(repository *repository.GopherRepository) *GopherService {
	return &GopherService{
		repository: repository,
	}
}

func (s *GopherService) Create(ctx context.Context, gopher *model.Gopher) error {
	return s.repository.Create(ctx, gopher)
}

func (s *GopherService) List(ctx context.Context) ([]model.Gopher, error) {
	log.CtxLogger(ctx).Info().Msg("called GopherService.List()")

	GopherListCounter.Inc()

	return s.repository.FindAll(ctx)
}
```

To collect this metric, we need to register it with `fxmetrics.AsMetricsCollector()` in `internal/services.go`:

```go title="internal/services.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/orm/healthcheck"
	"github.com/foo/bar/internal/repository"
	"github.com/foo/bar/internal/service"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// orm probe
		fxhealthcheck.AsCheckerProbe(healthcheck.NewOrmProbe),
		// services
		fx.Provide(
			// gophers repository
			repository.NewGopherRepository,
			// gophers service
			service.NewGopherService,
		),
		// gophers list metric
		fxmetrics.AsMetricsCollector(service.GopherListCounter),
	)
}
```

If you call `[GET] http://localhost:8080/gophers`, you can then check the metrics on the [core metrics endpoint](http://localhost:8081/metrics):

```shell title="[GET] http://localhost:8081/metrics"
# ...
# HELP gophers_list_total The number of times gophers were listed
# TYPE gophers_list_total counter
gophers_list_total 1
```

## Testing

At this stage, we are able to create and list gophers, and we have observability signals to monitor this.

The next step is to provide tests for your application, to ensure it's behaving as expected.

### Configuration

Yokai's [bootstrapper](../modules/fxcore.md#bootstrap) provides a `RunTest()` function to start your application in `test` mode.

This will automatically set the env var `APP_ENV=test`, and will [load your test configuration](../modules/fxconfig.md#dynamic-env-overrides).

For our tests, we can configure:

- the [log](../modules/fxlog.md#testing) module to send logs to a `test buffer`
- the [trace](../modules/fxtrace.md#testing) module to send trace spans to a `test exporter`
- the [ORM](../modules/fxorm.md#testing) module to use an [SQLite database](https://www.sqlite.org/index.html), in memory, to make our tests easily portable on any CI pipeline (no need to spin up a MySQL instance)

Let's set the testing configuration in `config/config.test.yaml` and activate the `debug`:

```yaml title="config/config.test.yaml"
app:
  debug: true
modules:
  log:
    level: debug
    output: test
  trace:
    processor:
      type: test
  orm:
    driver: sqlite
    dsn: ":memory:"
```

We also need to update the in bootstrapper the `RunTest()` function to apply your model migrations via `RunFxOrmAutoMigrate()`:

```go title="internal/bootstrap.go"
package internal

import (
	"testing"
	
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxorm"
	"github.com/foo/bar/internal/model"
	"go.uber.org/fx"
)

// ...

func RunTest(tb testing.TB, options ...fx.Option) {
	// ...

	Bootstrapper.RunTestApp(
		tb,
		fx.Options(options...),
		fxorm.RunFxOrmAutoMigrate(&model.Gopher{}),
	)
}
```

This will enable your tests to use the SQLite database automatically with the schema matching your model.

### Implementation

We can now provide `functional` tests for your application endpoints.

Let's create our `TestListGophersHandlerSuccess` test in the `gopher_test` package:

```go title="internal/handler/gopher/list_test.go"
package gopher_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/foo/bar/internal"
	"github.com/foo/bar/internal/model"
	"github.com/foo/bar/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestListGophersHandlerSuccess(t *testing.T) {
	// extraction
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry
	var repo *repository.GopherRepository

	// run test
	internal.RunTest(
		t,
		fx.Populate(&httpServer, &logBuffer, &traceExporter, &metricsRegistry, &repo),
	)

	// populate database
	err := repo.Create(context.Background(), &model.Gopher{
		Name: "bob",
		Job:  "builder",
	})
	assert.NoError(t, err)

	err = repo.Create(context.Background(), &model.Gopher{
		Name: "alice",
		Job:  "doctor",
	})
	assert.NoError(t, err)

	// [GET] /gophers response assertion
	req := httptest.NewRequest(http.MethodGet, "/gophers", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var gophers []*model.Gopher
	err = json.Unmarshal(rec.Body.Bytes(), &gophers)
	assert.NoError(t, err)

	assert.Len(t, gophers, 2)
	assert.Equal(t, gophers[0].Name, "bob")
	assert.Equal(t, gophers[0].Job, "builder")
	assert.Equal(t, gophers[1].Name, "alice")
	assert.Equal(t, gophers[1].Job, "doctor")

	// logs assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called ListGophersHandler",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called GopherService.List()",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(t, traceExporter, "ListGophersHandler span")

	// metrics assertion
	expectedMetric := `
		# HELP gophers_list_total The number of times gophers were listed
		# TYPE gophers_list_total counter
		gophers_list_total 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"gophers_list_total",
	)
	assert.NoError(t, err)
}
```

In this functional test:

- we start the application in test mode
- we populate the test database with fixtures
- we send an HTTP request
- we assert on the HTTP response status and body
- we assert on the observability signals (logs, traces and metrics)

You can then run `make test`:

```shell
=== RUN   TestListGophersHandlerSuccess
--- PASS: TestListGophersHandlerSuccess (0.00s)
PASS
```

This tutorial will only cover testing of the `ListGopherHandler` as example, you need to provide other `functional` tests and the classic `unit` and `integration` tests for the rest of your application.

Thanks to Yokai's [dependency injection system](../modules/fxcore.md#dependency-injection) and [testing tools](../modules/fxcore.md#testing), it's easy to provide mocks as dependencies for your implementations.
