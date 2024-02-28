package transport_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/httpclient/transport"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetricsTransportRoundTrip(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, server.URL, nil)

	trans := transport.NewMetricsTransport(nil)
	assert.IsType(t, &transport.MetricsTransport{}, trans)
	assert.Implements(t, (*http.RoundTripper)(nil), trans)

	resp, err := trans.RoundTrip(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// requests counter assertions
	expectedCounterMetric := fmt.Sprintf(`
			# HELP client_requests_total Number of performed HTTP requests
			# TYPE client_requests_total counter
			client_requests_total{method="GET",status="5xx",url="%s"} 1
		`,
		server.URL,
	)

	err = testutil.GatherAndCompare(
		prometheus.DefaultGatherer,
		strings.NewReader(expectedCounterMetric),
		"client_requests_total",
	)
	assert.NoError(t, err)
}

func TestMetricsTransportRoundTripWithBaseAndConfig(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, server.URL, nil)

	base := &http.Transport{}

	trans := transport.NewMetricsTransportWithConfig(
		base,
		&transport.MetricsTransportConfig{
			Registry:            registry,
			Namespace:           "foo",
			Subsystem:           "bar",
			Buckets:             []float64{1, 2, 3},
			NormalizeHTTPStatus: false,
		},
	)

	assert.Equal(t, base, trans.Base())

	resp, err := trans.RoundTrip(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// requests counter assertions
	expectedCounterMetric := fmt.Sprintf(`
			# HELP foo_bar_client_requests_total Number of performed HTTP requests
			# TYPE foo_bar_client_requests_total counter
			foo_bar_client_requests_total{method="GET",status="204",url="%s"} 1
		`,
		server.URL,
	)

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedCounterMetric),
		"foo_bar_client_requests_total",
	)
	assert.NoError(t, err)
}
