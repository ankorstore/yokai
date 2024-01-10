package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofIndexHandler is an [echo.HandlerFunc] for pprof index dashboard.
func PprofIndexHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Index(c.Response().Writer, c.Request())

		return nil
	}
}
