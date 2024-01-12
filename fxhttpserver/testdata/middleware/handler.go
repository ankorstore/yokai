package middleware

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
)

type TestHandlerMiddleware struct {
	config *config.Config
}

func NewTestHandlerMiddleware(config *config.Config) *TestHandlerMiddleware {
	return &TestHandlerMiddleware{
		config: config,
	}
}

func (m *TestHandlerMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.CtxLogger(c.Request().Context()).Info().Msgf("HANDLER middleware for app: %s", m.config.AppName())

			c.Response().Header().Add("handler-middleware", "true")

			return next(c)
		}
	}
}
