---
icon: material/school-outline
---

# :material-school-outline: HTTP application tutorial

> How to build, step by step, an HTTP application with Yokai.

## Overview

In this tutorial, we will create an `HTTP REST API` to manage [gophers](https://go.dev/blog/gopher).

You can find a complete implementation in the [HTTP application demo](../../applications/demos#http-application-demo).

## Application setup

In this tutorial, we will create our application in the `github.com/foo/bar` example repository.

### Repository creation

To create your `github.com/foo/bar` repository, you can use the [HTTP application template](../../applications/templates#http-application-template).

It provides:

- a ready to extend Yokai application, with the [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module installed
- a ready to use [dev environment](https://github.com/ankorstore/yokai-http-template/blob/main/docker-compose.yaml), based on [Air](https://github.com/cosmtrek/air) (for live reloading)

### Repository content

Once your repository is created, you should have the following the content:

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `handler/`: handler and test examples
	- `bootstrap.go`: bootstrap (modules, lifecycles, etc)
	- `routing.go`: routing
	- `services.go`: dependency injection

And a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make test   # run tests
make lint   # run linter
```

## Application discovery

You can start your application by running:

```shell
make fresh
```

After a short time, the application will expose:

- [http://localhost:8080](http://localhost:8080): application example endpoint
- [http://localhost:8081](http://localhost:8081): application core dashboard

### Example endpoint

When you use the template, an example endpoint is provided on [http://localhost:8080](http://localhost:8080):

```shell title="GET http://localhost:8080"
Welcome to http-app.
```

To ease development, [Air](https://github.com/cosmtrek/air) is watching any changes you perform on `Go code` or `config files` to perform hot reload.

Let's rename your application in `gopher-api` by updating `app.name`:

```yaml title="config/config.yaml"
app:
	name: gopher-api
	version: 0.1.0
	# ...
```

Calling again [http://localhost:8080](http://localhost:8080) should now return:

```shell title="GET http://localhost:8080"
Welcome to gopher-api.
```

### Core dashboard

Yokai is providing a core dashboard on [http://localhost:8081](http://localhost:8081):

![](../../assets/images/http-tutorial-core-dash-light.png#only-light)
![](../../assets/images/http-tutorial-core-dash-dark.png#only-dark)

From there, you can get:

- an overview of your application
- information and tooling about your application: build, config, metrics, pprof, etc.
- access to the configured health check endpoints
- access to the loaded modules information (when exposed)

Here we can see for example the [fxhttpserver](../modules/fxhttpserver.md) information in the `Modules` section:

- server port
- active routes
- error handler
- etc

See [fxcore](../modules/fxcore.md) documentation for more information.

## Application implementation

Let's start your application implementation, by:

- adding database support
- implementing endpoints to create and list gophers

### Database setup

#### MySQL installation

Let's update your `docker-compose.yaml` to add a [MySQL](https://www.mysql.com/) container to your stack:

```yaml title="docker-compose.yaml"
version: '3.9'

services:
  gohper-api-app:
    container_name: gohper-api-app
    build:
      dockerfile: dev.Dockerfile
      context: .
    networks:
      - gohper-api
    ports:
      - "8080:8080"
      - "8081:8081"
    expose:
      - "8080"
      - "8081"
    volumes:
      - .:/app
    env_file:
      - .env

  gohper-api-database:
    container_name: gohper-api-database
    image: mysql:8
    restart: always
    networks:
      - gohper-api
    volumes:
      - gohper-api-database-data:/var/lib/mysql
    env_file:
      - .env

volumes:
  gohper-api-database-data:
    driver: local

networks:
  gohper-api:
    driver: bridge
```

And the configuration in your `.env` file:

```env title=".env"
APP_ENV=dev
APP_DEBUG=true
MYSQL_HOST=gohper-api-database
MYSQL_PORT=3306
MYSQL_DATABASE=gohper-api
MYSQL_USER=user
MYSQL_PASSWORD=password
MYSQL_ROOT_PASSWORD=rootpassword
```

You can then refresh your stack to bring this up:

```shell
make fresh
```

#### ORM module installation

Yokai provides the [fxorm](../modules/fxorm.md) module, extending your application with [GORM](https://gorm.io/).

You can install it:

```shell
go get github.com/ankorstore/yokai/fxorm
```

Then activate it in your application bootstrapper:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxorm"
)

// ...

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// load fxorm module
	fxorm.FxOrmModule,
	// ...
)
```

You can then provide the module configuration:

```yaml title="configs/config.yaml"
modules:
  orm:
    driver: mysql
    dsn: ${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}?parseTime=true
    log:
      enabled: true
      level: info
      values: true
    trace:
      enabled: true
      values: true
```

#### Model creation

To manage [gophers](https://go.dev/blog/gopher), we need to [create a model](https://gorm.io/docs/models.html):

```go title="internal/model/gopher.go"
package model

import (
	"gorm.io/gorm"
)

type Gopher struct {
	gorm.Model
	Name string `json:"name" form:"name"`
	Job  string `json:"job" form:"job"`
}
```

#### Model migrations

The [fxorm](../modules/fxorm.md) module [provides ways](../modules/fxorm.md#migrations) to apply your [schemas migrations](https://gorm.io/docs/migration.html).

To run the migrations automatically at bootstrap, we just need to pass our model to `RunFxOrmAutoMigrate()`:

```go title="internal/bootstrap.go"
package internal

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxorm"
	"github.com/foo/bar/internal/model"
)

// ...

func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp(
		// run ORM migrations for the Gopher model
		fxorm.RunFxOrmAutoMigrate(&model.Gopher{}),
	)
}
```

If you check the logs with `make logs`, you should see the migration happening:

```shell
INF starting ORM auto migration service=gopher-api
INF ORM auto migration success service=gopher-api
```

#### Health check

Yokai's [fxhealthcheck](../modules/fxhealthcheck.md) module allows the core HTTP server to expose health check endpoints, useful if your application runs on [Kunernetes](https://kubernetes.io/). It will execute the [registered probes](../modules/fxhealthcheck.md#usage).

The [fxorm](../modules/fxorm.md#health-check) provides a ready to use [OrmProbe](https://github.com/ankorstore/yokai/blob/main/orm/healthcheck/probe.go), that will `ping` the database connection to check if it's healthy.

To activate it, you can use the `fxhealthcheck.AsCheckerProbe()` function in `internal/services.go`:

```go title="internal/services.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/orm/healthcheck"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// orm probe
		fxhealthcheck.AsCheckerProbe(healthcheck.NewOrmProbe),
	)
}
```

This will register the ORM probe for `startup`, `liveness` and `readiness` checks.

You can check that it's properly activated on the [core dashboard](http://localhost:8081):

![](../../assets/images/http-tutorial-core-hc-light.png#only-light)
![](../../assets/images/http-tutorial-core-hc-dark.png#only-dark)

### Repository implementation

We can create a `GopherRepository` to manage our gophers, with:

- the `Create()` function to create a gopher 
- and the `FindAll()` function to list all gophers

```go title="internal/repository/gopher.go"
package repository

import (
	"context"
	"sync"

	"github.com/foo/bar/internal/model"
	"gorm.io/gorm"
)

type GopherRepository struct {
	mutex sync.Mutex
	db    *gorm.DB
}

func NewGopherRepository(db *gorm.DB) *GopherRepository {
	return &GopherRepository{
		db: db,
	}
}

func (r *GopherRepository) Create(ctx context.Context, gopher *model.Gopher) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	res := r.db.WithContext(ctx).Create(gopher)

	return res.Error
}

func (r *GopherRepository) FindAll(ctx context.Context) ([]model.Gopher, error) {
	var gophers []model.Gopher

	res := r.db.WithContext(ctx).Find(&gophers)
	if res.Error != nil {
		return nil, res.Error
	}

	return gophers, nil
}
```

We then need to register the repository in `internal/services.go`:

```go title="internal/services.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/orm/healthcheck"
	"github.com/foo/bar/internal/repository"
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
		),
	)
}
```

This will automatically inject the `*gorm.DB` in the `GopherRepository` constructor.

### Service implementation

Now that we have a repository, let's create a `GopherService`, with:

- the `Create()` function to create a gopher
- and the `List()` function to list all gophers

```go title="internal/service/gopher.go"
package service

import (
	"context"

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

func (s *GopherService) Create(ctx context.Context, gopher *model.Gopher) error {
	return s.repository.Create(ctx, gopher)
}

func (s *GopherService) List(ctx context.Context) ([]model.Gopher, error) {
	return s.repository.FindAll(ctx)
}
```

We then need to register the service in `internal/services.go`:

```go title="internal/services.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
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
	)
}
```

This will automatically inject the `*repository.GopherRepository` in the `GopherService` constructor.

### HTTP handlers implementation

Now that we have a `GopherService` able to create and list gophers, let's expose it via HTTP endpoints in your application.

#### Create HTTP handler

Let's create a `CreateGopherHandler` to handle requests on `[POST] /gophers` to create gophers:

```go title="internal/handler/gopher/create.go"
package gopher

import (
	"fmt"
	"net/http"

	"github.com/foo/bar/internal/model"
	"github.com/foo/bar/internal/service"
	"github.com/labstack/echo/v4"
)

type CreateGopherHandler struct {
	service *service.GopherService
}

func NewCreateGopherHandler(service *service.GopherService) *CreateGopherHandler {
	return &CreateGopherHandler{
		service: service,
	}
}

func (h *CreateGopherHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		gopher := new(model.Gopher)
		if err := c.Bind(gopher); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("cannot bind gopher: %v", err))
		}

		err := h.service.Create(c.Request().Context(), gopher)
		if err != nil {
			return fmt.Errorf("cannot create gopher: %w", err)
		}

		return c.JSON(http.StatusCreated, gopher)
	}
}
```

We then need to register the handler for `[POST] /gophers` in `internal/routing.go`:

```go title="internal/routing.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/foo/bar/internal/handler"
	"github.com/foo/bar/internal/handler/gopher"
	"go.uber.org/fx"
)

func ProvideRouting() fx.Option {
	return fx.Options(
		fxhttpserver.AsHandler("GET", "", handler.NewExampleHandler),
		// gopher creation
		fxhttpserver.AsHandler("POST", "/gophers", gopher.NewCreateGopherHandler),
	)
}
```

Let's try to call it:

```shell title="POST http://localhost:8080/gophers"
curl -X POST http://localhost:8080/gophers -H 'Content-Type: application/json' -d '{"name":"bob","job":"builder"}'                   
{
  "ID": 1,
  "CreatedAt": "2024-02-06T10:29:26.497Z",
  "UpdatedAt": "2024-02-06T10:29:26.497Z",
  "DeletedAt": null,
  "name": "bob",
  "job": "builder"
}
```

You should receive a response with status `201` (created), and with the created gopher representation.

You can check the [fxhttpserver](../modules/fxhttpserver.md#handlers-registration) module documentation if you need more information about registering handlers.

#### List HTTP handler

Let's now create a `ListGopherHandler` to handle requests on `[GET] /gophers` to list gophers:

```go title="internal/handler/gopher/list.go"
package gopher

import (
	"fmt"
	"net/http"

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
		gophers, err := h.service.List(c.Request().Context())
		if err != nil {
			return fmt.Errorf("cannot list gophers: %w", err)
		}

		return c.JSON(http.StatusOK, gophers)
	}
}
```

We then need to register the handler for `[GET] /gophers` in `internal/routing.go`.

We can group our handlers registration with `fxhttpserver.AsHandlersGroup()`:

```go title="internal/routing.go"
package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/foo/bar/internal/handler"
	"github.com/foo/bar/internal/handler/gopher"
	"go.uber.org/fx"
)

func ProvideRouting() fx.Option {
	return fx.Options(
		fxhttpserver.AsHandler("GET", "", handler.NewExampleHandler),
		// gopher handlers group
		fxhttpserver.AsHandlersGroup(
			"/gophers",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "", gopher.NewListGophersHandler),
				fxhttpserver.NewHandlerRegistration("POST", "", gopher.NewCreateGopherHandler),
			},
		),
	)
}
```

You can check the [fxhttpserver](../modules/fxhttpserver.md#handlers-groups-registration) module documentation if you need more information about registering handlers groups.

Let's try to call it:

```shell title="GET http://localhost:8080/gophers"
curl http://localhost:8080/gophers                                                                                
[
  {
    "ID": 1,
    "CreatedAt": "2024-02-06T10:29:26.497Z",
    "UpdatedAt": "2024-02-06T10:29:26.497Z",
    "DeletedAt": null,
    "name": "bob",
    "job": "builder"
  }
]
```

You should receive a response with status `200` (ok), and with a list of gophers containing the one previously created.

## Application observability

At this stage, we are able to create and list gophers.

To provide a better understanding of what is happening at runtime, let's instrument it with:

- logs
- traces
- metrics

### Application logging

With Yokai, `logging` is `contextual`.

This means that you should [propagate the context](https://go.dev/blog/context) and retrieve the [logger](../modules/fxlog.md#usage) from it in order to produce `correlated` logs.

The [fxhttpserver](../modules/fxhttpserver.md#logging) module automatically injects a logger in the context provided to HTTP handlers.

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

You can get more information about ORM logging in the [fxorm](../modules/fxorm.md#logging) documentation.

### Application tracing

With Yokai, `tracing` is `contextual`.

This means that you should [propagate the context](https://go.dev/blog/context) and retrieve the [tracer provider](../modules/fxtrace.md#usage) from it in order to produce `correlated` trace spans.

The [fxhttpserver](../modules/fxhttpserver.md#logging) module automatically injects the tracer provider in the context provided to HTTP handlers.

First let's activate the [fxtrace](../modules/fxtrace.md#configuration) exporter to `stdout`:

```yaml title="configs/config.yaml"
modules:
  trace:
    processor: stdout
```

Let's then add trace spans to our `ListGophersHandler` with `trace.CtxTracerProvider()`:

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

You can get more information about ORM tracing in the [fxorm](../modules/fxorm.md#tracing) documentation.

### Application metrics

Yokai, via the [fxmetrics](../modules/fxmetrics.md) module, is collecting and exposing automatically metrics.

The core HTTP server of your application will expose them by default on [http://localhost:8081/metrics](http://localhost:8081/metrics), but you can also see them on your [core dashboard](http://localhost:8081):

![](../../assets/images/http-tutorial-core-metrics-light.png#only-light)
![](../../assets/images/http-tutorial-core-metrics-dark.png#only-dark)

You can see that, by default, the [fxhttpserver](../modules/fxhttpserver.md#metrics) module automatically collects HTTP requests metrics on your HTTP handlers.

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

## Application testing

At this stage, we are able to create and list gophers, and we have observability signals to monitor this.

The next step is to provide tests for your application, to ensure it's behaving as expected.

### Tests configuration

Yokai's [bootstrapper](../modules/fxcore.md#bootstrap) provides a `RunTest()` function to start your application in `test` mode.

This will automatically set the env var `APP_ENV=test`, and will [load your test configuration](../modules/fxconfig.md#dynamic-env-overrides).

For our tests, we can configure:

- the [fxlog](../modules/fxlog.md#testing) module to send logs to a `test buffer`
- the [fxtrace](../modules/fxtrace.md#testing) module to send trace spans to a `test exporter`
- the [fxorm](../modules/fxorm.md#testing) module to use an [SQLite database](https://www.sqlite.org/index.html), in memory, to make our tests easily portable on any CI pipeline (no need to spin up a MySQL instance)

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

### Tests implementation

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
