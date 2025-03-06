package httpclienttest

type TestHTTPServerOptions struct {
	RoundtripsStack []TestHTTPRoundTrip
}

func DefaultTestHTTPServerOptions() TestHTTPServerOptions {
	return TestHTTPServerOptions{
		RoundtripsStack: []TestHTTPRoundTrip{},
	}
}

type TestHTTPServerOptionFunc func(o *TestHTTPServerOptions)

func WithTestHTTPRoundTrip(reqFunc TestHTTPRequestFunc, responseFunc TestHTTPResponseFunc) TestHTTPServerOptionFunc {
	return func(o *TestHTTPServerOptions) {
		o.RoundtripsStack = append(o.RoundtripsStack, TestHTTPRoundTrip{
			RequestFunc:  reqFunc,
			ResponseFunc: responseFunc,
		})
	}
}
