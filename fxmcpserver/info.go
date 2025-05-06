package fxmcpserver

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/fxmcpserver/server/stdio"
)

// MCPServerModuleInfo is the MCP server module info.
type MCPServerModuleInfo struct {
	config      *config.Config
	registry    *server.MCPServerRegistry
	sseServer   *sse.MCPSSEServer
	stdioServer *stdio.MCPStdioServer
}

// NewMCPServerModuleInfo returns a new MCPServerModuleInfo instance.
func NewMCPServerModuleInfo(
	config *config.Config,
	registry *server.MCPServerRegistry,
	sseServer *sse.MCPSSEServer,
	stdioServer *stdio.MCPStdioServer,
) *MCPServerModuleInfo {
	return &MCPServerModuleInfo{
		config:      config,
		registry:    registry,
		sseServer:   sseServer,
		stdioServer: stdioServer,
	}
}

// Name return the name of the module info.
func (i *MCPServerModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *MCPServerModuleInfo) Data() map[string]any {
	sseServerInfo := i.sseServer.Info()
	stdioServerInfo := i.stdioServer.Info()
	mcpRegistryInfo := i.registry.Info()

	return map[string]any{
		"transports": map[string]any{
			"sse":   sseServerInfo,
			"stdio": stdioServerInfo,
		},
		"capabilities": map[string]any{
			"tools":     mcpRegistryInfo.Capabilities.Tools,
			"prompts":   mcpRegistryInfo.Capabilities.Prompts,
			"resources": mcpRegistryInfo.Capabilities.Resources,
		},
		"registrations": map[string]any{
			"tools":             mcpRegistryInfo.Registrations.Tools,
			"prompts":           mcpRegistryInfo.Registrations.Prompts,
			"resources":         mcpRegistryInfo.Registrations.Resources,
			"resourceTemplates": mcpRegistryInfo.Registrations.ResourceTemplates,
		},
	}
}
