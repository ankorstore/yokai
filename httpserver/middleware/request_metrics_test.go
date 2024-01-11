package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRequestMetricsMiddlewareWithDefaults(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		time.Sleep(1 * time.Millisecond)

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestMetricsMiddleware()
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())

	// requests counter assertions
	expectedCounterMetric := `
		# HELP requests_total Number of processed HTTP requests
		# TYPE requests_total counter
        requests_total{handler="/not-found",method="GET",status="2xx"} 1
	`

	err = testutil.GatherAndCompare(
		prometheus.DefaultGatherer,
		strings.NewReader(expectedCounterMetric),
		"requests_total",
	)
	assert.NoError(t, err)
}

func TestRequestMetricsMiddlewareWithSkipper(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		time.Sleep(1 * time.Millisecond)

		return c.String(http.StatusOK, "ok")
	}

	m := middleware.RequestMetricsMiddlewareWithConfig(middleware.RequestMetricsMiddlewareConfig{
		Registry: registry,
		Skipper: func(echo.Context) bool {
			return true
		},
		NormalizeHTTPStatus: false,
	})
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())

	// requests counter assertions
	expectedCounterMetric := `
		# HELP requests_total Number of processed HTTP requests
		# TYPE requests_total counter
        requests_total{handler="/not-found",method="GET",status="200"} 1
	`

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedCounterMetric),
		"requests_total",
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "metric output does not match expectation")
}

func TestRequestMetricsMiddlewareWithCustomOptions(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		time.Sleep(1 * time.Millisecond)

		return c.String(http.StatusNotFound, "not found")
	}

	m := middleware.RequestMetricsMiddlewareWithConfig(middleware.RequestMetricsMiddlewareConfig{
		Registry:            registry,
		Namespace:           "namespace",
		Subsystem:           "subsystem",
		Buckets:             []float64{0.01, 1, 10},
		NormalizeHTTPStatus: true,
	})
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "not found", rec.Body.String())

	// requests counter assertions
	expectedCounterMetric := `
		# HELP namespace_subsystem_requests_total Number of processed HTTP requests
        # TYPE namespace_subsystem_requests_total counter
        namespace_subsystem_requests_total{handler="/not-found",method="GET",status="4xx"} 1
	`

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedCounterMetric),
		"namespace_subsystem_requests_total",
	)
	assert.NoError(t, err)
}
