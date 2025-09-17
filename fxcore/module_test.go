package fxcore_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxcore/testdata/probes"
	"github.com/ankorstore/yokai/fxcore/testdata/tasks"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/fx"
)

func TestModuleWithServerDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("MODULES_CORE_SERVER_EXPOSE", "false")

	var core *fxcore.Core

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core))

	assert.Nil(t, core.HttpServer())
}

func TestModuleWithMetricsDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("METRICS_ENABLED", "false")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /metrics
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/metrics",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /metrics",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/metrics"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleWithMetricsEnabledAndCollected(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter, &metricsRegistry))

	// [GET] / twice to generate some metrics
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	core.HttpServer().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()

	core.HttpServer().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// [GET] /metrics to check the metrics
	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// assert [GET] / counter = 2
	assert.Contains(t, rec.Body.String(), `core_http_server_requests_duration_seconds_bucket{method="GET",path="/",le="1"} 2`)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"uri":     "/metrics",
		"status":  200,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /metrics",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/metrics"),
		semconv.HTTPStatusCode(http.StatusOK),
	)

	// assert metrics
	expectedMetric := `
		# HELP core_http_server_requests_total Number of processed HTTP requests
		# TYPE core_http_server_requests_total counter
		core_http_server_requests_total{method="GET",path="/",status="2xx"} 2
		core_http_server_requests_total{method="GET",path="/metrics",status="2xx"} 1
	`

	err := testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"core_http_server_requests_total",
	)
	assert.NoError(t, err)
}

func TestModuleWithMetricsEnabledAndCollectedWithNamespace(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")
	t.Setenv("METRICS_NAMESPACE", "foo")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter, &metricsRegistry))

	// [GET] / twice to generate some metrics
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	core.HttpServer().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()

	core.HttpServer().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// [GET] /metrics to check the metrics
	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// assert [GET] / counter = 2
	assert.Contains(t, rec.Body.String(), `foo_core_http_server_requests_duration_seconds_bucket{method="GET",path="/",le="1"} 2`)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"uri":     "/metrics",
		"status":  200,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /metrics",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/metrics"),
		semconv.HTTPStatusCode(http.StatusOK),
	)

	// assert metrics
	expectedMetric := `
		# HELP foo_core_http_server_requests_total Number of processed HTTP requests
		# TYPE foo_core_http_server_requests_total counter
		foo_core_http_server_requests_total{method="GET",path="/",status="2xx"} 2
		foo_core_http_server_requests_total{method="GET",path="/metrics",status="2xx"} 1
	`

	err := testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"foo_core_http_server_requests_total",
	)
	assert.NoError(t, err)
}

func TestModuleWithHealthcheckDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("STARTUP_ENABLED", "false")
	t.Setenv("LIVENESS_ENABLED", "false")
	t.Setenv("READINESS_ENABLED", "false")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /healthz
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/healthz",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"GET /healthz",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/healthz"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)

	// [GET] /livez
	logBuffer.Reset()
	traceExporter.Reset()

	req = httptest.NewRequest(http.MethodGet, "/livez", nil)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/livez",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"GET /livez",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/livez"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)

	// [GET] /readyz
	logBuffer.Reset()
	traceExporter.Reset()

	req = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/readyz",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"GET /readyz",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/readyz"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleWithHealthcheckEnabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("STARTUP_ENABLED", "true")
	t.Setenv("LIVENESS_ENABLED", "true")
	t.Setenv("READINESS_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(
		t,
		fxhealthcheck.AsCheckerProbe(probes.NewSuccessProbe),
		fxhealthcheck.AsCheckerProbe(probes.NewFailureProbe, healthcheck.Liveness),
		fx.Populate(&core, &logBuffer, &traceExporter),
	)

	// [GET] /healthz
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t,
		`{"success":true,"probes":{"successProbe":{"success":true,"message":"success"}}}`,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
	)

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"method":  "GET",
		"uri":     "/healthz",
		"message": "request logger",
	})

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"GET /healthz",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/healthz"),
		semconv.HTTPStatusCode(http.StatusOK),
	)

	// [GET] /livez
	logBuffer.Reset()

	req = httptest.NewRequest(http.MethodGet, "/livez", nil)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t,
		`{"success":false,"probes":{"failureProbe":{"success":false,"message":"failure"},"successProbe":{"success":true,"message":"success"}}}`,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "error",
		"service":      "core-app",
		"module":       "core",
		"successProbe": "success: true, message: success",
		"failureProbe": "success: false, message: failure",
		"message":      "healthcheck failure",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"service": "core-app",
		"module":  "core",
		"method":  "GET",
		"uri":     "/livez",
		"message": "request logger",
	})

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"GET /livez",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/livez"),
		semconv.HTTPStatusCode(http.StatusInternalServerError),
	)

	// [GET] /readyz
	logBuffer.Reset()

	req = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t,
		`{"success":true,"probes":{"successProbe":{"success":true,"message":"success"}}}`,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
	)

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"method":  "GET",
		"uri":     "/readyz",
		"message": "request logger",
	})

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"GET /readyz",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/readyz"),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestModuleWithDebugConfigDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("CONFIG_ENABLED", "false")
	t.Setenv("APP_DEBUG", "false")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/config
	req := httptest.NewRequest(http.MethodGet, "/debug/config", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/config",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/config",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/config"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleWithDebugConfigEnabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("CONFIG_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/config
	req := httptest.NewRequest(http.MethodGet, "/debug/config", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
		`"name":"core-app"`,
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/config",
		"status":  200,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/config",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/config"),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestModuleWithDebugPprofDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PPROF_ENABLED", "false")
	t.Setenv("APP_DEBUG", "false")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/pprof/
	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/pprof/",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/pprof/",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/pprof/"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleWithDebugPprofEnabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PPROF_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/pprof/
	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
		`<title>/debug/pprof/</title>`,
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/pprof/",
		"status":  200,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/pprof/",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/pprof/"),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestModuleWithDebugRoutesDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("ROUTES_ENABLED", "false")
	t.Setenv("APP_DEBUG", "false")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/routes
	req := httptest.NewRequest(http.MethodGet, "/debug/routes", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/routes",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/routes",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/routes"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleWithDebugRoutesEnabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("ROUTES_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/routes
	req := httptest.NewRequest(http.MethodGet, "/debug/routes", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
		`"path":"/debug/routes"`,
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/routes",
		"status":  200,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/routes",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/routes"),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestModuleWithDebugStatsDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("STATS_ENABLED", "false")
	t.Setenv("APP_DEBUG", "false")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/stats/
	req := httptest.NewRequest(http.MethodGet, "/debug/stats/", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/stats/",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/stats/",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/stats/"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleWithDebugStatsEnabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("STATS_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/stats/
	req := httptest.NewRequest(http.MethodGet, "/debug/stats/", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
		`<title>Statsviz</title>`,
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/stats/",
		"status":  200,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/stats/",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/stats/"),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestModuleWithDebugBuildDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("BUILD_ENABLED", "false")
	t.Setenv("APP_DEBUG", "false")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/build
	req := httptest.NewRequest(http.MethodGet, "/debug/build", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/build",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/build",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/build"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleWithDebugBuildEnabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("BUILD_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/build
	req := httptest.NewRequest(http.MethodGet, "/debug/build", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
		`"version"`,
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/build",
		"status":  200,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/build",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/build"),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestModuleWithDebugModulesDisabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("MODULES_ENABLED", "false")
	t.Setenv("APP_DEBUG", "false")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/modules/core
	req := httptest.NewRequest(http.MethodGet, "/debug/modules/core", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/modules/core",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/modules/core",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/modules/core"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleWithDebugModulesEnabled(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("MODULES_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] /debug/modules/core
	req := httptest.NewRequest(http.MethodGet, "/debug/modules/core", nil)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		strings.ReplaceAll(strings.ReplaceAll(rec.Body.String(), " ", ""), "\n", ""),
		`"name":"core-app"`,
	)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/modules/core",
		"status":  200,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/modules/:name",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/modules/core"),
		semconv.HTTPStatusCode(http.StatusOK),
	)

	// [GET] /debug/modules/invalid
	req = httptest.NewRequest(http.MethodGet, "/debug/modules/invalid", nil)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), `fx module info with name invalid was not found`)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"service": "core-app",
		"module":  "core",
		"uri":     "/debug/modules/invalid",
		"status":  404,
		"message": "request logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /debug/modules/:name",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/debug/modules/invalid"),
		semconv.HTTPStatusCode(http.StatusNotFound),
	)
}

