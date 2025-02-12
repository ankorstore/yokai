package httpserver_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandling(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, false).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("custom error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"custom error"}`)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "custom error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithObfuscate(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(true, false).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("custom error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"Internal Server Error"}`)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "custom error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, true).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("custom error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"custom error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*errors.errorString custom error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "custom error",
		"stack":   "*errors.errorString custom error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithObfuscateAndStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(true, true).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return fmt.Errorf("custom error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"Internal Server Error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*errors.errorString custom error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "custom error",
		"stack":   "*errors.errorString custom error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithHeadRequest(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, false).Handle()

	httpServer.HEAD("/test", func(c echo.Context) error {
		return fmt.Errorf("custom error")
	})

	req := httptest.NewRequest(http.MethodHead, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Empty(t, rec.Body.String())

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "custom error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithAlreadyCommittedResponse(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, false).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		err := fmt.Errorf("custom error")
		c.Error(err)

		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"custom error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "custom error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithHttpError(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, false).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"bad request error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "code=400, message=bad request error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithWrappedError(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, false).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("wrapped error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"wrapped error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "code=400, message=wrapped error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithWrappedErrorWithStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, true).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("wrapped error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"wrapped error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=400, message=wrapped error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "code=400, message=wrapped error",
		"stack":   "*echo.HTTPError code=400, message=wrapped error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithWrappedErrorWithObfuscateAndStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(true, true).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("wrapped error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"Bad Request"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=400, message=wrapped error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "code=400, message=wrapped error",
		"stack":   "*echo.HTTPError code=400, message=wrapped error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithHttpErrorWithStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, true).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"bad request error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=400, message=bad request error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "bad request error",
		"stack":   "*echo.HTTPError code=400, message=bad request error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithHttpErrorWithObfuscateAndStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(true, true).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"Bad Request"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=400, message=bad request error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "bad request error",
		"stack":   "*echo.HTTPError code=400, message=bad request error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithInternalHttpError(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, false).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		internalError := echo.NewHTTPError(http.StatusInternalServerError, "internal error")

		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  internalError.Error(),
			Internal: internalError,
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `{"message":"internal error"}`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "internal error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithInternalHttpErrorWithStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(false, true).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		internalError := echo.NewHTTPError(http.StatusInternalServerError, "internal error")

		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  internalError.Error(),
			Internal: internalError,
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"internal error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=500, message=code=500, message=internal error, internal=code=500, message=internal error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "internal error",
		"stack":   "*echo.HTTPError code=500, message=code=500, message=internal error, internal=code=500, message=internal error",
		"message": "error handler",
	})
}

func TestErrorHandlingWithInternalHttpErrorWithObfuscateAndStack(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)
	httpServer.HTTPErrorHandler = httpserver.NewJsonErrorHandler(true, true).Handle()

	httpServer.GET("/test", func(c echo.Context) error {
		internalError := echo.NewHTTPError(http.StatusInternalServerError, "internal error")

		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  internalError.Error(),
			Internal: internalError,
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"Internal Server Error"`)
	assert.Contains(t, rec.Body.String(), `"stack":"*echo.HTTPError code=500, message=code=500, message=internal error, internal=code=500, message=internal error`)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "internal error",
		"stack":   "*echo.HTTPError code=500, message=code=500, message=internal error, internal=code=500, message=internal error",
		"message": "error handler",
	})
}
