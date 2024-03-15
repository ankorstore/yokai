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
	"github.com/stretchr/testify/mock"
)

func TestMetricsTransportRoundTrip(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	trans := transport.NewMetricsTransport(nil)
	assert.IsType(t, &transport.MetricsTransport{}, trans)
	assert.Implements(t, (*http.RoundTripper)(nil), trans)

	req := httptest.NewRequest(http.MethodGet, server.URL, nil)

	resp, err := trans.RoundTrip(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// requests counter assertions
	expectedCounterMetric := fmt.Sprintf(
		`
			# HELP http_client_requests_total Number of performed HTTP requests
			# TYPE http_client_requests_total counter
			http_client_requests_total{host="%s",method="GET",path="",status="5xx"} 1
		`,
		server.URL,
	)

	err = testutil.GatherAndCompare(
		prometheus.DefaultGatherer,
		strings.NewReader(expectedCounterMetric),
		"http_client_requests_total",
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

	base := &http.Transport{}

	trans := transport.NewMetricsTransportWithConfig(
		base,
		&transport.MetricsTransportConfig{
			Registry:             registry,
			Namespace:            "foo",
			Subsystem:            "bar",
			Buckets:              []float64{1, 2, 3},
			NormalizeRequestPath: true,
			NormalizeRequestPathMasks: map[string]string{
				`/foo/(.+)/bar\?page=(.+)`: "/foo/{fooId}/bar?page={pageId}",
			},
			NormalizeResponseStatus: false,
		},
	)

	assert.Equal(t, base, trans.Base())

	// requests
	urls := []string{
		server.URL,
		fmt.Sprintf("%s/foo/1/bar?page=1#baz", server.URL),
		fmt.Sprintf("%s/foo/2/bar?page=2#baz", server.URL),
		fmt.Sprintf("%s/foo/3/bar?page=3#baz", server.URL),
		fmt.Sprintf("%s/foo/4/baz", server.URL),
	}

	for _, url := range urls {
		req := httptest.NewRequest(http.MethodGet, url, nil)

		resp, err := trans.RoundTrip(req)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	}

	// requests counter assertions
	expectedCounterMetric := fmt.Sprintf(
		`
			# HELP foo_bar_http_client_requests_total Number of performed HTTP requests
			# TYPE foo_bar_http_client_requests_total counter
    		foo_bar_http_client_requests_total{host="%s",method="GET",path="",status="204"} 1
    		foo_bar_http_client_requests_total{host="%s",method="GET",path="/foo/4/baz",status="204"} 1
    		foo_bar_http_client_requests_total{host="%s",method="GET",path="/foo/{fooId}/bar?page={pageId}",status="204"} 3
		`,
		server.URL,
		server.URL,
		server.URL,
	)

	err := testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedCounterMetric),
		"foo_bar_http_client_requests_total",
	)
	assert.NoError(t, err)
}

func TestMetricsTransportRoundTripWithFailure(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	base := new(transportMock)
	base.On("RoundTrip", mock.Anything).Return(nil, fmt.Errorf("custom http error"))

	trans := transport.NewMetricsTransportWithConfig(
		base,
		&transport.MetricsTransportConfig{
			Registry:  registry,
			Namespace: "foo",
			Subsystem: "bar",
		},
	)

	assert.Equal(t, base, trans.Base())

	// request
	req := httptest.NewRequest(http.MethodGet, server.URL, nil)

	//nolint:bodyclose
	resp, err := trans.RoundTrip(req)
	assert.Nil(t, resp)
	assert.Error(t, err)

	// requests counter assertions
	expectedCounterMetric := fmt.Sprintf(
		`
			# HELP foo_bar_http_client_requests_total Number of performed HTTP requests
			# TYPE foo_bar_http_client_requests_total counter
			foo_bar_http_client_requests_total{host="%s",method="GET",path="",status="error"} 1
		`,
		server.URL,
	)

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedCounterMetric),
		"foo_bar_http_client_requests_total",
	)
	assert.NoError(t, err)
}
