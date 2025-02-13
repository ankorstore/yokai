package handler

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type TestErrorHandler struct {
}

func NewTestErrorHandler() *TestErrorHandler {
	return &TestErrorHandler{}
}

func (h *TestErrorHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		return errors.New("test error")
	}
}
