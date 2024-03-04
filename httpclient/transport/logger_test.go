package transport_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ankorstore/yokai/httpclient/transport"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type transportMock struct {
	mock.Mock
}

func (m *transportMock) RoundTrip(*http.Request) (*http.Response, error) {
	args := m.Called()

	return nil, args.Error(1)
}

func TestNewLoggerTransport(t *testing.T) {
	t.Parallel()

	trans := transport.NewLoggerTransport(nil)

	assert.IsType(t, &transport.LoggerTransport{}, trans)
	assert.Implements(t, (*http.RoundTripper)(nil), trans)
}

func TestLoggerTransportBase(t *testing.T) {
	t.Parallel()

	base := &http.Transport{}

	trans := transport.NewLoggerTransport(base)

	assert.Equal(t, base, trans.Base())
}

func TestLoggerTransportRoundTrip(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, server.URL, nil)
	req = req.WithContext(logger.WithContext(context.Background()))

	resp, err := transport.NewLoggerTransport(nil).RoundTrip(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"method":  "GET",
		"url":     server.URL,
		"message": "http client request",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"url":     server.URL,
		"code":    http.StatusNoContent,
		"message": "http client response",
	})
}

func TestLoggerTransportRoundTripWithConfig(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedStatus, err := strconv.Atoi(r.Header.Get("expected-response-code"))
		assert.NoError(t, err)

		w.WriteHeader(expectedStatus)

		_, err = w.Write([]byte(r.Header.Get("expected-response-body")))
		assert.NoError(t, err)
	}))
	defer server.Close()

	// transport
	trans := transport.NewLoggerTransportWithConfig(nil, &transport.LoggerTransportConfig{
		LogRequest:                       true,
		LogResponse:                      true,
		LogRequestBody:                   true,
		LogResponseBody:                  true,
		LogRequestLevel:                  zerolog.DebugLevel,
		LogResponseLevel:                 zerolog.DebugLevel,
		LogResponseLevelFromResponseCode: true,
	})

	// 200 response
	data := []byte(`{"input":"data"}`)
	req := httptest.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(data))
	req.Header.Add("expected-response-code", "200")
	req.Header.Add("expected-response-body", `{"output":"ok"}`)
	req = req.WithContext(logger.WithContext(context.Background()))

	resp, err := trans.RoundTrip(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "debug",
		"method":  "POST",
		"url":     server.URL,
		"request": `{"input":"data"}`,
		"message": "http client request",
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":    "debug",
		"url":      server.URL,
		"code":     http.StatusOK,
		"response": `{"output":"ok"}`,
		"message":  "http client response",
	})

	// 400 response
	data = []byte(`{"input":"data"}`)
	req = httptest.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(data))
	req.Header.Add("expected-response-code", "400")
	req.Header.Add("expected-response-body", `{"output":"bad request"}`)
	req = req.WithContext(logger.WithContext(context.Background()))

	resp, err = trans.RoundTrip(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "debug",
		"method":  "POST",
		"url":     server.URL,
		"request": `{"input":"data"}`,
		"message": "http client request",
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":    "warn",
		"url":      server.URL,
		"code":     http.StatusBadRequest,
		"response": `{"output":"bad request"}`,
		"message":  "http client response",
	})

	// 500 response
	data = []byte(`{"input":"data"}`)
	req = httptest.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(data))
	req.Header.Add("expected-response-code", "500")
	req.Header.Add("expected-response-body", `{"output":"error"}`)
	req = req.WithContext(logger.WithContext(context.Background()))

	resp, err = trans.RoundTrip(req)
	assert.NoError(t, err)

	err = resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "debug",
		"method":  "POST",
		"url":     server.URL,
		"request": `{"input":"data"}`,
		"message": "http client request",
	})

	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":    "error",
		"url":      server.URL,
		"code":     http.StatusInternalServerError,
		"response": `{"output":"error"}`,
		"message":  "http client response",
	})
}

func TestLoggerTransportRoundTripWithFailure(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, server.URL, nil)
	req = req.WithContext(logger.WithContext(context.Background()))

	base := new(transportMock)
	base.On("RoundTrip", mock.Anything).Return(nil, fmt.Errorf("custom http error"))

	//nolint:bodyclose
	resp, err := transport.NewLoggerTransport(base).RoundTrip(req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "custom http error", err.Error())

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"method":  "GET",
		"url":     server.URL,
		"message": "http client request",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"error":   "custom http error",
		"level":   "error",
		"message": "http client failure",
	})
}
