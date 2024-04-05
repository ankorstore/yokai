package fxhttpserver_test

import (
	"github.com/labstack/echo/v4"
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/stretchr/testify/assert"
)

func TestNewFxHttpServerModuleInfo(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("./testdata/config"),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.Debug = true

	info := fxhttpserver.NewFxHttpServerModuleInfo(httpServer, cfg)
	assert.IsType(t, &fxhttpserver.FxHttpServerModuleInfo{}, info)

	assert.Equal(t, fxhttpserver.ModuleName, info.Name())
	assert.Equal(
		t,
		map[string]interface{}{
			"address":      fxhttpserver.DefaultAddress,
			"debug":        true,
			"binder":       "*echo.DefaultBinder",
			"serializer":   "*echo.DefaultJSONSerializer",
			"renderer":     "<nil>",
			"errorHandler": "echo.HTTPErrorHandler",
			"routes":       []*echo.Route{},
		},
		info.Data(),
	)
}
