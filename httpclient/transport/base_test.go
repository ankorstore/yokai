package transport_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/httpclient/transport"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseTransport(t *testing.T) {
	t.Parallel()

	trans := transport.NewBaseTransport()

	assert.IsType(t, &transport.BaseTransport{}, trans)
	assert.Implements(t, (*http.RoundTripper)(nil), trans)
}

func TestBaseTransportBaseWithConfig(t *testing.T) {
	t.Parallel()

	trans := transport.NewBaseTransportWithConfig(
		&transport.BaseTransportConfig{
			MaxIdleConnections:        50,
			MaxConnectionsPerHost:     50,
			MaxIdleConnectionsPerHost: 50,
		},
	)

	assert.Equal(t, 50, trans.Base().MaxIdleConns)
	assert.Equal(t, 50, trans.Base().MaxConnsPerHost)
	assert.Equal(t, 50, trans.Base().MaxIdleConnsPerHost)
}

func TestBaseTransportRoundTrip(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, server.URL, nil)

	resp, err := transport.NewBaseTransport().RoundTrip(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
