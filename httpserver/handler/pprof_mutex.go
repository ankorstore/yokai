package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofMutexHandler is an [echo.HandlerFunc] for pprof mutex.
func PprofMutexHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Handler("mutex").ServeHTTP(c.Response().Writer, c.Request())

		return nil
	}
}
