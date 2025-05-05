package server

import (
	"context"
	"strings"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/fxmcpserver/server/stdio"
	"github.com/ankorstore/yokai/healthcheck"
)

// MCPServerProbe is a probe compatible with the healthcheck module.
type MCPServerProbe struct {
	config      *config.Config
	sseServer   *sse.MCPSSEServer
	stdioServer *stdio.MCPStdioServer
}

// NewMCPServerProbe returns a new MCPServerProbe.
func NewMCPServerProbe(
	config *config.Config,
	sseServer *sse.MCPSSEServer,
	stdioServer *stdio.MCPStdioServer,
) *MCPServerProbe {
	return &MCPServerProbe{
		config:      config,
		sseServer:   sseServer,
		stdioServer: stdioServer,
	}
}

// Name returns the name of the MCPServerProbe.
func (p *MCPServerProbe) Name() string {
	return "mcpserver"
}

// Check returns a successful healthcheck.CheckerProbeResult if the exposed MCP servers are running.
func (p *MCPServerProbe) Check(context.Context) *healthcheck.CheckerProbeResult {
	success := true
	var messages []string

	if p.config.GetBool("modules.mcp.server.transport.sse.expose") {
		if p.sseServer.Running() {
			messages = append(messages, "MCP SSE server is running")
		} else {
			success = false
			messages = append(messages, "MCP SSE server is not running")
		}
	}

	if p.config.GetBool("modules.mcp.server.transport.stdio.expose") {
		if p.stdioServer.Running() {
			messages = append(messages, "MCP Stdio server is running")
		} else {
			success = false
			messages = append(messages, "MCP Stdio server is not running")
		}
	}

	return &healthcheck.CheckerProbeResult{
		Success: success,
		Message: strings.Join(messages, ", "),
	}
}
