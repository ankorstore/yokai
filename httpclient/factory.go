package httpclient

import (
	"net/http"
)

// HttpClientFactory is the interface for [http.Client] factories.
type HttpClientFactory interface {
	Create(opts ...HttpClientOption) (*http.Client, error)
}

// DefaultHttpClientFactory is the default [HttpClientFactory] implementation.
type DefaultHttpClientFactory struct{}

// NewDefaultHttpClientFactory returns a [DefaultHttpClientFactory], implementing [HttpClientFactory].
func NewDefaultHttpClientFactory() HttpClientFactory {
	return &DefaultHttpClientFactory{}
}

// Create returns a new [http.Client], and accepts a list of [HttpClientOption].
// For example:
//
//	var client, _ = httpclient.NewDefaultHttpClientFactory().Create()
//
//	// equivalent to:
//	var client, _ = httpclient.NewDefaultHttpClientFactory().Create(
//		httpclient.WithTransport(transport.NewBaseTransport()), // base http transport (optimized)
//		httpclient.WithTimeout(30*time.Second),                 // 30 seconds timeout
//		httpclient.WithCheckRedirect(nil),                      // default redirection checks
//		httpclient.WithCookieJar(nil),                          // default cookie jar
//	)
func (f *DefaultHttpClientFactory) Create(options ...HttpClientOption) (*http.Client, error) {
	appliedOpts := DefaultHttpClientOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	return &http.Client{
		Transport:     appliedOpts.Transport,
		CheckRedirect: appliedOpts.CheckRedirect,
		Jar:           appliedOpts.Jar,
		Timeout:       appliedOpts.Timeout,
	}, nil
}
