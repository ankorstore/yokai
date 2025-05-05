package server_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	fs "github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMCPServerHooksProvider_Provide(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../testdata/config"),
	)
	assert.NoError(t, err)

	pro := fs.NewDefaultMCPServerHooksProvider(reg, cfg)

	hooks := pro.Provide()

	assert.IsType(t, (*server.Hooks)(nil), hooks)
}
