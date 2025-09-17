package server

import (
	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
)

const (
	DefaultServerName    = "MCP Server"
	DefaultServerVersion = "1.0.0"
)

var _ MCPServerFactory = (*DefaultMCPServerFactory)(nil)

// MCPServerFactory is the interface for server.MCPServer factories.
type MCPServerFactory interface {
	Create(options ...server.ServerOption) *server.MCPServer
}

// DefaultMCPServerFactory is the default MCPServerFactory implementation.
type DefaultMCPServerFactory struct {
	config *config.Config
}

// NewDefaultMCPServerFactory returns a new DefaultMCPServerFactory instance.
func NewDefaultMCPServerFactory(config *config.Config) *DefaultMCPServerFactory {
	return &DefaultMCPServerFactory{
		config: config,
	}
}

// Create creates a new server.MCPServer instance.
func (f *DefaultMCPServerFactory) Create(options ...server.ServerOption) *server.MCPServer {
	name := f.config.GetString("modules.mcp.server.name")
	if name == "" {
		name = DefaultServerName
	}

	version := f.config.GetString("modules.mcp.server.version")
	if version == "" {
		version = DefaultServerVersion
	}

	srvOptions := []server.ServerOption{
		server.WithLogging(),
		server.WithRecovery(),
	}

	instructions := f.config.GetString("modules.mcp.server.instructions")
	if instructions != "" {
		srvOptions = append(srvOptions, server.WithInstructions(instructions))
	}

	srvOptions = append(srvOptions, options...)

	return server.NewMCPServer(name, version, srvOptions...)
}
