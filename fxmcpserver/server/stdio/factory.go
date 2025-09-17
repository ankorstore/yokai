package stdio

import (
	"os"

	"github.com/mark3labs/mcp-go/server"
)

var _ MCPStdioServerFactory = (*DefaultMCPStdioServerFactory)(nil)

// MCPStdioServerFactory is the interface for MCP Stdio server factories.
type MCPStdioServerFactory interface {
	Create(mcpServer *server.MCPServer, options ...server.StdioOption) *MCPStdioServer
}

// DefaultMCPStdioServerFactory is the default MCPStdioServerFactory implementation.
type DefaultMCPStdioServerFactory struct{}

// NewDefaultMCPStdioServerFactory returns a new DefaultMCPStdioServerFactory instance.
func NewDefaultMCPStdioServerFactory() *DefaultMCPStdioServerFactory {
	return &DefaultMCPStdioServerFactory{}
}

// Create returns a new MCPStdioServer instance.
func (f *DefaultMCPStdioServerFactory) Create(mcpServer *server.MCPServer, options ...server.StdioOption) *MCPStdioServer {
	srvConfig := MCPStdioServerConfig{
		In:  os.Stdin,
		Out: os.Stdout,
	}

	return NewMCPStdioServer(mcpServer, srvConfig, options...)
}
