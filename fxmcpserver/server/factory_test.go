package server_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	fs "github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMCPServerFactory_Create(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../testdata/config"),
	)
	assert.NoError(t, err)

	fac := fs.NewDefaultMCPServerFactory(cfg)

	srv := fac.Create()

	assert.IsType(t, (*server.MCPServer)(nil), srv)
}
