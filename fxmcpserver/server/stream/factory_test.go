package stream_test

import (
	"github.com/ankorstore/yokai/fxmcpserver/server/stream"
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMCPStreamableHTTPServerFactory_Create(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../../testdata/config"),
	)
	assert.NoError(t, err)

	mcpSrv := &server.MCPServer{}

	fac := stream.NewDefaultMCPStreamableHTTPServerFactory(cfg)

	srv := fac.Create(mcpSrv)

	assert.IsType(t, (*server.StreamableHTTPServer)(nil), srv.Server())

	assert.Equal(t, ":0", srv.Config().Address)
	assert.True(t, srv.Config().Stateless)
	assert.Equal(t, stream.DefaultBasePath, srv.Config().BasePath)
	assert.True(t, srv.Config().KeepAlive)
	assert.Equal(t, stream.DefaultKeepAliveInterval, srv.Config().KeepAliveInterval)

	assert.False(t, srv.Running())
}
