package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofAllocsHandler is an [echo.HandlerFunc] for pprof allocs.
func PprofAllocsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Handler("allocs").ServeHTTP(c.Response().Writer, c.Request())

		return nil
	}
}
