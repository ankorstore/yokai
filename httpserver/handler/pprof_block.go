package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofBlockHandler is an [echo.HandlerFunc] for pprof block.
func PprofBlockHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Handler("block").ServeHTTP(c.Response().Writer, c.Request())

		return nil
	}
}
