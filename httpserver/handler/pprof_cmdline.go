package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofCmdlineHandler is an [echo.HandlerFunc] for pprof cmdline.
func PprofCmdlineHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Cmdline(c.Response().Writer, c.Request())

		return nil
	}
}
