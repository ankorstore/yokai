package handler

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PprofSymbolHandler is an [echo.HandlerFunc] for pprof symbol.
func PprofSymbolHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		pprof.Symbol(c.Response().Writer, c.Request())

		return nil
	}
}
