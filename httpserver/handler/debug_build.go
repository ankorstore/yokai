package handler

import (
	"net/http"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

// DebugBuildHandler is an [echo.HandlerFunc] that returns build information.
func DebugBuildHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "error reading build info")
		}

		modInfo := map[string]string{}
		for _, dep := range buildInfo.Deps {
			modInfo[dep.Path] = dep.Version
		}

		return c.JSON(http.StatusOK, echo.Map{
			"env": echo.Map{
				"arch": runtime.GOARCH,
				"os":   runtime.GOOS,
				"vars": os.Environ(),
			},
			"go": echo.Map{
				"main":    buildInfo.Main.Path,
				"modules": modInfo,
				"version": buildInfo.GoVersion,
			},
		})
	}
}
