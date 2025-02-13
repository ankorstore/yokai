package errorhandler

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

type TestErrorHandler struct {
	config *config.Config
}

func NewTestErrorHandler(config *config.Config) *TestErrorHandler {
	return &TestErrorHandler{
		config: config,
	}
}

func (h *TestErrorHandler) Handle() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		c.String(http.StatusInternalServerError, fmt.Sprintf("error handled in test error handler of %s: %s", h.config.AppName(), err))
	}
}
