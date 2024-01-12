package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/ankorstore/yokai/fxhttpserver/testdata/service"
)

type TestTemplateHandler struct {
	service *service.TestService
}

func NewTestTemplateHandler(service *service.TestService) *TestTemplateHandler {
	return &TestTemplateHandler{
		service: service,
	}
}

func (h *TestTemplateHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Info("in template handler")

		return c.Render(http.StatusOK, "test.html", map[string]interface{}{
			"name": h.service.GetAppName(),
		})
	}
}
