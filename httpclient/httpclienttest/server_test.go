package httpclienttest_test

import (
	"net/http"
	"testing"

	"github.com/ankorstore/yokai/httpclient/httpclienttest"
	"github.com/stretchr/testify/assert"
)

func TestTestHTTPServer(t *testing.T) {
	t.Parallel()

	t.Run("test failure on request fn failed assertion", func(t *testing.T) {
		t.Parallel()

		testReqFn := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.True(tb, false)

			return nil
		}

		testRespFn := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			return nil
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(mt, testReqFn, testRespFn)
		defer testServer.Close()

		resp, err := testServer.Client().Get(testServer.URL)
		assert.NoError(t, err)

		assert.True(t, mt.Failed())

		err = resp.Body.Close()
		assert.NoError(t, err)
	})
}
