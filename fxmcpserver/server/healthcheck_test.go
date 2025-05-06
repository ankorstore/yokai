package server_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/config"
	fs "github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/fxmcpserver/server/stdio"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestMCPServerProbe(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../testdata/config"),
	)
	assert.NoError(t, err)

	mcpSrv := server.NewMCPServer("test-server", "1.0.0")

	sseSrv := sse.NewDefaultMCPSSEServerFactory(cfg).Create(mcpSrv)
	stdioSrv := stdio.NewDefaultMCPStdioServerFactory().Create(mcpSrv)

	probe := fs.NewMCPServerProbe(cfg, sseSrv, stdioSrv)

	res := probe.Check(context.Background())

	assert.False(t, res.Success)
	assert.Equal(t, "MCP SSE server is not running", res.Message)
}
