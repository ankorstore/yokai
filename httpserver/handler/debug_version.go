package handler

import (
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

// DebugVersionHandler is an [echo.HandlerFunc] that returns version information.
func DebugVersionHandler(config *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"application": config.AppName(),
			"version":     config.AppVersion(),
		})
	}
}
