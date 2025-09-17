package sse

import (
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
)

const (
	DefaultAddr              = ":8082"
	DefaultBaseURL           = ""
	DefaultBasePath          = ""
	DefaultSSEEndpoint       = "/sse"
	DefaultMessageEndpoint   = "/message"
	DefaultKeepAliveInterval = 10 * time.Second
)

var _ MCPSSEServerFactory = (*DefaultMCPSSEServerFactory)(nil)

// MCPSSEServerFactory is the interface for MCP SSE server factories.
type MCPSSEServerFactory interface {
	Create(mcpServer *server.MCPServer, options ...server.SSEOption) *MCPSSEServer
}

// DefaultMCPSSEServerFactory is the default MCPSSEServerFactory implementation.
type DefaultMCPSSEServerFactory struct {
	config *config.Config
}

// NewDefaultMCPSSEServerFactory returns a new DefaultMCPSSEServerFactory instance.
func NewDefaultMCPSSEServerFactory(config *config.Config) *DefaultMCPSSEServerFactory {
	return &DefaultMCPSSEServerFactory{
		config: config,
	}
}

// Create returns a new MCPSSEServer instance.
func (f *DefaultMCPSSEServerFactory) Create(mcpServer *server.MCPServer, options ...server.SSEOption) *MCPSSEServer {
	addr := f.config.GetString("modules.mcp.server.transport.sse.address")
	if addr == "" {
		addr = DefaultAddr
	}

	baseURL := f.config.GetString("modules.mcp.server.transport.sse.base_url")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	basePath := f.config.GetString("modules.mcp.server.transport.sse.base_path")
	if basePath == "" {
		basePath = DefaultBasePath
	}

	sseEndpoint := f.config.GetString("modules.mcp.server.transport.sse.sse_endpoint")
	if sseEndpoint == "" {
		sseEndpoint = DefaultSSEEndpoint
	}

	messageEndpoint := f.config.GetString("modules.mcp.server.transport.sse.message_endpoint")
	if messageEndpoint == "" {
		messageEndpoint = DefaultMessageEndpoint
	}

	keepAlive := f.config.GetBool("modules.mcp.server.transport.sse.keep_alive")

	keepAliveInterval := DefaultKeepAliveInterval
	keepAliveIntervalConfig := f.config.GetInt("modules.mcp.server.transport.sse.keep_alive_interval")
	if keepAliveIntervalConfig != 0 {
		keepAliveInterval = time.Duration(keepAliveIntervalConfig) * time.Second
	}

	srvConfig := MCPSSEServerConfig{
		Address:           addr,
		BaseURL:           baseURL,
		BasePath:          basePath,
		SSEEndpoint:       sseEndpoint,
		MessageEndpoint:   messageEndpoint,
		KeepAlive:         keepAlive,
		KeepAliveInterval: keepAliveInterval,
	}

	srvOptions := []server.SSEOption{
		server.WithBaseURL(srvConfig.BaseURL),
		server.WithStaticBasePath(srvConfig.BasePath),
		server.WithSSEEndpoint(srvConfig.SSEEndpoint),
		server.WithMessageEndpoint(srvConfig.MessageEndpoint),
	}

	if srvConfig.KeepAlive {
		srvOptions = append(srvOptions, server.WithKeepAliveInterval(srvConfig.KeepAliveInterval))
	}

	srvOptions = append(srvOptions, options...)

	return NewMCPSSEServer(mcpServer, srvConfig, srvOptions...)
}
