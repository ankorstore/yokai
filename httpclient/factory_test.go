package httpclient_test

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/ankorstore/yokai/httpclient"
	"github.com/ankorstore/yokai/httpclient/transport"
	"github.com/stretchr/testify/assert"
)

func TestDefaultHttpClientFactory(t *testing.T) {
	t.Parallel()

	factory := httpclient.NewDefaultHttpClientFactory()

	assert.IsType(t, &httpclient.DefaultHttpClientFactory{}, factory)
	assert.Implements(t, (*httpclient.HttpClientFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	factory := httpclient.NewDefaultHttpClientFactory()

	checkRedirectFunc := func(req *http.Request, via []*http.Request) error {
		return fmt.Errorf("custom error")
	}
	jar, _ := cookiejar.New(nil)
	timeout := time.Second * 20

	options := []httpclient.HttpClientOption{
		httpclient.WithCheckRedirect(checkRedirectFunc),
		httpclient.WithCookieJar(jar),
		httpclient.WithTimeout(timeout),
	}

	client, err := factory.Create(options...)

	assert.NoError(t, err)
	assert.IsType(t, &transport.BaseTransport{}, client.Transport)
	assert.Equal(t, jar, client.Jar)
	assert.Equal(t, timeout, client.Timeout)

	req, err := http.NewRequest(http.MethodGet, "https://test.com", nil)
	assert.NoError(t, err)

	err = client.CheckRedirect(req, []*http.Request{})

	assert.Error(t, err)
	assert.Equal(t, "custom error", err.Error())
}
