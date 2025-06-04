package stream

import (
	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
	"time"
)

const (
	DefaultAddr              = ":8083"
	DefaultBasePath          = "/mcp"
	DefaultKeepAliveInterval = 10 * time.Second
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

	stateless := f.config.GetBool("modules.mcp.server.transport.stream.stateless")

	basePath := f.config.GetString("modules.mcp.server.transport.stream.base_path")
	if basePath == "" {
		basePath = DefaultBasePath
	}

	keepAlive := f.config.GetBool("modules.mcp.server.transport.stream.keep_alive")

	keepAliveInterval := DefaultKeepAliveInterval
	keepAliveIntervalConfig := f.config.GetInt("modules.mcp.server.transport.stream.keep_alive_interval")
	if keepAliveIntervalConfig != 0 {
		keepAliveInterval = time.Duration(keepAliveIntervalConfig) * time.Second
	}

	srvConfig := MCPStreamableHTTPServerConfig{
		Address:           addr,
		Stateless:         stateless,
		BasePath:          basePath,
		KeepAlive:         keepAlive,
		KeepAliveInterval: keepAliveInterval,
	}

	srvOptions := []server.StreamableHTTPOption{
		server.WithStateLess(srvConfig.Stateless),
		server.WithEndpointPath(srvConfig.BasePath),
	}

	if srvConfig.KeepAlive {
		srvOptions = append(srvOptions, server.WithHeartbeatInterval(srvConfig.KeepAliveInterval))
	}

	srvOptions = append(srvOptions, options...)

	return NewMCPStreamableHTTPServer(mcpServer, srvConfig, srvOptions...)
}
