package httpserver

import (
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/labstack/echo/v4"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TracerName is the httpserver tracer name.
const TracerName = "httpserver"

// CtxRequestIdKey is a contextual struct key.
type CtxRequestIdKey struct{}

// CtxRequestId returns the contextual request id.
func CtxRequestId(c echo.Context) string {
	if rid, ok := c.Request().Context().Value(CtxRequestIdKey{}).(string); ok {
		return rid
	} else {
		return ""
	}
}

// CtxRequestId returns the contextual [log.Logger].
func CtxLogger(c echo.Context) *log.Logger {
	return log.CtxLogger(c.Request().Context())
}

// CtxTracer returns the contextual [Tracer].
//
// [Tracer]: https://go.opentelemetry.io/otel/trace
func CtxTracer(c echo.Context) oteltrace.Tracer {
	return trace.CtxTracerProvider(c.Request().Context()).Tracer(TracerName)
}
