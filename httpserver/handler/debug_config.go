package handler

import (
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

// DebugConfigHandler is an [echo.HandlerFunc] that returns config information.
func DebugConfigHandler(config *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, config.AllSettings())
	}
}
