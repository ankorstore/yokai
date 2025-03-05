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

func NewTestHTTPServer(tb testing.TB, roundtripStack ...TestHTTPRoundTrip) *httptest.Server {
	tb.Helper()

	var mu sync.Mutex

	if len(roundtripStack) == 0 {
		tb.Fatal("test HTTP server: empty roundtrips stack")

		return nil
	}

	stackPosition := 0

	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()

				if stackPosition >= len(roundtripStack) {
					tb.Error("test HTTP server: roundtrips stack exhausted")
				}

				err := roundtripStack[stackPosition].RequestFunc(tb, r)
				assert.NoError(tb, err)

				err = roundtripStack[stackPosition].ResponseFunc(tb, w)
				assert.NoError(tb, err)

				stackPosition++
			},
		),
	)
}
