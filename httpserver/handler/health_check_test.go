package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/httpserver/handler"
	"github.com/ankorstore/yokai/httpserver/testdata/probes"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	checker, err := healthcheck.NewDefaultCheckerFactory().Create(
		healthcheck.WithProbe(probes.NewSuccessProbe()),
		healthcheck.WithProbe(probes.NewFailureProbe(), healthcheck.Liveness, healthcheck.Readiness),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Logger = httpserver.NewEchoLogger(logger)

	// [GET] /healthz => startup probes (should not log success)
	httpServer.GET("/healthz", handler.HealthCheckHandler(checker, healthcheck.Startup))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		rec.Body.String(),
		`{"success":true,"probes":{"successProbe":{"success":true,"message":"some success"}}}`,
	)

	logBufferRecords, err := logBuffer.Records()
	assert.NoError(t, err)
	assert.Len(t, logBufferRecords, 0)

	// [GET] /livez => liveness probes (should log failure)
	logBuffer.Reset()

	httpServer.GET("/livez", handler.HealthCheckHandler(checker, healthcheck.Liveness))

	req = httptest.NewRequest(http.MethodGet, "/livez", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec = httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(
		t,
		rec.Body.String(),
		`{"success":false,"probes":{"failureProbe":{"success":false,"message":"some failure"},"successProbe":{"success":true,"message":"some success"}}}`,
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "error",
		"successProbe": "success: true, message: some success",
		"failureProbe": "success: false, message: some failure",
		"message":      "healthcheck failure",
	})

	// [GET] /readyz => readiness probes (should log failure)
	logBuffer.Reset()

	httpServer.GET("/readyz", handler.HealthCheckHandler(checker, healthcheck.Readiness))

	req = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	req = req.WithContext(logger.WithContext(context.Background()))
	rec = httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(
		t,
		rec.Body.String(),
		`{"success":false,"probes":{"failureProbe":{"success":false,"message":"some failure"},"successProbe":{"success":true,"message":"some success"}}}`,
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "error",
		"successProbe": "success: true, message: some success",
		"failureProbe": "success: false, message: some failure",
		"message":      "healthcheck failure",
	})
}
