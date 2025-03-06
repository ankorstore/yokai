package httpclienttest

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestHTTPRequestFunc func(tb testing.TB, req *http.Request) error
type TestHTTPResponseFunc func(tb testing.TB, w http.ResponseWriter) error

type TestHTTPRoundTrip struct {
	RequestFunc  TestHTTPRequestFunc
	ResponseFunc TestHTTPResponseFunc
}

func NewTestHTTPServer(tb testing.TB, options ...TestHTTPServerOptionFunc) *httptest.Server {
	tb.Helper()

	var mu sync.Mutex

	serverOptions := DefaultTestHTTPServerOptions()
	for _, opt := range options {
		opt(&serverOptions)
	}

	if len(serverOptions.RoundtripsStack) == 0 {
		tb.Error("test HTTP server: empty roundtrips stack")

		return nil
	}

	stackPosition := 0

	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()

				if stackPosition >= len(serverOptions.RoundtripsStack) {
					tb.Error("test HTTP server: roundtrips stack exhausted")

					return
				}

				err := serverOptions.RoundtripsStack[stackPosition].RequestFunc(tb, r)
				assert.NoError(tb, err)

				err = serverOptions.RoundtripsStack[stackPosition].ResponseFunc(tb, w)
				assert.NoError(tb, err)

				stackPosition++
			},
		),
	)
}
