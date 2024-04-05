package fxcore

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/handler"
	httpservermiddleware "github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/ankorstore/yokai/log"
	"github.com/arl/statsviz"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const (
	ModuleName                      = "core"
	DefaultAddress                  = ":8081"
	DefaultMetricsPath              = "/metrics"
	DefaultHealthCheckStartupPath   = "/healthz"
	DefaultHealthCheckLivenessPath  = "/livez"
	DefaultHealthCheckReadinessPath = "/readyz"
	DefaultDebugConfigPath          = "/debug/config"
	DefaultDebugPProfPath           = "/debug/pprof"
	DefaultDebugBuildPath           = "/debug/build"
	DefaultDebugRoutesPath          = "/debug/routes"
	DefaultDebugStatsPath           = "/debug/stats"
	DefaultDebugModulesPath         = "/debug/modules"
	ThemeLight                      = "light"
	ThemeDark                       = "dark"
)

//go:embed templates/*
var templatesFS embed.FS

// FxCoreModule is the [Fx] core module.
//
// [Fx]: https://github.com/uber-go/fx
var FxCoreModule = fx.Module(
	ModuleName,
	fxgenerate.FxGenerateModule,
	fxconfig.FxConfigModule,
	fxlog.FxLogModule,
	fxtrace.FxTraceModule,
	fxmetrics.FxMetricsModule,
	fxhealthcheck.FxHealthcheckModule,
	fx.Provide(
		NewFxModuleInfoRegistry,
		NewFxCore,
		fx.Annotate(
			NewFxCoreModuleInfo,
			fx.As(new(interface{})),
			fx.ResultTags(`group:"core-module-infos"`),
		),
	),
	fx.Invoke(func(logger *log.Logger, core *Core) {
		logger.Debug().Msg("starting core")
	}),
)

// FxCoreDashboardTheme is the theme for the core dashboard.
type FxCoreDashboardTheme struct {
	Theme string `form:"theme" json:"theme"`
}

// FxCoreParam allows injection of the required dependencies in [NewFxCore].
//
//nolint:containedctx
type FxCoreParam struct {
	fx.In
	Context         context.Context
	LifeCycle       fx.Lifecycle
	Generator       uuid.UuidGenerator
	TracerProvider  oteltrace.TracerProvider
	Checker         *healthcheck.Checker
	Config          *config.Config
	Logger          *log.Logger
	Registry        *FxModuleInfoRegistry
	MetricsRegistry *prometheus.Registry
}

// NewFxCore returns a new [Core].
func NewFxCore(p FxCoreParam) (*Core, error) {
	appDebug := p.Config.AppDebug()

	// logger
	coreLogger := httpserver.NewEchoLogger(
		log.FromZerolog(p.Logger.ToZerolog().With().Str("module", ModuleName).Logger()),
	)

	// server
	coreServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithDebug(appDebug),
		httpserver.WithBanner(false),
		httpserver.WithRecovery(true),
		httpserver.WithLogger(coreLogger),
		httpserver.WithRenderer(NewDashboardRenderer(templatesFS, "templates/dashboard.html")),
		httpserver.WithHttpErrorHandler(
			httpserver.JsonErrorHandler(
				p.Config.GetBool("modules.core.server.errors.obfuscate") || !appDebug,
				p.Config.GetBool("modules.core.server.errors.stack") || appDebug,
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create core http server: %w", err)
	}

	// middlewares
	coreServer = withMiddlewares(coreServer, p)

	// handlers
	coreServer, err = withHandlers(coreServer, p)
	if err != nil {
		return nil, fmt.Errorf("failed to register core http server handlers: %w", err)
	}

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := p.Config.GetString("modules.core.server.address")
			if address == "" {
				address = DefaultAddress
			}

			//nolint:errcheck
			go coreServer.Start(address)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return coreServer.Shutdown(ctx)
		},
	})

	return NewCore(p.Config, p.Checker, coreServer), nil
}

