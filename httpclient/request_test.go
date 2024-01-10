package httpclient_test

import (
	"net/http"
	"testing"

	"github.com/ankorstore/yokai/httpclient"
	"github.com/stretchr/testify/assert"
)

func TestCopyRequestHeaders(t *testing.T) {
	t.Parallel()

	source, err := http.NewRequest(http.MethodGet, "https://test.com", nil)
	assert.NoError(t, err)

	source.Header.Add("foo", "foo1")
	source.Header.Add("Bar", "bar")
	source.Header.Add("ignore", "ignore")

	dest, err := http.NewRequest(http.MethodGet, "https://other-test.com", nil)
	assert.NoError(t, err)

	source.Header.Add("foo", "foo2")
	dest.Header.Add("Baz", "baz")

	httpclient.CopyRequestHeaders(source, dest, "foo", "bar")

	assert.Equal(t, "foo1", dest.Header.Get("foo"))
	assert.Equal(t, []string{"foo1", "foo2"}, dest.Header.Values("foo"))

	assert.Equal(t, "bar", dest.Header.Get("bar"))
	assert.Equal(t, []string{"bar"}, dest.Header.Values("bar"))

	assert.Equal(t, "baz", dest.Header.Get("baz"))
	assert.Equal(t, []string{"baz"}, dest.Header.Values("baz"))
}

func TestCopyObservabilityRequestHeaders(t *testing.T) {
	t.Parallel()

	source, err := http.NewRequest(http.MethodGet, "https://test.com", nil)
	assert.NoError(t, err)

	source.Header.Add("x-request-id", "test-request-id")
	source.Header.Add("traceparent", "test-traceparent")

	dest, err := http.NewRequest(http.MethodGet, "https://other-test.com", nil)
	assert.NoError(t, err)

	httpclient.CopyObservabilityRequestHeaders(source, dest)

	assert.Equal(t, "test-request-id", dest.Header.Get("x-request-id"))
	assert.Equal(t, "test-traceparent", dest.Header.Get("traceparent"))
}
