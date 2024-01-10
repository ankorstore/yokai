package transport

import (
	"net/http"
)

// BaseTransport is a wrapper around [http.Transport] with some [BaseTransportConfig] configuration.
type BaseTransport struct {
	transport *http.Transport
	config    *BaseTransportConfig
}

// BaseTransportConfig is the configuration of the [BaseTransport].
type BaseTransportConfig struct {
	MaxIdleConnections        int
	MaxConnectionsPerHost     int
	MaxIdleConnectionsPerHost int
}

// NewBaseTransport returns a [BaseTransport] instance with optimized default [BaseTransportConfig] configuration.
func NewBaseTransport() *BaseTransport {
	return NewBaseTransportWithConfig(
		&BaseTransportConfig{
			MaxIdleConnections:        100,
			MaxConnectionsPerHost:     100,
			MaxIdleConnectionsPerHost: 100,
		},
	)
}

// NewBaseTransportWithConfig returns a [BaseTransport] instance for a provided [BaseTransportConfig] configuration.
func NewBaseTransportWithConfig(config *BaseTransportConfig) *BaseTransport {
	//nolint:forcetypeassert
	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.MaxIdleConns = config.MaxIdleConnections
	transport.MaxConnsPerHost = config.MaxConnectionsPerHost
	transport.MaxIdleConnsPerHost = config.MaxIdleConnectionsPerHost

	return &BaseTransport{
		transport: transport,
		config:    config,
	}
}

// Base returns the wrapped [http.Transport].
func (t *BaseTransport) Base() *http.Transport {
	return t.transport
}

// RoundTrip performs a request / response round trip, based on the wrapped [http.Transport].
func (t *BaseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.transport.RoundTrip(req)
}
