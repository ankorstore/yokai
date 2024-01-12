package fxcore_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/stretchr/testify/assert"
)

func TestNewCore(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("./testdata/config"),
	)
	assert.NoError(t, err)

	checker, err := healthcheck.NewDefaultCheckerFactory().Create()
	assert.NoError(t, err)

	server, err := httpserver.NewDefaultHttpServerFactory().Create()
	assert.NoError(t, err)

	core := fxcore.NewCore(cfg, checker, server)
	assert.IsType(t, &fxcore.Core{}, core)

	assert.Equal(t, cfg, core.Config())
	assert.Equal(t, checker, core.Checker())
	assert.Equal(t, server, core.HttpServer())
}
