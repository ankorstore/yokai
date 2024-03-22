package middleware

import (
	"fmt"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/trace"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// RequestTracerMiddlewareConfig is the configuration for the [RequestTracerMiddleware].
type RequestTracerMiddlewareConfig struct {
	Skipper                     middleware.Skipper
	TracerProvider              oteltrace.TracerProvider
	TextMapPropagator           propagation.TextMapPropagator
	RequestUriPrefixesToExclude []string
}

// DefaultRequestTracerMiddlewareConfig is the default configuration for the [RequestTracerMiddleware].
var DefaultRequestTracerMiddlewareConfig = RequestTracerMiddlewareConfig{
	Skipper:                     middleware.DefaultSkipper,
	TracerProvider:              otel.GetTracerProvider(),
	TextMapPropagator:           propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	RequestUriPrefixesToExclude: []string{},
}

// RequestTracerMiddleware returns a [RequestTracerMiddleware] with the [DefaultRequestTracerMiddlewareConfig].
func RequestTracerMiddleware(serviceName string) echo.MiddlewareFunc {
	return RequestTracerMiddlewareWithConfig(serviceName, DefaultRequestTracerMiddlewareConfig)
}

// RequestTracerMiddlewareWithConfig returns a [RequestTracerMiddleware] for a provided [RequestTracerMiddlewareConfig].
func RequestTracerMiddlewareWithConfig(serviceName string, config RequestTracerMiddlewareConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRequestTracerMiddlewareConfig.Skipper
	}

	if config.TracerProvider == nil {
		config.TracerProvider = DefaultRequestTracerMiddlewareConfig.TracerProvider
	}

	if config.TextMapPropagator == nil {
		config.TextMapPropagator = DefaultRequestTracerMiddlewareConfig.TextMapPropagator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// req
			request := c.Request()

			// header carrier context propagation
			ctx := config.TextMapPropagator.Extract(request.Context(), propagation.HeaderCarrier(request.Header))

			// tracer provider context propagation
			ctx = trace.WithContext(ctx, config.TracerProvider)

			c.SetRequest(request.WithContext(ctx))

			// skip
			if config.Skipper(c) || httpserver.MatchPrefix(config.RequestUriPrefixesToExclude, request.URL.Path) {
				return next(c)
			}

			// request tracing preparation
			spanOptions := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(semconv.HTTPRoute(request.URL.Path)),
				oteltrace.WithAttributes(httpconv.ServerRequest(serviceName, request)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}

			path := c.Path()
			if path == "" {
				path = request.URL.Path
			}

			spanName := fmt.Sprintf("%s %s", request.Method, path)

			ctx, span := config.TracerProvider.Tracer(serviceName).Start(ctx, spanName, spanOptions...)
			defer span.End()

			c.SetRequest(request.WithContext(ctx))

			// call next in chain
			err := next(c)
			if err != nil {
				span.SetAttributes(attribute.String("handler.error", err.Error()))
				c.Error(err)
			}

			// response span annotation
			status := c.Response().Status
			span.SetStatus(httpconv.ServerStatus(status))

			if status > 0 {
				span.SetAttributes(semconv.HTTPStatusCode(status))
			}

			return err
		}
	}
}
