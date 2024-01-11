package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofGoroutineHandler is an [echo.HandlerFunc] for pprof goroutine.
func PprofGoroutineHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(c.Response().Writer, c.Request())

		return nil
	}
}
