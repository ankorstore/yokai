package stdio

import (
	"context"
	"io"
	"sync"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
)

// MCPStdioServerConfig is the MCP Stdio server configuration.
type MCPStdioServerConfig struct {
	In  io.Reader
	Out io.Writer
}

// MCPStdioServer is the MCP Stdio server.
type MCPStdioServer struct {
	server  *server.StdioServer
	config  MCPStdioServerConfig
	mutex   sync.RWMutex
	running bool
}

// NewMCPStdioServer returns a new MCPStdioServer instance.
func NewMCPStdioServer(mcpServer *server.MCPServer, config MCPStdioServerConfig, opts ...server.StdioOption) *MCPStdioServer {
	stdioServer := server.NewStdioServer(mcpServer)

	for _, opt := range opts {
		opt(stdioServer)
	}

	return &MCPStdioServer{
		server: stdioServer,
		config: config,
	}
}

// Server returns the MCPStdioServer underlying server.
func (s *MCPStdioServer) Server() *server.StdioServer {
	return s.server
}

// Config returns the MCPStdioServer config.
func (s *MCPStdioServer) Config() MCPStdioServerConfig {
	return s.config
}

// Start starts the MCPStdioServer.
func (s *MCPStdioServer) Start(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msg("starting MCP Stdio server")

	s.mutex.Lock()
	s.running = true
	s.mutex.Unlock()

	err := s.server.Listen(ctx, s.config.In, s.config.Out)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to start MCP Stdio server")

		s.running = false
	}

	return err
}

// Running returns true if the MCPStdioServer is running.
func (s *MCPStdioServer) Running() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.running
}

// Info returns the MCPStdioServer information.
func (s *MCPStdioServer) Info() map[string]any {
	return map[string]any{
		"status": map[string]any{
			"running": s.running,
		},
	}
}
