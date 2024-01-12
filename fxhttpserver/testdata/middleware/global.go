package middleware

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
)

type TestGlobalMiddleware struct {
	config *config.Config
}

func NewTestGlobalMiddleware(config *config.Config) *TestGlobalMiddleware {
	return &TestGlobalMiddleware{
		config: config,
	}
}

func (m *TestGlobalMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.CtxLogger(c.Request().Context()).Info().Msgf("GLOBAL middleware for app: %s", m.config.AppName())

			c.Response().Header().Add("global-middleware", "true")

			return next(c)
		}
	}
}
