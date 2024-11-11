package fxhttpserver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/httpserver"
	httpservermiddleware "github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gommonlog "github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const (
	ModuleName     = "httpserver"
	DefaultAddress = ":8080"
)

// FxHttpServerModule is the [Fx] httpserver module.
//
// [Fx]: https://github.com/uber-go/fx
var FxHttpServerModule = fx.Module(
	ModuleName,
	fx.Provide(
		httpserver.NewDefaultHttpServerFactory,
		NewFxHttpServerRegistry,
		NewFxHttpServer,
		fx.Annotate(
			NewFxHttpServerModuleInfo,
			fx.As(new(interface{})),
			fx.ResultTags(`group:"core-module-infos"`),
		),
	),
)

// FxHttpServerParam allows injection of the required dependencies in [NewFxHttpServer].
type FxHttpServerParam struct {
	fx.In
	LifeCycle       fx.Lifecycle
	Factory         httpserver.HttpServerFactory
	Generator       uuid.UuidGenerator
	Registry        *HttpServerRegistry
	Config          *config.Config
	Logger          *log.Logger
	TracerProvider  trace.TracerProvider
	MetricsRegistry *prometheus.Registry
}

// NewFxHttpServer returns a new [echo.Echo].
func NewFxHttpServer(p FxHttpServerParam) (*echo.Echo, error) {
	appDebug := p.Config.AppDebug()

	// logger
	echoLogger := httpserver.NewEchoLogger(
		log.FromZerolog(p.Logger.ToZerolog().With().Str("module", ModuleName).Logger()),
	)

	// renderer
	var renderer echo.Renderer
	if p.Config.GetBool("modules.http.server.templates.enabled") {
		renderer = httpserver.NewHtmlTemplateRenderer(p.Config.GetString("modules.http.server.templates.path"))
	}

	// server
	httpServer, err := p.Factory.Create(
		httpserver.WithDebug(appDebug),
		httpserver.WithBanner(false),
		httpserver.WithLogger(echoLogger),
		httpserver.WithRenderer(renderer),
		httpserver.WithHttpErrorHandler(
			httpserver.JsonErrorHandler(
				p.Config.GetBool("modules.http.server.errors.obfuscate") || !appDebug,
				p.Config.GetBool("modules.http.server.errors.stack") || appDebug,
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http server: %w", err)
	}

	// middlewares registrations
	httpServer = withDefaultMiddlewares(httpServer, p)

	// groups, handlers & middlewares registrations
	httpServer, err = withRegisteredResources(httpServer, p)
	if err != nil {
		return httpServer, fmt.Errorf("failed to register http server resources: %w", err)
	}

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if !p.Config.IsTestEnv() {
				address := p.Config.GetString("modules.http.server.address")
				if address == "" {
					address = DefaultAddress
				}

				//nolint:errcheck
				go httpServer.Start(address)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if !p.Config.IsTestEnv() {
				return httpServer.Shutdown(ctx)
			}

			return nil
		},
	})

	return httpServer, nil
}

