package middleware

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
)

type TestGroupMiddleware struct {
	config *config.Config
}

func NewTestGroupMiddleware(config *config.Config) *TestGroupMiddleware {
	return &TestGroupMiddleware{
		config: config,
	}
}

func (m *TestGroupMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.CtxLogger(c.Request().Context()).Info().Msgf("GROUP middleware for app: %s", m.config.AppName())

			c.Response().Header().Add("group-middleware", "true")

			return next(c)
		}
	}
}
