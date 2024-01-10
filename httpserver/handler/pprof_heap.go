package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofHeapHandler is an [echo.HandlerFunc] for pprof heap.
func PprofHeapHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Handler("heap").ServeHTTP(c.Response(), c.Request())

		return nil
	}
}
