package httpserver

import (
	"github.com/labstack/echo/v4"
)

// HttpServerFactory is the interface for [echo.Echo] factories.
type HttpServerFactory interface {
	Create(options ...HttpServerOption) (*echo.Echo, error)
}

// DefaultHttpServerFactory is the default [HttpServerFactory] implementation.
type DefaultHttpServerFactory struct{}

// NewDefaultHttpServerFactory returns a [DefaultHttpServerFactory], implementing [HttpServerFactory].
func NewDefaultHttpServerFactory() HttpServerFactory {
	return &DefaultHttpServerFactory{}
}

// Create returns a new [echo.Echo], and accepts a list of [HttpServerOption].
// For example:
//
//	var server, _ = httpserver.NewDefaultHttpServerFactory().Create()
//
// is equivalent to:
//
//	var server, _ = httpserver.NewDefaultHttpServerFactory().Create(
//		httpserver.WithDebug(false),                                  // debug disabled by default
//		httpserver.WithBanner(false),                                 // banner disabled by default
//		httpserver.WithLogger(log.New("default")),                    // echo default logger
//		httpserver.WithBinder(&echo.DefaultBinder{}),                 // echo default binder
//		httpserver.WithJsonSerializer(&echo.DefaultJSONSerializer{}), // echo default json serializer
//		httpserver.WithHttpErrorHandler(nil),                         // echo default error handler
//	)
func (f *DefaultHttpServerFactory) Create(options ...HttpServerOption) (*echo.Echo, error) {
	appliedOpts := DefaultHttpServerOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	httpServer := echo.New()

	httpServer.Debug = appliedOpts.Debug
	httpServer.HideBanner = !appliedOpts.Banner

	httpServer.Logger = appliedOpts.Logger
	httpServer.Binder = appliedOpts.Binder
	httpServer.JSONSerializer = appliedOpts.JsonSerializer

	if appliedOpts.HttpErrorHandler != nil {
		httpServer.HTTPErrorHandler = appliedOpts.HttpErrorHandler
	}

	if appliedOpts.Renderer != nil {
		httpServer.Renderer = appliedOpts.Renderer
	}

	return httpServer, nil
}
