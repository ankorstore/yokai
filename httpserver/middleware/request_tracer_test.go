package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestRequestTracerMiddlewareWithDefaults(t *testing.T) {
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
		_, span := trace.CtxTracerProvider(c.Request().Context()).Tracer("test").Start(c.Request().Context(), "test span")
		defer span.End()

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestTracerMiddleware("test")
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	tracetest.AssertHasTraceSpan(t, exporter, "GET /test")
	tracetest.AssertHasTraceSpan(t, exporter, "test span")
}

func TestRequestTracerMiddlewareWithOptions(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(false),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		_, span := trace.CtxTracerProvider(c.Request().Context()).Tracer("test").Start(c.Request().Context(), "test span")
		defer span.End()

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestTracerMiddlewareWithConfig("test", middleware.RequestTracerMiddlewareConfig{
		TracerProvider: tracerProvider,
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"POST /test",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, "test-request-id"),
	)
	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"test span",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, "test-request-id"),
	)
}

func TestRequestTracerMiddlewareWithSkipper(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		_, span := trace.CtxTracerProvider(c.Request().Context()).Tracer("test").Start(c.Request().Context(), "test span")
		defer span.End()

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestTracerMiddlewareWithConfig("test", middleware.RequestTracerMiddlewareConfig{
		Skipper: func(echo.Context) bool {
			return true
		},
		TracerProvider: tracerProvider,
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	tracetest.AssertHasNotTraceSpan(t, exporter, "DELETE /test")
	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"test span",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, "test-request-id"),
	)
}

func TestRequestTracerMiddlewareWithCustomRequestUriToExclude(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		_, span := trace.CtxTracerProvider(c.Request().Context()).Tracer("test").Start(c.Request().Context(), "test span")
		defer span.End()

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestTracerMiddlewareWithConfig("test", middleware.RequestTracerMiddlewareConfig{
		RequestUriPrefixesToExclude: []string{
			"/test",
		},
		TracerProvider: tracerProvider,
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	tracetest.AssertHasNotTraceSpan(t, exporter, "PUT /test")
	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"test span",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, "test-request-id"),
	)
}

func TestRequestTracerMiddlewareWithFailingHandler(t *testing.T) {
	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(false),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		_, span := trace.CtxTracerProvider(c.Request().Context()).Tracer("test").Start(c.Request().Context(), "test span")
		defer span.End()

		return echo.NewHTTPError(http.StatusInternalServerError, "custom error")
	}

	m := middleware.RequestTracerMiddlewareWithConfig("test", middleware.RequestTracerMiddlewareConfig{
		TracerProvider: tracerProvider,
	})
	h := m(handler)

	err = h(ctx)
	assert.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"PATCH /test",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, "test-request-id"),
		attribute.String("handler.error", "code=500, message=custom error"),
	)
	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"test span",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, "test-request-id"),
	)
}
