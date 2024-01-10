package middleware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/labstack/echo/v4"
	gommonlog "github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestRequestLoggerMiddlewareWithDefaults(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")

	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		// echo logger
		c.Logger().Info("test-echo-logger")

		// zero logger
		httpserver.CtxLogger(c).Info().Msg("test-zero-logger")

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestLoggerMiddleware()
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test-echo-logger",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test-zero-logger",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithSkipper(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		// from echo logger
		c.Logger().Info("test-echo-logger")

		// from zero logger
		httpserver.CtxLogger(c).Info().Msg("test zero logger")

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		Skipper: func(echo.Context) bool {
			return true
		},
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	hasRecord, err := logBuffer.HasRecord(map[string]interface{}{
		"level":        "info",
		"method":       "GET",
		"uri":          "/test",
		"status":       200,
		"message":      "request",
		"x-request-id": "test-request-id",
	})
	assert.NoError(t, err)
	assert.False(t, hasRecord)
}

func TestRequestLoggerMiddlewareWithCustomRequestHeadersToLog(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add("x-custom-header1", "value-header1")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		RequestHeadersToLog: map[string]string{
			"x-custom-header1": "custom-header1",
			"x-custom-header2": "custom-header1",
		},
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":          "info",
		"message":        "test",
		"custom-header1": "value-header1",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":          "info",
		"method":         "GET",
		"uri":            "/test",
		"status":         200,
		"message":        "request logger",
		"custom-header1": "value-header1",
	})
}

func TestRequestLoggerMiddlewareWithCustomRequestUriToExclude(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		RequestUriPrefixesToExclude: []string{
			"/test",
		},
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test",
	})

	hasRecord, err := logBuffer.HasRecord(map[string]interface{}{
		"level":   "info",
		"method":  "GET",
		"uri":     "/test",
		"status":  200,
		"message": "request logger",
	})
	assert.NoError(t, err)
	assert.False(t, hasRecord)
}

func TestRequestLoggerMiddlewareWithCustomRequestUriToExcludeWithResponseError(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return c.String(http.StatusInternalServerError, "custom error")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		RequestUriPrefixesToExclude: []string{
			"/test",
		},
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"method":  "GET",
		"uri":     "/test",
		"status":  500,
		"message": "request logger",
	})
}

func TestRequestLoggerMiddlewareWithCustomRequestUriToExcludeWithHttpError(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return echo.NewHTTPError(http.StatusInternalServerError, "custom error")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		RequestUriPrefixesToExclude: []string{
			"/test",
		},
	})
	h := m(handler)

	err = h(ctx)
	assert.Error(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"method":  "GET",
		"uri":     "/test",
		"status":  500,
		"error":   "code=500, message=custom error",
		"message": "request logger",
	})
}

func TestRequestLoggerMiddlewareWithCustomRequestUriToExcludeWithGenericError(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return fmt.Errorf("generic error")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		RequestUriPrefixesToExclude: []string{
			"/test",
		},
	})
	h := m(handler)

	err = h(ctx)
	assert.Error(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"method":  "GET",
		"uri":     "/test",
		"status":  500,
		"error":   "generic error",
		"message": "request logger",
	})
}

func TestRequestLoggerMiddlewareWithFailingHandler(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "custom error")
	}

	m := middleware.RequestLoggerMiddleware()
	h := m(handler)

	err = h(ctx)
	assert.Error(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"error":     "code=500, message=custom error",
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithLogLevelFromResponseOnHttpResponse2xx(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return c.String(http.StatusOK, "data")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		LogLevelFromResponseOrErrorCode: true,
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithLogLevelFromResponseOnHttpResponse4xx(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return c.String(http.StatusBadRequest, "bad request")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		LogLevelFromResponseOrErrorCode: true,
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "warn",
		"method":    "GET",
		"uri":       "/test",
		"status":    400,
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithLogLevelFromResponseOnHttpResponse5xx(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return c.String(http.StatusInternalServerError, "custom error")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		LogLevelFromResponseOrErrorCode: true,
	})
	h := m(handler)

	err = h(ctx)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithLogLevelFromResponseOnHttpError2xx(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return echo.NewHTTPError(http.StatusOK, "custom error")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		LogLevelFromResponseOrErrorCode: true,
	})
	h := m(handler)

	err = h(ctx)
	assert.Error(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"method":    "GET",
		"uri":       "/test",
		"status":    200,
		"error":     "code=200, message=custom error",
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithLogLevelFromResponseOnHttpError4xx(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return echo.NewHTTPError(http.StatusBadRequest, "http bad request")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		LogLevelFromResponseOrErrorCode: true,
	})
	h := m(handler)

	err = h(ctx)
	assert.Error(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "warn",
		"method":    "GET",
		"uri":       "/test",
		"status":    400,
		"error":     "code=400, message=http bad request",
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithLogLevelFromResponseOnHttpError5xx(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return echo.NewHTTPError(http.StatusInternalServerError, "http error")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		LogLevelFromResponseOrErrorCode: true,
	})
	h := m(handler)

	err = h(ctx)
	assert.Error(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"error":     "code=500, message=http error",
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithLogLevelFromResponseOnGenericError(t *testing.T) {
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(middleware.HeaderXRequestId, "test-request-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		c.Logger().Info("test")

		return fmt.Errorf("generic error")
	}

	m := middleware.RequestLoggerMiddlewareWithConfig(middleware.RequestLoggerMiddlewareConfig{
		LogLevelFromResponseOrErrorCode: true,
	})
	h := m(handler)

	err = h(ctx)
	assert.Error(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "test",
		"requestID": "test-request-id",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"method":    "GET",
		"uri":       "/test",
		"status":    500,
		"error":     "generic error",
		"message":   "request logger",
		"requestID": "test-request-id",
	})
}

func TestRequestLoggerMiddlewareWithoutContextLogger(t *testing.T) {
	buffer := new(bytes.Buffer)

	logger := gommonlog.New("echo")
	logger.SetOutput(buffer)
	logger.DisableColor()

	httpServer := echo.New()
	httpServer.Logger = logger

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestLoggerMiddleware()
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)
}
