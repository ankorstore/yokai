package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Options are options for the [HttpServerFactory] implementations.
type Options struct {
	Debug            bool
	Banner           bool
	Recovery         bool
	Logger           echo.Logger
	Binder           echo.Binder
	JsonSerializer   echo.JSONSerializer
	HttpErrorHandler echo.HTTPErrorHandler
	Renderer         echo.Renderer
}

// DefaultHttpServerOptions are the default options used in the [DefaultHttpServerFactory].
func DefaultHttpServerOptions() Options {
	return Options{
		Debug:            false,
		Banner:           false,
		Recovery:         true,
		Logger:           log.New("default"),
		Binder:           &echo.DefaultBinder{},
		JsonSerializer:   &echo.DefaultJSONSerializer{},
		HttpErrorHandler: nil,
		Renderer:         nil,
	}
}

// HttpServerOption are functional options for the [HttpServerFactory] implementations.
type HttpServerOption func(o *Options)

// WithDebug is used to activate the server debug mode.
func WithDebug(d bool) HttpServerOption {
	return func(o *Options) {
		o.Debug = d
	}
}

// WithBanner is used to activate the server banner.
func WithBanner(b bool) HttpServerOption {
	return func(o *Options) {
		o.Banner = b
	}
}

// WithRecovery is used to activate the server automatic panic recovery.
func WithRecovery(r bool) HttpServerOption {
	return func(o *Options) {
		o.Recovery = r
	}
}

// WithLogger is used to specify a [echo.Logger] to be used by the server.
func WithLogger(l echo.Logger) HttpServerOption {
	return func(o *Options) {
		o.Logger = l
	}
}

// WithBinder is used to specify a [echo.Binder] to be used by the server.
func WithBinder(b echo.Binder) HttpServerOption {
	return func(o *Options) {
		o.Binder = b
	}
}

// WithJsonSerializer is used to specify a [echo.JSONSerializer] to be used by the server.
func WithJsonSerializer(s echo.JSONSerializer) HttpServerOption {
	return func(o *Options) {
		o.JsonSerializer = s
	}
}

// WithHttpErrorHandler is used to specify a [echo.HTTPErrorHandler] to be used by the server.
func WithHttpErrorHandler(h echo.HTTPErrorHandler) HttpServerOption {
	return func(o *Options) {
		o.HttpErrorHandler = h
	}
}

// WithRenderer is used to specify a [echo.Renderer] to be used by the server.
func WithRenderer(r echo.Renderer) HttpServerOption {
	return func(o *Options) {
		o.Renderer = r
	}
}