func TestModuleDashboard(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("MODULES_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [GET] / with light theme cookie
	cookie := &http.Cookie{Name: "theme", Value: "light"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `data-bs-theme="light"`)

	// [GET] / with dark theme cookie
	cookie = &http.Cookie{Name: "theme", Value: "dark"}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `data-bs-theme="dark"`)

	// [GET] / with invalid theme cookie
	cookie = &http.Cookie{Name: "theme", Value: "invalid"}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `data-bs-theme="light"`)

	// [GET] / with no theme cookie
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `data-bs-theme="light"`)
}

//nolint:bodyclose
func TestModuleDashboardTheme(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("MODULES_ENABLED", "true")
	t.Setenv("METRICS_ENABLED", "true")
	t.Setenv("METRICS_COLLECT", "true")

	var core *fxcore.Core
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxcore.NewBootstrapper().RunTestApp(t, fx.Populate(&core, &logBuffer, &traceExporter))

	// [POST] /theme to switch to dark
	data := `{"theme": "dark"}`
	req := httptest.NewRequest(http.MethodPost, "/theme", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMovedPermanently, rec.Code)

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "theme" {
			assert.Equal(t, "dark", cookie.Value)
		}
	}

	// [POST] /theme to switch to light
	data = `{"theme": "light"}`
	req = httptest.NewRequest(http.MethodPost, "/theme", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMovedPermanently, rec.Code)

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "theme" {
			assert.Equal(t, "light", cookie.Value)
		}
	}

	// [POST] /theme to switch to invalid
	data = `{"theme": "invalid"}`
	req = httptest.NewRequest(http.MethodPost, "/theme", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMovedPermanently, rec.Code)

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "theme" {
			assert.Equal(t, "light", cookie.Value)
		}
	}

	// [POST] /theme to switch with invalid
	data = `{"theme": "invalid}`
	req = httptest.NewRequest(http.MethodPost, "/theme", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMovedPermanently, rec.Code)

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "theme" {
			assert.Equal(t, "light", cookie.Value)
		}
	}
}

func TestModuleDashboardTasks(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TASKS_ENABLED", "true")

	var core *fxcore.Core

	fxcore.NewBootstrapper().RunTestApp(
		t,
		fxcore.AsTasks(
			tasks.NewErrorTask,
			tasks.NewSuccessTask,
		),
		fx.Populate(&core),
	)

	// [GET] /tasks/success
	req := httptest.NewRequest(http.MethodPost, "/tasks/success", bytes.NewBuffer([]byte("test input")))
	rec := httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	res := fxcore.TaskResult{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.True(t, res.Success)
	assert.Equal(t, "task success", res.Message)
	assert.Equal(
		t,
		map[string]any{
			"app":   "core-app",
			"input": "test input",
		},
		res.Details,
	)

	// [GET] /tasks/error
	req = httptest.NewRequest(http.MethodPost, "/tasks/error", bytes.NewBuffer([]byte("test input")))
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	res = fxcore.TaskResult{}
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.False(t, res.Success)
	assert.Equal(t, "task error", res.Message)
	assert.Nil(t, res.Details)

	// [GET] /tasks/invalid
	req = httptest.NewRequest(http.MethodPost, "/tasks/invalid", bytes.NewBuffer([]byte("test input")))
	rec = httptest.NewRecorder()
	core.HttpServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	res = fxcore.TaskResult{}
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.False(t, res.Success)
	assert.Equal(t, "task invalid not found", res.Message)
	assert.Nil(t, res.Details)
}