func withMiddlewares(coreServer *echo.Echo, p FxCoreParam) *echo.Echo {
	// CORS middleware
	coreServer.Use(middleware.CORS())

	// request id middleware
	coreServer.Use(httpservermiddleware.RequestIdMiddlewareWithConfig(
		httpservermiddleware.RequestIdMiddlewareConfig{
			Generator: p.Generator,
		},
	))

	// request logger middleware
	requestHeadersToLog := map[string]string{
		httpservermiddleware.HeaderXRequestId: httpservermiddleware.LogFieldRequestId,
	}

	for headerName, fieldName := range p.Config.GetStringMapString("modules.core.server.log.headers") {
		requestHeadersToLog[headerName] = fieldName
	}

	coreServer.Use(httpservermiddleware.RequestLoggerMiddlewareWithConfig(
		httpservermiddleware.RequestLoggerMiddlewareConfig{
			RequestHeadersToLog:             requestHeadersToLog,
			RequestUriPrefixesToExclude:     p.Config.GetStringSlice("modules.core.server.log.exclude"),
			LogLevelFromResponseOrErrorCode: p.Config.GetBool("modules.core.server.log.level_from_response"),
		},
	))

	// request tracer middleware
	if p.Config.GetBool("modules.core.server.trace.enabled") {
		coreServer.Use(httpservermiddleware.RequestTracerMiddlewareWithConfig(
			p.Config.AppName(),
			httpservermiddleware.RequestTracerMiddlewareConfig{
				TracerProvider:              httpserver.AnnotateTracerProvider(p.TracerProvider),
				RequestUriPrefixesToExclude: p.Config.GetStringSlice("modules.core.server.trace.exclude"),
			},
		))
	}

	// request metrics middleware
	if p.Config.GetBool("modules.core.server.metrics.collect.enabled") {
		var buckets []float64
		if bucketsConfig := p.Config.GetString("modules.core.server.metrics.buckets"); bucketsConfig != "" {
			for _, s := range Split(bucketsConfig) {
				f, err := strconv.ParseFloat(s, 64)
				if err == nil {
					buckets = append(buckets, f)
				}
			}
		}

		metricsMiddlewareConfig := httpservermiddleware.RequestMetricsMiddlewareConfig{
			Registry:                p.MetricsRegistry,
			Namespace:               Sanitize(p.Config.GetString("modules.core.server.metrics.collect.namespace")),
			Subsystem:               Sanitize(ModuleName),
			Buckets:                 buckets,
			NormalizeRequestPath:    p.Config.GetBool("modules.core.server.metrics.normalize.request_path"),
			NormalizeResponseStatus: p.Config.GetBool("modules.core.server.metrics.normalize.response_status"),
		}

		coreServer.Use(httpservermiddleware.RequestMetricsMiddlewareWithConfig(metricsMiddlewareConfig))
	}

	return coreServer
}

