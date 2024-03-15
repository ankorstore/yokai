package httpserver_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	testRequestId   = "33084b3e-9b90-926c-af19-3859d70bd296"
	testTraceId     = "c4ca71e03e42c2c3d54293a6e2608bfa"
	testSpanId      = "8d0fdc8a74baaaea"
	testTraceParent = fmt.Sprintf("00-%s-%s-01", testTraceId, testSpanId)
)

func TestDefaultHttpServerFactory(t *testing.T) {
	t.Parallel()

	factory := httpserver.NewDefaultHttpServerFactory()

	assert.IsType(t, &httpserver.DefaultHttpServerFactory{}, factory)
	assert.Implements(t, (*httpserver.HttpServerFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	logger, err := log.NewDefaultLoggerFactory().Create()
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)
	binder := &echo.DefaultBinder{}
	jsonSerializer := &echo.DefaultJSONSerializer{}
	httpErrorHandler := func(err error, c echo.Context) {}
	render := httpserver.NewHtmlTemplateRenderer("testdata/templates/*.html")

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithDebug(true),
		httpserver.WithBanner(true),
		httpserver.WithRecovery(true),
		httpserver.WithLogger(echoLogger),
		httpserver.WithBinder(binder),
		httpserver.WithJsonSerializer(jsonSerializer),
		httpserver.WithHttpErrorHandler(httpErrorHandler),
		httpserver.WithRenderer(render),
	)

	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	assert.True(t, httpServer.Debug)
	assert.False(t, httpServer.HideBanner)
	assert.Equal(t, echoLogger, httpServer.Logger)
	assert.Equal(t, binder, httpServer.Binder)
	assert.Equal(t, jsonSerializer, httpServer.JSONSerializer)
	assert.NotNil(t, httpServer.HTTPErrorHandler)
	assert.NotNil(t, httpServer.Renderer)
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOn2xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `ok`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"GET /test",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, testRequestId),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOn4xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusBadRequest, "bad request")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `bad request`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"method":    "GET",
		"uri":       "/test",
		"status":    400,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"GET /test",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, testRequestId),
		semconv.HTTPStatusCode(http.StatusBadRequest),
	)
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOn5xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusInternalServerError, "server error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `server error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"GET /test",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, testRequestId),
		semconv.HTTPStatusCode(http.StatusInternalServerError),
	)
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOn2xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `ok`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"GET /test",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, testRequestId),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOn4xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusBadRequest, "bad request")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `bad request`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "warn",
		"method":    "GET",
		"uri":       "/test",
		"status":    400,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"GET /test",
		attribute.String(httpserver.TraceSpanAttributeHttpRequestId, testRequestId),
		semconv.HTTPStatusCode(http.StatusBadRequest),
	)
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOn5xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusInternalServerError, "server error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `server error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnHttpError2xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusOK, "http error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"http error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=200, message=http error",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=200, message=http error",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnHttpError4xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "http bad request")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"http bad request"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=400, message=http bad request",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=400, message=http bad request",
		"method":    "GET",
		"uri":       "/test",
		"status":    400,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnHttpError5xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "http custom error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"http custom error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=http custom error",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=500, message=http custom error",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnHttpError5xxWithStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, true)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "http custom error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"http custom error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=500, message=http custom error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=http custom error",
		"message":   "error handler",
		"stack":     "*echo.HTTPError code=500, message=http custom error",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=500, message=http custom error",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnHttpError5xxWithObfuscateAndStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(true, true)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "http custom error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"Internal Server Error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=500, message=http custom error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=http custom error",
		"message":   "error handler",
		"stack":     "*echo.HTTPError code=500, message=http custom error",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=500, message=http custom error",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnComplexHttpError2xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusOK, echo.Map{"some": "data"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"some":"data"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=200, message=map[some:data]",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=200, message=map[some:data]",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnComplexHttpError4xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"some": "data"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"some":"data"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=400, message=map[some:data]",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=400, message=map[some:data]",
		"method":    "GET",
		"uri":       "/test",
		"status":    400,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnComplexHttpError5xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"some": "data"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"some":"data"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=map[some:data]",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=500, message=map[some:data]",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnComplexHttpError5xxWithStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, true)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"some": "data"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"some":"data"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=500, message=map[some:data]`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=map[some:data]",
		"stack":     "*echo.HTTPError code=500, message=map[some:data]",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=500, message=map[some:data]",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnComplexHttpError5xxWithObfuscateAndStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(true, true)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"some": "data"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"Internal Server Error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=500, message=map[some:data]`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=map[some:data]",
		"stack":     "*echo.HTTPError code=500, message=map[some:data]",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=500, message=map[some:data]",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnHttpError2xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusOK, "http error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"http error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=200, message=http error",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=200, message=http error",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnHttpError4xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "http bad request")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"http bad request"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=400, message=http bad request",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "warn",
		"error":     "code=400, message=http bad request",
		"method":    "GET",
		"uri":       "/test",
		"status":    400,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnHttpError5xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "http custom error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"http custom error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=http custom error",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=http custom error",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnComplexHttpError2xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusOK, echo.Map{"some": "data"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"some":"data"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=200, message=map[some:data]",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"error":     "code=200, message=map[some:data]",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnComplexHttpError4xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"some": "data"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"some":"data"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=400, message=map[some:data]",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "warn",
		"error":     "code=400, message=map[some:data]",
		"method":    "GET",
		"uri":       "/test",
		"status":    400,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnComplexHttpError5xx(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"some": "data"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"some":"data"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=map[some:data]",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "code=500, message=map[some:data]",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnGenericError(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("generic error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"generic error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "generic error",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"uri":       "/test",
		"status":    500,
		"error":     "generic error",
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnGenericError(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, false)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("generic error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"generic error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "generic error",
		"message":   "error handler",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"uri":       "/test",
		"status":    500,
		"error":     "generic error",
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnGenericErrorWithStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, true)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("wrapped error: %w", fmt.Errorf("generic error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"wrapped error: generic error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*fmt.wrapError wrapped error: generic error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "wrapped error: generic error",
		"message":   "error handler",
		"stack":     "*fmt.wrapError wrapped error: generic error",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"uri":       "/test",
		"status":    500,
		"error":     "wrapped error: generic error",
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestLoggerAndTracerAndErrorHandlerOnGenericErrorWithObfuscateAndStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(true, true)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("wrapped error: %w", fmt.Errorf("generic error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"Internal Server Error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*fmt.wrapError wrapped error: generic error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "wrapped error: generic error",
		"message":   "error handler",
		"stack":     "*fmt.wrapError wrapped error: generic error",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"uri":       "/test",
		"status":    500,
		"error":     "wrapped error: generic error",
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnGenericErrorWithStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(false, true)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("wrapped error: %w", fmt.Errorf("generic error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"wrapped error: generic error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*fmt.wrapError wrapped error: generic error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "wrapped error: generic error",
		"message":   "error handler",
		"stack":     "*fmt.wrapError wrapped error: generic error",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"uri":       "/test",
		"status":    500,
		"error":     "wrapped error: generic error",
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithLeveledRequestLoggerAndTracerAndErrorHandlerOnGenericErrorWithObfuscateAndStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	exporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithHttpErrorHandler(httpserver.JsonErrorHandler(true, true)),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestTracerMiddlewareWithConfig(
		"test",
		middleware.RequestTracerMiddlewareConfig{
			TracerProvider: tracerProvider,
		},
	))
	httpServer.Use(middleware.RequestLoggerMiddlewareWithConfig(
		middleware.RequestLoggerMiddlewareConfig{
			LogLevelFromResponseOrErrorCode: true,
		},
	))

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("wrapped error: %w", fmt.Errorf("generic error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, testRequestId)
	req.Header.Add(middleware.HeaderTraceParent, testTraceParent)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"Internal Server Error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*fmt.wrapError wrapped error: generic error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"error":     "wrapped error: generic error",
		"message":   "error handler",
		"stack":     "*fmt.wrapError wrapped error: generic error",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"uri":       "/test",
		"status":    500,
		"error":     "wrapped error: generic error",
		"message":   "request logger",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})
}

func TestCreateWithRequestMetricsAndWithNormalization(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create()
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestMetricsMiddlewareWithConfig(middleware.RequestMetricsMiddlewareConfig{
		Registry:                registry,
		Namespace:               "namespace",
		Subsystem:               "subsystem",
		NormalizeResponseStatus: true,
		NormalizeRequestPath:    true,
	}))

	httpServer.GET("/foo/bar/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	req := httptest.NewRequest(http.MethodGet, "/foo/bar/baz?page=1", nil)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `baz`)

	expectedCounterMetric := `
		# HELP namespace_subsystem_http_server_requests_total Number of processed HTTP requests
        # TYPE namespace_subsystem_http_server_requests_total counter
        namespace_subsystem_http_server_requests_total{method="GET",path="/foo/bar/:id",status="2xx"} 1
	`

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedCounterMetric),
		"namespace_subsystem_http_server_requests_total",
	)
	assert.NoError(t, err)
}

func TestCreateWithRequestMetricsAndWithoutNormalization(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create()
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestMetricsMiddlewareWithConfig(middleware.RequestMetricsMiddlewareConfig{
		Registry:                registry,
		Namespace:               "namespace",
		Subsystem:               "subsystem",
		NormalizeResponseStatus: false,
		NormalizeRequestPath:    false,
	}))

	httpServer.GET("/foo/bar/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	req := httptest.NewRequest(http.MethodGet, "/foo/bar/baz?page=1", nil)

	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `baz`)

	expectedCounterMetric := `
		# HELP namespace_subsystem_http_server_requests_total Number of processed HTTP requests
        # TYPE namespace_subsystem_http_server_requests_total counter
        namespace_subsystem_http_server_requests_total{method="GET",path="/foo/bar/baz?page=1",status="200"} 1
	`

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedCounterMetric),
		"namespace_subsystem_http_server_requests_total",
	)
	assert.NoError(t, err)
}

func TestCreateWithPanicRecovery(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer, err := httpserver.NewDefaultHttpServerFactory().Create(
		httpserver.WithLogger(httpserver.NewEchoLogger(logger)),
		httpserver.WithRecovery(true),
	)
	assert.NoError(t, err)
	assert.IsType(t, &echo.Echo{}, httpServer)

	httpServer.Use(middleware.RequestLoggerMiddleware())

	httpServer.GET("/panic", func(c echo.Context) error {
		panic("custom panic")
	})

	defer func() {
		if r := recover(); r != nil {
			t.Error("should have recovered by itself")
		}
	}()

	for i := 0; i <= 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), `{"message":"Internal Server Error"}`)
	}
}
