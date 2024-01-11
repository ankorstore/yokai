package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofTraceHandler is an [echo.HandlerFunc] for pprof trace.
func PprofTraceHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Trace(c.Response().Writer, c.Request())

		return nil
	}
}