//nolint:cyclop,gocognit,gocyclo,maintidx
func withHandlers(coreServer *echo.Echo, p FxCoreParam) (*echo.Echo, error) {
	appDebug := p.Config.AppDebug()

	// dashboard
	dashboardEnabled := p.Config.GetBool("modules.core.server.dashboard.enabled")

	// dashboard overview
	overviewInfo, err := p.Registry.Find(ModuleName)
	if err != nil {
		return nil, err
	}

	// dashboard overview expositions
	overviewAppEnvExpose := p.Config.GetBool("modules.core.server.dashboard.overview.app_env")
	overviewAppDebugExpose := p.Config.GetBool("modules.core.server.dashboard.overview.app_debug")
	overviewAppVersionExpose := p.Config.GetBool("modules.core.server.dashboard.overview.app_version")
	overviewLogLevelExpose := p.Config.GetBool("modules.core.server.dashboard.overview.log_level")
	overviewLogOutputExpose := p.Config.GetBool("modules.core.server.dashboard.overview.log_output")
	overviewTraceSamplerExpose := p.Config.GetBool("modules.core.server.dashboard.overview.trace_sampler")
	overviewTraceProcessorExpose := p.Config.GetBool("modules.core.server.dashboard.overview.trace_processor")

	// template expositions
	metricsExpose := p.Config.GetBool("modules.core.server.metrics.expose")
	startupExpose := p.Config.GetBool("modules.core.server.healthcheck.startup.expose")
	livenessExpose := p.Config.GetBool("modules.core.server.healthcheck.liveness.expose")
	readinessExpose := p.Config.GetBool("modules.core.server.healthcheck.readiness.expose")
	configExpose := p.Config.GetBool("modules.core.server.debug.config.expose")
	pprofExpose := p.Config.GetBool("modules.core.server.debug.pprof.expose")
	routesExpose := p.Config.GetBool("modules.core.server.debug.routes.expose")
	statsExpose := p.Config.GetBool("modules.core.server.debug.stats.expose")
	buildExpose := p.Config.GetBool("modules.core.server.debug.build.expose")
	modulesExpose := p.Config.GetBool("modules.core.server.debug.modules.expose")

	// template paths
	metricsPath := p.Config.GetString("modules.core.server.metrics.path")
	startupPath := p.Config.GetString("modules.core.server.healthcheck.startup.path")
	livenessPath := p.Config.GetString("modules.core.server.healthcheck.liveness.path")
	readinessPath := p.Config.GetString("modules.core.server.healthcheck.readiness.path")
	configPath := p.Config.GetString("modules.core.server.debug.config.path")
	pprofPath := p.Config.GetString("modules.core.server.debug.pprof.path")
	routesPath := p.Config.GetString("modules.core.server.debug.routes.path")
	statsPath := p.Config.GetString("modules.core.server.debug.stats.path")
	buildPath := p.Config.GetString("modules.core.server.debug.build.path")
	modulesPath := p.Config.GetString("modules.core.server.debug.modules.path")

	// metrics
	if metricsExpose {
		if metricsPath == "" {
			metricsPath = DefaultMetricsPath
		}

		coreServer.GET(metricsPath, echo.WrapHandler(promhttp.HandlerFor(p.MetricsRegistry, promhttp.HandlerOpts{})))

		coreServer.Logger.Debug("registered metrics handler")
	}

	// healthcheck startup
	if startupExpose {
		if startupPath == "" {
			startupPath = DefaultHealthCheckStartupPath
		}

		coreServer.GET(startupPath, handler.HealthCheckHandler(p.Checker, healthcheck.Startup))

		coreServer.Logger.Debug("registered healthcheck startup handler")
	}

	// healthcheck liveness
	if livenessExpose {
		if livenessPath == "" {
			livenessPath = DefaultHealthCheckLivenessPath
		}

		coreServer.GET(livenessPath, handler.HealthCheckHandler(p.Checker, healthcheck.Liveness))

		coreServer.Logger.Debug("registered healthcheck liveness handler")
	}

	// healthcheck readiness
	if readinessExpose {
		if readinessPath == "" {
			readinessPath = DefaultHealthCheckReadinessPath
		}

		coreServer.GET(readinessPath, handler.HealthCheckHandler(p.Checker, healthcheck.Readiness))

		coreServer.Logger.Debug("registered healthcheck readiness handler")
	}

	// debug config
	if configExpose || appDebug {
		if configPath == "" {
			configPath = DefaultDebugConfigPath
		}

		coreServer.GET(configPath, handler.DebugConfigHandler(p.Config))

		coreServer.Logger.Debug("registered debug config handler")
	}

	// debug pprof
	if pprofExpose || appDebug {
		if pprofPath == "" {
			pprofPath = DefaultDebugPProfPath
		}

		pprofGroup := coreServer.Group(pprofPath)

		pprofGroup.GET("/", handler.PprofIndexHandler())
		pprofGroup.GET("/allocs", handler.PprofAllocsHandler())
		pprofGroup.GET("/block", handler.PprofBlockHandler())
		pprofGroup.GET("/cmdline", handler.PprofCmdlineHandler())
		pprofGroup.GET("/goroutine", handler.PprofGoroutineHandler())
		pprofGroup.GET("/heap", handler.PprofHeapHandler())
		pprofGroup.GET("/mutex", handler.PprofMutexHandler())
		pprofGroup.GET("/profile", handler.PprofProfileHandler())
		pprofGroup.GET("/symbol", handler.PprofSymbolHandler())
		pprofGroup.POST("/symbol", handler.PprofSymbolHandler())
		pprofGroup.GET("/threadcreate", handler.PprofThreadCreateHandler())
		pprofGroup.GET("/trace", handler.PprofTraceHandler())

		coreServer.Logger.Debug("registered debug pprof handlers")
	}

	// debug routes
	if routesExpose || appDebug {
		if routesPath == "" {
			routesPath = DefaultDebugRoutesPath
		}

		coreServer.GET(routesPath, handler.DebugRoutesHandler(coreServer))

		coreServer.Logger.Debug("registered debug routes handler")
	}

	// debug stats
	if statsExpose || appDebug {
		if statsPath == "" {
			statsPath = DefaultDebugStatsPath
		}

		mux := http.NewServeMux()

		err := statsviz.Register(mux, statsviz.Root(statsPath))
		if err != nil {
			coreServer.Logger.Error("failed to register debug stats handler")
		} else {
			statsGroup := coreServer.Group(statsPath)

			statsGroup.GET("/", echo.WrapHandler(mux))
			statsGroup.GET("/*", echo.WrapHandler(mux))

			coreServer.Logger.Debug("registered debug stats handler")
		}
	}

	// debug build
	if buildExpose || appDebug {
		if buildPath == "" {
			buildPath = DefaultDebugBuildPath
		}

		coreServer.GET(buildPath, handler.DebugBuildHandler())

		coreServer.Logger.Debug("registered debug build handler")
	}

	// debug modules
	if modulesExpose || appDebug {
		if modulesPath == "" {
			modulesPath = DefaultDebugModulesPath
		}

		coreServer.GET(fmt.Sprintf("%s/:name", modulesPath), func(c echo.Context) error {
			info, err := p.Registry.Find(c.Param("name"))
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			return c.JSON(http.StatusOK, info.Data())
		})

		coreServer.Logger.Debug("registered debug modules handler")
	}

	// dashboard
	if dashboardEnabled || appDebug {
		// theme
		coreServer.POST("/theme", func(c echo.Context) error {
			themeCookie := new(http.Cookie)
			themeCookie.Name = "theme"

			var theme FxCoreDashboardTheme
			if err = c.Bind(&theme); err != nil {
				themeCookie.Value = ThemeLight
			} else {
				switch theme.Theme {
				case ThemeDark:
					themeCookie.Value = ThemeDark
				case ThemeLight:
					themeCookie.Value = ThemeLight
				default:
					themeCookie.Value = ThemeLight
				}
			}

			c.SetCookie(themeCookie)

			return c.Redirect(http.StatusMovedPermanently, "/")
		})

		coreServer.Logger.Debug("registered dashboard theme handler")

		// render
		coreServer.GET("/", func(c echo.Context) error {
			var theme string
			themeCookie, err := c.Cookie("theme")
			if err == nil {
				switch themeCookie.Value {
				case ThemeDark:
					theme = ThemeDark
				case ThemeLight:
					theme = ThemeLight
				default:
					theme = ThemeLight
				}
			} else {
				theme = ThemeLight
			}

			return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
				"overviewInfo":                 overviewInfo,
				"overviewAppEnvExpose":         overviewAppEnvExpose,
				"overviewAppDebugExpose":       overviewAppDebugExpose,
				"overviewAppVersionExpose":     overviewAppVersionExpose,
				"overviewLogLevelExpose":       overviewLogLevelExpose,
				"overviewLogOutputExpose":      overviewLogOutputExpose,
				"overviewTraceSamplerExpose":   overviewTraceSamplerExpose,
				"overviewTraceProcessorExpose": overviewTraceProcessorExpose,
				"metricsExpose":                metricsExpose,
				"metricsPath":                  metricsPath,
				"startupExpose":                startupExpose,
				"startupPath":                  startupPath,
				"livenessExpose":               livenessExpose,
				"livenessPath":                 livenessPath,
				"readinessExpose":              readinessExpose,
				"readinessPath":                readinessPath,
				"configExpose":                 configExpose || appDebug,
				"configPath":                   configPath,
				"pprofExpose":                  pprofExpose || appDebug,
				"pprofPath":                    pprofPath,
				"routesExpose":                 routesExpose || appDebug,
				"routesPath":                   routesPath,
				"statsExpose":                  statsExpose || appDebug,
				"statsPath":                    statsPath,
				"buildExpose":                  buildExpose || appDebug,
				"buildPath":                    buildPath,
				"modulesExpose":                modulesExpose || appDebug,
				"modulesPath":                  modulesPath,
				"modulesNames":                 p.Registry.Names(),
				"theme":                        theme,
			})
		})

		coreServer.Logger.Debug("registered dashboard handler")
	}

	return coreServer, nil
}
