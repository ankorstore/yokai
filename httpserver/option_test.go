package httpserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestWithDebug(t *testing.T) {
	t.Parallel()

	opt := httpserver.DefaultHttpServerOptions()
	httpserver.WithDebug(true)(&opt)

	assert.Equal(t, true, opt.Debug)
}

func TestWithBanner(t *testing.T) {
	t.Parallel()

	opt := httpserver.DefaultHttpServerOptions()
	httpserver.WithBanner(true)(&opt)

	assert.Equal(t, true, opt.Banner)
}

func TestWithRecovery(t *testing.T) {
	t.Parallel()

	opt := httpserver.DefaultHttpServerOptions()
	httpserver.WithRecovery(false)(&opt)

	assert.Equal(t, false, opt.Recovery)
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	opt := httpserver.DefaultHttpServerOptions()
	logger := log.New("test")
	httpserver.WithLogger(logger)(&opt)

	assert.Equal(t, logger, opt.Logger)
}

func TestWithBinder(t *testing.T) {
	t.Parallel()

	opt := httpserver.DefaultHttpServerOptions()
	binder := &echo.DefaultBinder{}
	httpserver.WithBinder(binder)(&opt)

	assert.Equal(t, binder, opt.Binder)
}

func TestWithJsonSerializer(t *testing.T) {
	t.Parallel()

	opt := httpserver.DefaultHttpServerOptions()
	jsonSerializer := &echo.DefaultJSONSerializer{}
	httpserver.WithJsonSerializer(jsonSerializer)(&opt)

	assert.Equal(t, jsonSerializer, opt.JsonSerializer)
}

func TestWithHttpErrorHandler(t *testing.T) {
	t.Parallel()

	opt := httpserver.DefaultHttpServerOptions()
	httpErrorHandler := func(err error, c echo.Context) {}
	httpserver.WithHttpErrorHandler(httpErrorHandler)(&opt)

	assert.NotNil(t, opt.HttpErrorHandler)
}

func TestWithRenderer(t *testing.T) {
	t.Parallel()

	opt := httpserver.DefaultHttpServerOptions()
	httpserver.WithRenderer(httpserver.NewHtmlTemplateRenderer("testdata/templates/*.html"))(&opt)

	assert.NotNil(t, opt.Renderer)
}
