package fxcore

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/labstack/echo/v4"
)

// Core is the core component, holding the core config, health checker and http server.
type Core struct {
	config     *config.Config
	checker    *healthcheck.Checker
	httpServer *echo.Echo
}

// NewCore returns a new [Core].
func NewCore(config *config.Config, checker *healthcheck.Checker, httpServer *echo.Echo) *Core {
	return &Core{
		config:     config,
		checker:    checker,
		httpServer: httpServer,
	}
}

// Config returns the core config.
func (c *Core) Config() *config.Config {
	return c.config
}

// Checker returns the core health checker.
func (c *Core) Checker() *healthcheck.Checker {
	return c.checker
}

// HttpServer returns the core http server.
func (c *Core) HttpServer() *echo.Echo {
	return c.httpServer
}
