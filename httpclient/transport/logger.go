package transport

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
)

// LoggerTransport is a wrapper around [http.RoundTripper] with some [LoggerTransportConfig] configuration.
type LoggerTransport struct {
	transport http.RoundTripper
	config    *LoggerTransportConfig
}

// LoggerTransportConfig is the configuration of the [LoggerTransport].
type LoggerTransportConfig struct {
	LogRequest                       bool
	LogResponse                      bool
	LogRequestBody                   bool
	LogResponseBody                  bool
	LogRequestLevel                  zerolog.Level
	LogResponseLevel                 zerolog.Level
	LogResponseLevelFromResponseCode bool
}

// NewLoggerTransport returns a [LoggerTransport] instance with default [LoggerTransportConfig] configuration.
func NewLoggerTransport(base http.RoundTripper) *LoggerTransport {
	return NewLoggerTransportWithConfig(
		base,
		&LoggerTransportConfig{
			LogRequest:                       false,
			LogResponse:                      false,
			LogRequestBody:                   false,
			LogResponseBody:                  false,
			LogRequestLevel:                  zerolog.InfoLevel,
			LogResponseLevel:                 zerolog.InfoLevel,
			LogResponseLevelFromResponseCode: false,
		},
	)
}

// NewLoggerTransportWithConfig returns a [LoggerTransport] instance for a provided [LoggerTransportConfig] configuration.
func NewLoggerTransportWithConfig(base http.RoundTripper, config *LoggerTransportConfig) *LoggerTransport {
	if base == nil {
		base = NewBaseTransport()
	}

	return &LoggerTransport{
		transport: base,
		config:    config,
	}
}

// Base returns the wrapped [http.RoundTripper].
func (t *LoggerTransport) Base() http.RoundTripper {
	return t.transport
}

// RoundTrip performs a request / response round trip, based on the wrapped [http.RoundTripper].
//
//nolint:cyclop
func (t *LoggerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	logger := log.CtxLogger(req.Context())

	if t.config.LogRequest {
		reqEvt := logger.WithLevel(t.config.LogRequestLevel)

		reqDump, err := httputil.DumpRequestOut(req, t.config.LogRequestBody)
		if err == nil {
			reqEvt.Bytes("request", reqDump)
		}

		reqEvt.
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Msg("http client request")
	}

	start := time.Now()
	resp, err := t.transport.RoundTrip(req)
	latency := time.Since(start).String()

	if err != nil {
		logger.Error().Err(err).Str("latency", latency).Msg("http client failure")

		return resp, err
	}

	if t.config.LogResponse {
		var respEvt *zerolog.Event

		if t.config.LogResponseLevelFromResponseCode {
			switch {
			case resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError:
				respEvt = logger.Warn()
			case resp.StatusCode >= http.StatusInternalServerError:
				respEvt = logger.Error()
			default:
				respEvt = logger.WithLevel(t.config.LogResponseLevel)
			}
		} else {
			respEvt = logger.WithLevel(t.config.LogResponseLevel)
		}

		respDump, err := httputil.DumpResponse(resp, t.config.LogResponseBody)
		if err == nil {
			respEvt.Bytes("response", respDump)
		}

		respEvt.
			Str("method", resp.Request.Method).
			Str("url", resp.Request.URL.String()).
			Int("code", resp.StatusCode).
			Str("latency", latency).
			Msg("http client response")
	}

	return resp, err
}
