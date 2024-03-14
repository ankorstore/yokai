package middleware

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/ankorstore/yokai/httpserver/normalization"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	HttpServerMetricsRequestsCount    = "httpserver_requests_total"
	HttpServerMetricsRequestsDuration = "httpserver_requests_duration_seconds"
	HttpServerMetricsNotFoundPath     = "/not-found"
)

// RequestMetricsMiddlewareConfig is the configuration for the [RequestMetricsMiddleware].
type RequestMetricsMiddlewareConfig struct {
	Skipper                 middleware.Skipper
	Registry                prometheus.Registerer
	Namespace               string
	Buckets                 []float64
	Subsystem               string
	NormalizeRequestPath    bool
	NormalizeResponseStatus bool
}

// DefaultRequestMetricsMiddlewareConfig is the default configuration for the [RequestMetricsMiddleware].
var DefaultRequestMetricsMiddlewareConfig = RequestMetricsMiddlewareConfig{
	Skipper:                 middleware.DefaultSkipper,
	Registry:                prometheus.DefaultRegisterer,
	Namespace:               "",
	Subsystem:               "",
	Buckets:                 prometheus.DefBuckets,
	NormalizeRequestPath:    true,
	NormalizeResponseStatus: true,
}

// RequestMetricsMiddleware returns a [RequestMetricsMiddleware] with the [DefaultRequestMetricsMiddlewareConfig].
func RequestMetricsMiddleware() echo.MiddlewareFunc {
	return RequestMetricsMiddlewareWithConfig(DefaultRequestMetricsMiddlewareConfig)
}

// RequestMetricsMiddlewareWithConfig returns a [RequestMetricsMiddleware] for a provided [RequestMetricsMiddlewareConfig].
func RequestMetricsMiddlewareWithConfig(config RequestMetricsMiddlewareConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRequestMetricsMiddlewareConfig.Skipper
	}

	if config.Registry == nil {
		config.Registry = DefaultRequestMetricsMiddlewareConfig.Registry
	}

	if config.Namespace == "" {
		config.Namespace = DefaultRequestMetricsMiddlewareConfig.Namespace
	}

	if config.Subsystem == "" {
		config.Subsystem = DefaultRequestMetricsMiddlewareConfig.Subsystem
	}

	if len(config.Buckets) == 0 {
		config.Buckets = DefaultRequestMetricsMiddlewareConfig.Buckets
	}

	httpRequestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      HttpServerMetricsRequestsCount,
			Help:      "Number of processed HTTP requests",
		},
		[]string{
			"status",
			"method",
			"path",
		},
	)

	httpRequestsDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      HttpServerMetricsRequestsDuration,
			Help:      "Time spent processing HTTP requests",
			Buckets:   config.Buckets,
		},
		[]string{
			"method",
			"path",
		},
	)

	config.Registry.MustRegister(httpRequestsCounter, httpRequestsDuration)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// skipper
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()

			var path string
			if config.NormalizeRequestPath {
				path = c.Path()
			} else {
				path = req.URL.Path
				if req.URL.RawQuery != "" {
					path = fmt.Sprintf("%s?%s", path, req.URL.RawQuery)
				}
			}

			// to avoid high cardinality on 404s
			if reflect.ValueOf(c.Handler()).Pointer() == reflect.ValueOf(echo.NotFoundHandler).Pointer() {
				path = HttpServerMetricsNotFoundPath
			}

			timer := prometheus.NewTimer(httpRequestsDuration.WithLabelValues(req.Method, path))
			err := next(c)
			timer.ObserveDuration()

			if err != nil {
				c.Error(err)
			}

			status := ""
			if config.NormalizeResponseStatus {
				status = normalization.NormalizeStatus(c.Response().Status)
			} else {
				status = strconv.Itoa(c.Response().Status)
			}

			httpRequestsCounter.WithLabelValues(status, req.Method, path).Inc()

			return err
		}
	}
}
