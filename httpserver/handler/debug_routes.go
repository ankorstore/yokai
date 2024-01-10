package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DebugRoutesHandler is an [echo.HandlerFunc] that returns routing information.
func DebugRoutesHandler(httpServer *echo.Echo) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, httpServer.Routes())
	}
}
