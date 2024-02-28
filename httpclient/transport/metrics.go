package transport

import (
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai/httpclient/status"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	HttpClientMetricsRequestsCount    = "client_requests_total"
	HttpClientMetricsRequestsDuration = "client_request_duration_seconds"
)

// MetricsTransport is a wrapper around [http.RoundTripper] with some [MetricsTransportConfig] configuration.
type MetricsTransport struct {
	transport        http.RoundTripper
	config           *MetricsTransportConfig
	requestsCounter  *prometheus.CounterVec
	requestsDuration *prometheus.HistogramVec
}

// MetricsTransportConfig is the configuration of the [MetricsTransport].
type MetricsTransportConfig struct {
	Registry            prometheus.Registerer
	Namespace           string
	Subsystem           string
	Buckets             []float64
	NormalizeHTTPStatus bool
}

// NewMetricsTransport returns a [MetricsTransport] instance with default [MetricsTransportConfig] configuration.
func NewMetricsTransport(base http.RoundTripper) *MetricsTransport {
	return NewMetricsTransportWithConfig(
		base,
		&MetricsTransportConfig{
			Registry:            prometheus.DefaultRegisterer,
			Namespace:           "",
			Subsystem:           "",
			Buckets:             prometheus.DefBuckets,
			NormalizeHTTPStatus: true,
		},
	)
}

// NewMetricsTransportWithConfig returns a [MetricsTransport] instance for a provided [MetricsTransportConfig] configuration.
func NewMetricsTransportWithConfig(base http.RoundTripper, config *MetricsTransportConfig) *MetricsTransport {
	if base == nil {
		base = NewBaseTransport()
	}

	if config.Registry == nil {
		config.Registry = prometheus.DefaultRegisterer
	}

	requestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      HttpClientMetricsRequestsCount,
			Help:      "Number of performed HTTP requests",
		},
		[]string{
			"url",
			"method",
			"status",
		},
	)

	requestsDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Name:      HttpClientMetricsRequestsDuration,
			Help:      "Time spent performing HTTP requests",
			Buckets:   config.Buckets,
		},
		[]string{
			"url",
			"method",
		},
	)

	config.Registry.MustRegister(requestsCounter, requestsDuration)

	return &MetricsTransport{
		transport:        base,
		config:           config,
		requestsCounter:  requestsCounter,
		requestsDuration: requestsDuration,
	}
}

// Base returns the wrapped [http.RoundTripper].
func (t *MetricsTransport) Base() http.RoundTripper {
	return t.transport
}

// RoundTrip performs a request / response round trip, based on the wrapped [http.RoundTripper].
func (t *MetricsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	timer := prometheus.NewTimer(t.requestsDuration.WithLabelValues(req.URL.String(), req.Method))
	resp, err := t.transport.RoundTrip(req)
	timer.ObserveDuration()

	respStatus := ""
	if t.config.NormalizeHTTPStatus {
		respStatus = status.NormalizeHTTPStatus(resp.StatusCode)
	} else {
		respStatus = strconv.Itoa(resp.StatusCode)
	}

	t.requestsCounter.WithLabelValues(req.URL.String(), req.Method, respStatus).Inc()

	return resp, err
}
