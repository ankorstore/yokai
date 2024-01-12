package fxhttpserver

import (
	"fmt"

	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

// FxHttpServerModuleInfo is a module info collector for fxcore.
type FxHttpServerModuleInfo struct {
	Port         int
	Debug        bool
	Logger       string
	Binder       string
	Serializer   string
	Renderer     string
	ErrorHandler string
	Routes       []*echo.Route
}

// NewFxHttpServerModuleInfo returns a new [FxHttpServerModuleInfo].
func NewFxHttpServerModuleInfo(httpServer *echo.Echo, cfg *config.Config) *FxHttpServerModuleInfo {
	port := cfg.GetInt("modules.http.server.port")
	if port == 0 {
		port = DefaultPort
	}

	return &FxHttpServerModuleInfo{
		Port:         port,
		Debug:        httpServer.Debug,
		Logger:       fmt.Sprintf("%T", httpServer.Logger),
		Binder:       fmt.Sprintf("%T", httpServer.Binder),
		Serializer:   fmt.Sprintf("%T", httpServer.JSONSerializer),
		Renderer:     fmt.Sprintf("%T", httpServer.Renderer),
		ErrorHandler: fmt.Sprintf("%T", httpServer.HTTPErrorHandler),
		Routes:       httpServer.Routes(),
	}
}

// Name return the name of the module info.
func (i *FxHttpServerModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *FxHttpServerModuleInfo) Data() map[string]interface{} {
	return map[string]interface{}{
		"port":         i.Port,
		"debug":        i.Debug,
		"binder":       i.Binder,
		"serializer":   i.Serializer,
		"renderer":     i.Renderer,
		"errorHandler": i.ErrorHandler,
		"routes":       i.Routes,
	}
}
