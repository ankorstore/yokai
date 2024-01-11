package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofThreadCreateHandler is an [echo.HandlerFunc] for pprof threadcreate.
func PprofThreadCreateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(c.Response().Writer, c.Request())

		return nil
	}
}
