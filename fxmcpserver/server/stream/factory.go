package stream

import (
	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
)

const (
	DefaultAddr = ":8083"
)

var _ MCPStreamableHTTPServerFactory = (*DefaultMCPStreamableHTTPServerFactory)(nil)

// MCPStreamableHTTPServerFactory is the interface for MCP StreamableHTTP server factories.
type MCPStreamableHTTPServerFactory interface {
	Create(mcpServer *server.MCPServer, options ...server.StreamableHTTPOption) *MCPStreamableHTTPServer
}

// DefaultMCPStreamableHTTPServerFactory is the default MCPStreamableHTTPServerFactory implementation.
type DefaultMCPStreamableHTTPServerFactory struct {
	config *config.Config
}

// NewDefaultMCPStreamableHTTPServerFactory returns a new DefaultMCPStreamableHTTPServerFactory instance.
func NewDefaultMCPStreamableHTTPServerFactory(config *config.Config) *DefaultMCPStreamableHTTPServerFactory {
	return &DefaultMCPStreamableHTTPServerFactory{
		config: config,
	}
}

// Create returns a new MCPStreamableHTTPServer instance.
func (f *DefaultMCPStreamableHTTPServerFactory) Create(mcpServer *server.MCPServer, options ...server.StreamableHTTPOption) *MCPStreamableHTTPServer {
	addr := f.config.GetString("modules.mcp.server.transport.stream.address")
	if addr == "" {
		addr = DefaultAddr
	}

	srvConfig := MCPStreamableHTTPServerConfig{
		Address: addr,
	}

	srvOptions := []server.StreamableHTTPOption{
		server.WithStateLess(true),
	}

	srvOptions = append(srvOptions, options...)

	return NewMCPStreamableHTTPServer(mcpServer, srvConfig, srvOptions...)
}
