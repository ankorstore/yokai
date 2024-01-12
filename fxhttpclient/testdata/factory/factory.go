package factory

import (
	"net/http"

	"github.com/ankorstore/yokai/httpclient"
)

type TestHttpClientFactory struct{}

func NewTestHttpClientFactory() httpclient.HttpClientFactory {
	return &TestHttpClientFactory{}
}

func (f *TestHttpClientFactory) Create(options ...httpclient.HttpClientOption) (*http.Client, error) {
	return http.DefaultClient, nil
}
