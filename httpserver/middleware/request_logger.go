package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

const (
	HeaderXRequestId  = "x-request-id"
	HeaderTraceParent = "traceparent"
	LogFieldRequestId = "requestID"
)

// RequestLoggerMiddlewareConfig is the configuration for the [RequestLoggerMiddleware].
type RequestLoggerMiddlewareConfig struct {
	Skipper                         middleware.Skipper
	LogLevelFromResponseOrErrorCode bool
	RequestHeadersToLog             map[string]string
	RequestUriPrefixesToExclude     []string
}

// DefaultRequestLoggerMiddlewareConfig is the default configuration for the [RequestLoggerMiddleware].
var DefaultRequestLoggerMiddlewareConfig = RequestLoggerMiddlewareConfig{
	Skipper:                         middleware.DefaultSkipper,
	LogLevelFromResponseOrErrorCode: false,
	RequestHeadersToLog:             map[string]string{HeaderXRequestId: LogFieldRequestId},
	RequestUriPrefixesToExclude:     []string{},
}

// RequestLoggerMiddleware returns a [RequestLoggerMiddleware] with the [DefaultRequestLoggerMiddlewareConfig].
func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return RequestLoggerMiddlewareWithConfig(DefaultRequestLoggerMiddlewareConfig)
}

// RequestLoggerMiddlewareWithConfig returns a [RequestLoggerMiddleware] for a provided [RequestLoggerMiddlewareConfig].
//
//nolint:gocognit,nestif
func RequestLoggerMiddlewareWithConfig(config RequestLoggerMiddlewareConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIdMiddlewareConfig.Skipper
	}

	if config.RequestHeadersToLog == nil {
		config.RequestHeadersToLog = DefaultRequestLoggerMiddlewareConfig.RequestHeadersToLog
	}

	if config.RequestUriPrefixesToExclude == nil {
		config.RequestUriPrefixesToExclude = DefaultRequestLoggerMiddlewareConfig.RequestUriPrefixesToExclude
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// skipper
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			// logger preparation
			headersToLog := map[string]interface{}{}
			for headerNameToLog, logFieldName := range config.RequestHeadersToLog {
				headerValueToLog := req.Header.Get(headerNameToLog)

				if headerValueToLog == "" {
					headerValueToLog = res.Header().Get(headerNameToLog)
				}

				if headerValueToLog != "" {
					headersToLog[logFieldName] = headerValueToLog
				}
			}

			// request id context propagation
			requestId := req.Header.Get(HeaderXRequestId)
			if requestId == "" {
				requestId = res.Header().Get(HeaderXRequestId)
			}
			ctx := context.WithValue(req.Context(), httpserver.CtxRequestIdKey{}, requestId)

			// logger context propagation
			var logger zerolog.Logger
			if echoLogger, ok := c.Logger().(*httpserver.EchoLogger); ok {
				logger = echoLogger.ToZerolog().With().Fields(headersToLog).Logger()
			} else {
				logger = log.CtxLogger(req.Context()).With().Fields(headersToLog).Logger()
			}

			c.SetRequest(c.Request().WithContext(logger.WithContext(ctx)))
			c.SetLogger(httpserver.NewEchoLogger(log.FromZerolog(logger)))

			// invoke next in chain
			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			// trigger error handler
			if err != nil {
				c.Error(err)
			}

			// response status
			status := res.Status
			if err != nil {
				var httpErr *echo.HTTPError
				if errors.As(err, &httpErr) {
					status = httpErr.Code
				} else {
					status = http.StatusInternalServerError
				}
			}

			// skip if matching exclusions and not error or code > 500
			if httpserver.MatchPrefix(config.RequestUriPrefixesToExclude, req.RequestURI) &&
				err == nil &&
				status < http.StatusInternalServerError {
				return nil
			}

			// log event preparation
			var evt *zerolog.Event
			if config.LogLevelFromResponseOrErrorCode {
				if err != nil {
					var he *echo.HTTPError
					if errors.As(err, &he) {
						switch {
						case he.Code >= http.StatusBadRequest && he.Code < http.StatusInternalServerError:
							evt = logger.Warn()
						case he.Code >= http.StatusInternalServerError:
							evt = logger.Error()
						default:
							evt = logger.Info()
						}
					} else {
						evt = logger.Error().Err(err)
					}

					evt.Str(zerolog.ErrorFieldName, err.Error())
				} else {
					switch {
					case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
						evt = logger.Warn()
					case status >= http.StatusInternalServerError:
						evt = logger.Error()
					default:
						evt = logger.Info()
					}
				}
			} else {
				evt = logger.Info()

				if err != nil {
					evt.Str(zerolog.ErrorFieldName, err.Error())
				}
			}

			// log event tracing
			spanContext := trace.SpanContextFromContext(c.Request().Context())

			if spanContext.HasTraceID() {
				evt.Str("traceID", spanContext.TraceID().String())
			}

			if spanContext.HasSpanID() {
				evt.Str("spanID", spanContext.SpanID().String())
			}

			// log event propagation
			evt.
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", status).
				Str("latency", latency.String()).
				Str("remoteIp", c.RealIP()).
				Str("referer", req.Referer()).
				Str("userAgent", req.UserAgent()).
				Msg("request logger")

			// error propagation
			return err
		}
	}
}
