package httpclient_test

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/ankorstore/yokai/httpclient"
	"github.com/stretchr/testify/assert"
)

func TestWithTransport(t *testing.T) {
	t.Parallel()

	opts := httpclient.DefaultHttpClientOptions()

	transport := &http.Transport{}
	httpclient.WithTransport(transport)(&opts)

	assert.Equal(t, transport, opts.Transport)
}

func TestWithCheckRedirect(t *testing.T) {
	t.Parallel()

	opts := httpclient.DefaultHttpClientOptions()

	req, _ := http.NewRequest(http.MethodGet, "https://test.com", nil)
	checkRedirectFunc := func(req *http.Request, via []*http.Request) error {
		return fmt.Errorf("custom error")
	}
	httpclient.WithCheckRedirect(checkRedirectFunc)(&opts)

	err := opts.CheckRedirect(req, []*http.Request{})
	assert.Error(t, err)
	assert.Equal(t, "custom error", err.Error())
}

func TestWithCookieJar(t *testing.T) {
	t.Parallel()

	opts := httpclient.DefaultHttpClientOptions()

	jar, _ := cookiejar.New(nil)
	httpclient.WithCookieJar(jar)(&opts)

	assert.Equal(t, jar, opts.Jar)
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()

	opts := httpclient.DefaultHttpClientOptions()

	timeout := time.Second * 20
	httpclient.WithTimeout(timeout)(&opts)

	assert.Equal(t, timeout, opts.Timeout)
}
