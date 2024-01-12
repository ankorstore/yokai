package factory

import (
	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
)

type TestHttpServerFactory struct{}

func NewTestHttpServerFactory() httpserver.HttpServerFactory {
	return &TestHttpServerFactory{}
}

func (f *TestHttpServerFactory) Create(options ...httpserver.HttpServerOption) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = false

	return e, nil
}
