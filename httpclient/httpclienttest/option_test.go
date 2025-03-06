package httpclienttest_test

import (
	"net/http"
	"testing"

	"github.com/ankorstore/yokai/httpclient/httpclienttest"
	"github.com/stretchr/testify/assert"
)

func TestTestHTTPServerOptions(t *testing.T) {
	t.Parallel()

	defaultOptions := httpclienttest.DefaultTestHTTPServerOptions()
	assert.Len(t, defaultOptions.RoundtripsStack, 0)

	option := httpclienttest.WithTestHTTPRoundTrip(
		func(tb testing.TB, req *http.Request) error {
			tb.Helper()

			return nil
		},
		func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			return nil
		},
	)

	option(&defaultOptions)
	assert.Len(t, defaultOptions.RoundtripsStack, 1)
}
