package fxhttpclient_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxhttpclient"
	"github.com/ankorstore/yokai/fxhttpclient/testdata/factory"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModule(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var httpClient *http.Client
	var logger *log.Logger
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxmetrics.FxMetricsModule,
		fxtrace.FxTraceModule,
		fxhttpclient.FxHttpClientModule,
		fx.Populate(&httpClient, &logger, &logBuffer, &traceExporter, &metricsRegistry),
	).RequireStart().RequireStop()

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedStatus, err := strconv.Atoi(r.Header.Get("expected-response-code"))
		assert.NoError(t, err)

		w.WriteHeader(expectedStatus)

		_, err = w.Write([]byte(r.Header.Get("expected-response-body")))
		assert.NoError(t, err)
	}))
	defer httpServer.Close()

	// 200 response
	data := []byte(`{"input":"data"}`)
	req := httptest.NewRequest(http.MethodPost, httpServer.URL, bytes.NewBuffer(data))
	req.RequestURI = ""
	req.Header.Add("expected-response-code", "200")
	req.Header.Add("expected-response-body", `{"output":"ok"}`)
	req = req.WithContext(logger.WithContext(context.Background()))

	resp, err := httpClient.Do(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"method":  "POST",
		"url":     httpServer.URL,
		"request": `{"input":"data"}`,
		"message": "http client request",
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":    "info",
		"url":      httpServer.URL,
		"code":     http.StatusOK,
		"response": `{"output":"ok"}`,
		"message":  "http client response",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"HTTP POST",
		semconv.HTTPMethod(http.MethodPost),
		semconv.HTTPURL(httpServer.URL),
		semconv.HTTPStatusCode(http.StatusOK),
	)

	expectedMetric := fmt.Sprintf(
		`
			# HELP app_httpclient_client_requests_total Number of performed HTTP requests
			# TYPE app_httpclient_client_requests_total counter
			app_httpclient_client_requests_total{method="POST",status="2xx",url="%s"} 1
		`,
		httpServer.URL,
	)

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"app_httpclient_client_requests_total",
	)
	assert.NoError(t, err)

	// 400 response
	data = []byte(`{"input":"data"}`)
	req = httptest.NewRequest(http.MethodPost, httpServer.URL, bytes.NewBuffer(data))
	req.RequestURI = ""
	req.Header.Add("expected-response-code", "400")
	req.Header.Add("expected-response-body", `{"output":"bad request"}`)
	req = req.WithContext(logger.WithContext(context.Background()))

	resp, err = httpClient.Do(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"method":  "POST",
		"url":     httpServer.URL,
		"request": `{"input":"data"}`,
		"message": "http client request",
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":    "warn",
		"url":      httpServer.URL,
		"code":     http.StatusBadRequest,
		"response": `{"output":"bad request"}`,
		"message":  "http client response",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"HTTP POST",
		semconv.HTTPMethod(http.MethodPost),
		semconv.HTTPURL(httpServer.URL),
		semconv.HTTPStatusCode(http.StatusBadRequest),
	)

	expectedMetric = fmt.Sprintf(
		`
			# HELP app_httpclient_client_requests_total Number of performed HTTP requests
			# TYPE app_httpclient_client_requests_total counter
			app_httpclient_client_requests_total{method="POST",status="2xx",url="%s"} 1
			app_httpclient_client_requests_total{method="POST",status="4xx",url="%s"} 1
		`,
		httpServer.URL,
		httpServer.URL,
	)

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"app_httpclient_client_requests_total",
	)
	assert.NoError(t, err)

	// 500 response
	data = []byte(`{"input":"data"}`)
	req = httptest.NewRequest(http.MethodPost, httpServer.URL, bytes.NewBuffer(data))
	req.RequestURI = ""
	req.Header.Add("expected-response-code", "500")
	req.Header.Add("expected-response-body", `{"output":"error"}`)
	req = req.WithContext(logger.WithContext(context.Background()))

	resp, err = httpClient.Do(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"method":  "POST",
		"url":     httpServer.URL,
		"request": `{"input":"data"}`,
		"message": "http client request",
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":    "error",
		"url":      httpServer.URL,
		"code":     http.StatusInternalServerError,
		"response": `{"output":"error"}`,
		"message":  "http client response",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"HTTP POST",
		semconv.HTTPMethod(http.MethodPost),
		semconv.HTTPURL(httpServer.URL),
		semconv.HTTPStatusCode(http.StatusInternalServerError),
	)

	expectedMetric = fmt.Sprintf(
		`
			# HELP app_httpclient_client_requests_total Number of performed HTTP requests
			# TYPE app_httpclient_client_requests_total counter
			app_httpclient_client_requests_total{method="POST",status="2xx",url="%s"} 1
			app_httpclient_client_requests_total{method="POST",status="4xx",url="%s"} 1
			app_httpclient_client_requests_total{method="POST",status="5xx",url="%s"} 1
		`,
		httpServer.URL,
		httpServer.URL,
		httpServer.URL,
	)

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"app_httpclient_client_requests_total",
	)
	assert.NoError(t, err)
}

func TestModuleDecoration(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")

	var httpClient *http.Client

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxhttpclient.FxHttpClientModule,
		fx.Decorate(factory.NewTestHttpClientFactory),
		fx.Populate(&httpClient),
	).RequireStart().RequireStop()

	assert.Equal(t, http.DefaultClient, httpClient)
}
