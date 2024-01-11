package middleware

import (
	"context"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestIdMiddlewareConfig is the configuration for the [RequestIdMiddleware].
type RequestIdMiddlewareConfig struct {
	Skipper         middleware.Skipper
	Generator       uuid.UuidGenerator
	RequestIdHeader string
}

// DefaultRequestIdMiddlewareConfig is the default configuration for the [RequestIdMiddleware].
var DefaultRequestIdMiddlewareConfig = RequestIdMiddlewareConfig{
	Skipper:         middleware.DefaultSkipper,
	Generator:       uuid.NewDefaultUuidGenerator(),
	RequestIdHeader: echo.HeaderXRequestID,
}

// RequestIdMiddleware returns a [RequestIdMiddleware] with the [DefaultRequestIdMiddlewareConfig].
func RequestIdMiddleware() echo.MiddlewareFunc {
	return RequestIdMiddlewareWithConfig(DefaultRequestIdMiddlewareConfig)
}

// RequestIdMiddlewareWithConfig returns a [RequestIdMiddleware] for a provided [RequestIdMiddlewareConfig].
func RequestIdMiddlewareWithConfig(config RequestIdMiddlewareConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIdMiddlewareConfig.Skipper
	}

	if config.Generator == nil {
		config.Generator = DefaultRequestIdMiddlewareConfig.Generator
	}

	if config.RequestIdHeader == "" {
		config.RequestIdHeader = echo.HeaderXRequestID
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			resp := c.Response()

			// request_id req / resp header propagation
			rid := req.Header.Get(config.RequestIdHeader)

			if rid == "" {
				rid = config.Generator.Generate()
				req.Header.Set(config.RequestIdHeader, rid)
			}

			resp.Header().Set(config.RequestIdHeader, rid)

			// request_id ctx propagation
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), httpserver.CtxRequestIdKey{}, rid)))

			return next(c)
		}
	}
}
