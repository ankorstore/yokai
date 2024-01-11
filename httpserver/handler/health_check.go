package handler

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
)

// HealthCheckHandler is an [echo.HandlerFunc] returns the execution result of a [healthcheck.Checker] for a [healthcheck.ProbeKind].
func HealthCheckHandler(checker *healthcheck.Checker, kind healthcheck.ProbeKind) echo.HandlerFunc {
	return func(c echo.Context) error {
		result := checker.Check(c.Request().Context(), kind)

		status := http.StatusOK
		if !result.Success {
			status = http.StatusInternalServerError

			evt := httpserver.CtxLogger(c).Error()
			for probeName, probeResult := range result.ProbesResults {
				evt.Str(probeName, fmt.Sprintf("success: %v, message: %s", probeResult.Success, probeResult.Message))
			}

			evt.Msg("healthcheck failure")
		}

		return c.JSON(status, result)
	}
}
