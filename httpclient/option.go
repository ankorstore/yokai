package httpclient

import (
	"net/http"
	"time"

	"github.com/ankorstore/yokai/httpclient/transport"
)

// Options are options for the [HttpClientFactory] implementations.
type Options struct {
	Transport     http.RoundTripper
	CheckRedirect func(req *http.Request, via []*http.Request) error
	Jar           http.CookieJar
	Timeout       time.Duration
}

// DefaultHttpClientOptions are the default options used in the [DefaultHttpClientFactory].
func DefaultHttpClientOptions() Options {
	return Options{
		Transport:     transport.NewBaseTransport(),
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 30,
	}
}

// HttpClientOption are functional options for the [HttpClientFactory] implementations.
type HttpClientOption func(o *Options)

// WithTransport is used to specify the [http.RoundTripper] to use by the [http.Client].
func WithTransport(t http.RoundTripper) HttpClientOption {
	return func(o *Options) {
		o.Transport = t
	}
}

// WithCheckRedirect is used to specify the check redirect func to use by the [http.Client].
func WithCheckRedirect(f func(req *http.Request, via []*http.Request) error) HttpClientOption {
	return func(o *Options) {
		o.CheckRedirect = f
	}
}

// WithCookieJar is used to specify the [http.CookieJar] to use by the [http.Client].
func WithCookieJar(j http.CookieJar) HttpClientOption {
	return func(o *Options) {
		o.Jar = j
	}
}

// WithTimeout is used to specify the timeout to use by the [http.Client].
func WithTimeout(t time.Duration) HttpClientOption {
	return func(o *Options) {
		o.Timeout = t
	}
}
