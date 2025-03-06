package httpclienttest_test

import (
	"errors"
	"net/http"
	"sync"
	"testing"

	"github.com/ankorstore/yokai/httpclient/httpclienttest"
	"github.com/stretchr/testify/assert"
)

//nolint:goconst,maintidx
func TestTestHTTPServer(t *testing.T) {
	t.Parallel()

	t.Run("test success with single roundtrip", func(t *testing.T) {
		t.Parallel()

		testReqFn := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/foo", r.URL.Path)

			return nil
		}

		testRespFn := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			w.Header().Set("foo", "bar")

			w.WriteHeader(http.StatusOK)

			return nil
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(mt, httpclienttest.WithTestHTTPRoundTrip(testReqFn, testRespFn))
		defer testServer.Close()

		resp, err := testServer.Client().Get(testServer.URL + "/foo")
		assert.NoError(t, err)

		assert.False(t, mt.Failed())

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "bar", resp.Header.Get("foo"))

		err = resp.Body.Close()
		assert.NoError(t, err)
	})

	t.Run("test success with multiple sequential roundtrips stack", func(t *testing.T) {
		t.Parallel()

		testReqFn1 := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/foo", r.URL.Path)

			return nil
		}

		testRespFn1 := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			w.Header().Set("foo", "foo")

			return nil
		}

		testReqFn2 := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/bar", r.URL.Path)

			return nil
		}

		testRespFn2 := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			w.Header().Set("bar", "bar")

			return nil
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(
			mt,
			httpclienttest.WithTestHTTPRoundTrip(testReqFn1, testRespFn1),
			httpclienttest.WithTestHTTPRoundTrip(testReqFn2, testRespFn2),
		)
		defer testServer.Close()

		resp, err := testServer.Client().Get(testServer.URL + "/foo")
		assert.NoError(t, err)

		assert.False(t, mt.Failed())

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "foo", resp.Header.Get("foo"))

		err = resp.Body.Close()
		assert.NoError(t, err)

		resp, err = testServer.Client().Get(testServer.URL + "/bar")
		assert.NoError(t, err)

		assert.False(t, mt.Failed())

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "bar", resp.Header.Get("bar"))

		err = resp.Body.Close()
		assert.NoError(t, err)
	})

	t.Run("test success with multiple concurrent roundtrips stack", func(t *testing.T) {
		t.Parallel()

		testReqFn := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/foo", r.URL.Path)

			return nil
		}

		testRespFn := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			w.Header().Set("foo", "foo")

			return nil
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(
			mt,
			httpclienttest.WithTestHTTPRoundTrip(testReqFn, testRespFn),
			httpclienttest.WithTestHTTPRoundTrip(testReqFn, testRespFn),
		)

		var wg sync.WaitGroup
		wg.Add(2)

		go func(tb testing.TB, twg *sync.WaitGroup) {
			tb.Helper()

			defer twg.Done()

			resp, err := testServer.Client().Get(testServer.URL + "/foo")
			assert.NoError(t, err)

			assert.False(t, mt.Failed())

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "foo", resp.Header.Get("foo"))

			err = resp.Body.Close()
			assert.NoError(t, err)
		}(mt, &wg)

		go func(tb testing.TB, twg *sync.WaitGroup) {
			tb.Helper()

			defer twg.Done()

			resp, err := testServer.Client().Get(testServer.URL + "/foo")
			assert.NoError(t, err)

			assert.False(t, mt.Failed())

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "foo", resp.Header.Get("foo"))

			err = resp.Body.Close()
			assert.NoError(t, err)
		}(mt, &wg)

		wg.Wait()

		testServer.Close()
	})

	t.Run("test error with empty roundtrips stack", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)

		httpclienttest.NewTestHTTPServer(mt)

		assert.True(t, mt.Failed())
	})

	t.Run("test error with exhausted roundtrip stack", func(t *testing.T) {
		t.Parallel()

		testReqFn := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/foo", r.URL.Path)

			return nil
		}

		testRespFn := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			w.Header().Set("foo", "bar")

			return nil
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(mt, httpclienttest.WithTestHTTPRoundTrip(testReqFn, testRespFn))
		defer testServer.Close()

		resp, err := testServer.Client().Get(testServer.URL + "/foo")
		assert.NoError(t, err)

		assert.False(t, mt.Failed())

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "bar", resp.Header.Get("foo"))

		err = resp.Body.Close()
		assert.NoError(t, err)

		resp, err = testServer.Client().Get(testServer.URL + "/foo")
		assert.NoError(t, err)

		assert.True(t, mt.Failed())

		err = resp.Body.Close()
		assert.NoError(t, err)
	})

	t.Run("test error with request func error return", func(t *testing.T) {
		t.Parallel()

		testReqFn := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/foo", r.URL.Path)

			return errors.New("request error foo")
		}

		testRespFn := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			return nil
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(mt, httpclienttest.WithTestHTTPRoundTrip(testReqFn, testRespFn))
		defer testServer.Close()

		resp, err := testServer.Client().Get(testServer.URL + "/foo")
		assert.NoError(t, err)

		assert.True(t, mt.Failed())

		err = resp.Body.Close()
		assert.NoError(t, err)
	})

	t.Run("test error with request func error assertion", func(t *testing.T) {
		t.Parallel()

		testReqFn := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/foo", r.URL.Path)

			return nil
		}

		testRespFn := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			return nil
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(mt, httpclienttest.WithTestHTTPRoundTrip(testReqFn, testRespFn))
		defer testServer.Close()

		resp, err := testServer.Client().Get(testServer.URL + "/invalid")
		assert.NoError(t, err)

		assert.True(t, mt.Failed())

		err = resp.Body.Close()
		assert.NoError(t, err)
	})

	t.Run("test error with response func error return", func(t *testing.T) {
		t.Parallel()

		testReqFn := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/foo", r.URL.Path)

			return nil
		}

		testRespFn := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			return errors.New("response error foo")
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(mt, httpclienttest.WithTestHTTPRoundTrip(testReqFn, testRespFn))
		defer testServer.Close()

		resp, err := testServer.Client().Get(testServer.URL + "/foo")
		assert.NoError(t, err)

		assert.True(t, mt.Failed())

		err = resp.Body.Close()
		assert.NoError(t, err)
	})

	t.Run("test error with response func error assertion", func(t *testing.T) {
		t.Parallel()

		testReqFn := func(tb testing.TB, r *http.Request) error {
			tb.Helper()

			assert.Equal(tb, "/foo", r.URL.Path)

			return nil
		}

		testRespFn := func(tb testing.TB, w http.ResponseWriter) error {
			tb.Helper()

			assert.True(tb, false)

			return nil
		}

		mt := new(testing.T)

		testServer := httpclienttest.NewTestHTTPServer(mt, httpclienttest.WithTestHTTPRoundTrip(testReqFn, testRespFn))
		defer testServer.Close()

		resp, err := testServer.Client().Get(testServer.URL + "/foo")
		assert.NoError(t, err)

		assert.True(t, mt.Failed())

		err = resp.Body.Close()
		assert.NoError(t, err)
	})
}