func withDefaultMiddlewares(httpServer *echo.Echo, p FxHttpServerParam) *echo.Echo {
	// request id middleware
	httpServer.Use(httpservermiddleware.RequestIdMiddlewareWithConfig(
		httpservermiddleware.RequestIdMiddlewareConfig{
			Generator: p.Generator,
		},
	))

	// request tracer middleware
	if p.Config.GetBool("modules.http.server.trace.enabled") {
		httpServer.Use(httpservermiddleware.RequestTracerMiddlewareWithConfig(
			p.Config.AppName(),
			httpservermiddleware.RequestTracerMiddlewareConfig{
				TracerProvider:              httpserver.AnnotateTracerProvider(p.TracerProvider),
				RequestUriPrefixesToExclude: p.Config.GetStringSlice("modules.http.server.trace.exclude"),
			},
		))
	}

	// request logger middleware
	requestHeadersToLog := map[string]string{
		httpservermiddleware.HeaderXRequestId: httpservermiddleware.LogFieldRequestId,
	}

	for headerName, fieldName := range p.Config.GetStringMapString("modules.http.server.log.headers") {
		requestHeadersToLog[headerName] = fieldName
	}

	httpServer.Use(httpservermiddleware.RequestLoggerMiddlewareWithConfig(
		httpservermiddleware.RequestLoggerMiddlewareConfig{
			RequestHeadersToLog:             requestHeadersToLog,
			RequestUriPrefixesToExclude:     p.Config.GetStringSlice("modules.http.server.log.exclude"),
			LogLevelFromResponseOrErrorCode: p.Config.GetBool("modules.http.server.log.level_from_response"),
		},
	))

	// request metrics middleware
	if p.Config.GetBool("modules.http.server.metrics.collect.enabled") {
		namespace := Sanitize(p.Config.GetString("modules.http.server.metrics.collect.namespace"))
		subsystem := Sanitize(p.Config.GetString("modules.http.server.metrics.collect.subsystem"))

		var buckets []float64
		if bucketsConfig := p.Config.GetString("modules.http.server.metrics.buckets"); bucketsConfig != "" {
			for _, s := range Split(bucketsConfig) {
				f, err := strconv.ParseFloat(s, 64)
				if err == nil {
					buckets = append(buckets, f)
				}
			}
		}

		metricsMiddlewareConfig := httpservermiddleware.RequestMetricsMiddlewareConfig{
			Registry:                p.MetricsRegistry,
			Namespace:               namespace,
			Subsystem:               subsystem,
			Buckets:                 buckets,
			NormalizeRequestPath:    p.Config.GetBool("modules.http.server.metrics.normalize.request_path"),
			NormalizeResponseStatus: p.Config.GetBool("modules.http.server.metrics.normalize.response_status"),
		}

		httpServer.Use(httpservermiddleware.RequestMetricsMiddlewareWithConfig(metricsMiddlewareConfig))
	}

	// recovery middleware
	httpServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableErrorHandler: true,
		LogLevel:            gommonlog.ERROR,
	}))

	return httpServer
}

//nolint:cyclop
func withRegisteredResources(httpServer *echo.Echo, p FxHttpServerParam) (*echo.Echo, error) {
	// register handler groups
	resolvedHandlersGroups, err := p.Registry.ResolveHandlersGroups()
	if err != nil {
		httpServer.Logger.Errorf("cannot resolve router handlers groups: %v", err)
	}

	for _, g := range resolvedHandlersGroups {
		group := httpServer.Group(g.Prefix(), g.Middlewares()...)

		for _, h := range g.Handlers() {
			methods, err := ExtractMethods(h.Method())
			if err != nil {
				return httpServer, err
			}

			for _, method := range methods {
				group.Add(
					method,
					h.Path(),
					h.Handler(),
					h.Middlewares()...,
				)

				httpServer.Logger.Debugf("registering handler in group for [%s] %s%s", method, g.Prefix(), h.Path())
			}
		}

		httpServer.Logger.Debugf("registered handlers group for prefix %s", g.Prefix())
	}

	// register middlewares
	resolvedMiddlewares, err := p.Registry.ResolveMiddlewares()
	if err != nil {
		httpServer.Logger.Errorf("cannot resolve router middlewares: %v", err)
	}

	for _, m := range resolvedMiddlewares {
		if m.Kind() == GlobalPre {
			httpServer.Pre(m.Middleware())
		}

		if m.Kind() == GlobalUse {
			httpServer.Use(m.Middleware())
		}

		httpServer.Logger.Debugf("registered %s middleware %T", m.Kind().String(), m.Middleware())
	}

	// register handlers
	resolvedHandlers, err := p.Registry.ResolveHandlers()
	if err != nil {
		httpServer.Logger.Errorf("cannot resolve router handlers: %v", err)
	}

	for _, h := range resolvedHandlers {
		methods, err := ExtractMethods(h.Method())
		if err != nil {
			return httpServer, err
		}

		for _, method := range methods {
			httpServer.Add(
				method,
				h.Path(),
				h.Handler(),
				h.Middlewares()...,
			)

			httpServer.Logger.Debugf("registered handler for [%s] %s", method, h.Path())
		}
	}

	return httpServer, nil
}
