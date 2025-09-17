package stdio_test

import (
	"os"
	"testing"

	"github.com/ankorstore/yokai/fxmcpserver/server/stdio"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMCPStdioServerFactory_Create(t *testing.T) {
	t.Parallel()

	mcpSrv := &server.MCPServer{}

	fac := stdio.NewDefaultMCPStdioServerFactory()

	srv := fac.Create(mcpSrv)

	assert.IsType(t, (*server.StdioServer)(nil), srv.Server())

	assert.Equal(t, os.Stdin, srv.Config().In)
	assert.Equal(t, os.Stdout, srv.Config().Out)

	assert.False(t, srv.Running())
}
