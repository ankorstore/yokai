package httpserver_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEmptyCtxRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	assert.Equal(t, "", httpserver.CtxRequestId(c))
}

func TestCtxRequestIdWithValue(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithServiceName("test service"),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	// with value
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		httpserver.CtxLogger(c).Info().Msgf("request id: %s", httpserver.CtxRequestId(c))

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestLoggerMiddleware()
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test service",
		"message": fmt.Sprintf("request id: %s", testRequestId),
	})
}

func TestCtxLogger(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithServiceName("test service"),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		httpserver.CtxLogger(c).Info().Msg("test message")

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestLoggerMiddleware()
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test service",
		"message": "test message",
	})
}

func TestCtxTracer(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(false),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		tracer := httpserver.CtxTracer(c)

		_, span := tracer.Start(c.Request().Context(), "test span")
		span.End()

		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	tracetest.AssertHasTraceSpan(t, exporter, "GET /test")
	tracetest.AssertHasTraceSpan(t, exporter, "test span")
}

func TestCtxTracerFromGlobalsAndWithMiddleware(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	_, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Use(middleware.RequestTracerMiddleware("test"))

	httpServer.GET("/test", func(c echo.Context) error {
		tracer := httpserver.CtxTracer(c)

		_, span := tracer.Start(c.Request().Context(), "test span")
		span.End()

		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	tracetest.AssertHasTraceSpan(t, exporter, "GET /test")
	tracetest.AssertHasTraceSpan(t, exporter, "test span")
}

func TestCtxTracerFromGlobalsAndWithoutMiddleware(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	_, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		tracer := httpserver.CtxTracer(c)

		_, span := tracer.Start(c.Request().Context(), "test span")
		span.End()

		return c.String(http.StatusOK, "ok")
	}

	err = handler(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	assert.False(t, exporter.HasSpan("GET /test"))
	tracetest.AssertHasTraceSpan(t, exporter, "test span")
}
