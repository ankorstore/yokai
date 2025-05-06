package sse_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMCPSSEServerFactory_Create(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../../testdata/config"),
	)
	assert.NoError(t, err)

	mcpSrv := &server.MCPServer{}

	fac := sse.NewDefaultMCPSSEServerFactory(cfg)

	srv := fac.Create(mcpSrv)

	assert.IsType(t, (*server.SSEServer)(nil), srv.Server())

	assert.Equal(t, ":0", srv.Config().Address)
	assert.Equal(t, sse.DefaultBaseURL, srv.Config().BaseURL)
	assert.Equal(t, sse.DefaultBasePath, srv.Config().BasePath)
	assert.Equal(t, sse.DefaultSSEEndpoint, srv.Config().SSEEndpoint)
	assert.Equal(t, sse.DefaultMessageEndpoint, srv.Config().MessageEndpoint)
	assert.True(t, srv.Config().KeepAlive)
	assert.Equal(t, sse.DefaultKeepAliveInterval, srv.Config().KeepAliveInterval)

	assert.False(t, srv.Running())
}
