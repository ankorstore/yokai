package middleware

import (
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	HttpServerMetricsRequestsCount    = "requests_total"
	HttpServerMetricsRequestsDuration = "request_duration_seconds"
	HttpServerMetricsNotFoundPath     = "/not-found"
)

// RequestMetricsMiddlewareConfig is the configuration for the [RequestMetricsMiddleware].
type RequestMetricsMiddlewareConfig struct {
	Skipper             middleware.Skipper
	Registry            prometheus.Registerer
	Namespace           string
	Buckets             []float64
	Subsystem           string
	NormalizeHTTPStatus bool
}

// DefaultRequestMetricsMiddlewareConfig is the default configuration for the [RequestMetricsMiddleware].
var DefaultRequestMetricsMiddlewareConfig = RequestMetricsMiddlewareConfig{
	Skipper:             middleware.DefaultSkipper,
	Registry:            prometheus.DefaultRegisterer,
	Namespace:           "",
	Subsystem:           "",
	Buckets:             prometheus.DefBuckets,
	NormalizeHTTPStatus: true,
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
			"handler",
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
			"handler",
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
			path := c.Path()

			// to avoid high cardinality
			if isNotFoundHandler(c.Handler()) {
				path = HttpServerMetricsNotFoundPath
			}

			timer := prometheus.NewTimer(httpRequestsDuration.WithLabelValues(req.Method, path))
			err := next(c)
			timer.ObserveDuration()

			if err != nil {
				c.Error(err)
			}

			status := ""
			if config.NormalizeHTTPStatus {
				status = normalizeHTTPStatus(c.Response().Status)
			} else {
				status = strconv.Itoa(c.Response().Status)
			}

			httpRequestsCounter.WithLabelValues(status, req.Method, path).Inc()

			return err
		}
	}
}

func normalizeHTTPStatus(status int) string {
	switch {
	case status < 200:
		return "1xx"
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	default:
		return "5xx"
	}
}

func isNotFoundHandler(handler echo.HandlerFunc) bool {
	return reflect.ValueOf(handler).Pointer() == reflect.ValueOf(echo.NotFoundHandler).Pointer()
}
