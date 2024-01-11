package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofProfileHandler is an [echo.HandlerFunc] for pprof profile.
func PprofProfileHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Profile(c.Response().Writer, c.Request())

		return nil
	}
}
