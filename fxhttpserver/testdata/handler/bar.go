package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/ankorstore/yokai/fxhttpserver/testdata/service"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
)

type TestBarHandler struct {
	service *service.TestService
}

func NewTestBarHandler(service *service.TestService) *TestBarHandler {
	return &TestBarHandler{
		service: service,
	}
}

func (h *TestBarHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := httpserver.CtxTracer(c).Start(c.Request().Context(), "bar span")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("in bar handler")

		return c.String(http.StatusOK, fmt.Sprintf("bar: %s", h.service.GetAppName()))
	}
}
