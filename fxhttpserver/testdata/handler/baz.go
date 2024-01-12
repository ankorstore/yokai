package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/ankorstore/yokai/fxhttpserver/testdata/service"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
)

type TestBazHandler struct {
	service *service.TestService
}

func NewTestBazHandler(service *service.TestService) *TestBazHandler {
	return &TestBazHandler{
		service: service,
	}
}

func (h *TestBazHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := httpserver.CtxTracer(c).Start(c.Request().Context(), "baz span")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("in baz handler")

		return c.String(http.StatusOK, fmt.Sprintf("baz: %s", h.service.GetAppName()))
	}
}
