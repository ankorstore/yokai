package stream

import (
	"context"
	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
	"sync"
	"time"
)

// MCPStreamableHTTPServerConfig is the MCP StreamableHTTP server configuration.
type MCPStreamableHTTPServerConfig struct {
	Address           string
	Stateless         bool
	BasePath          string
	KeepAlive         bool
	KeepAliveInterval time.Duration
}

// MCPStreamableHTTPServer is the MCP StreamableHTTP server.
type MCPStreamableHTTPServer struct {
	server  *server.StreamableHTTPServer
	config  MCPStreamableHTTPServerConfig
	mutex   sync.RWMutex
	running bool
}

// NewMCPStreamableHTTPServer returns a new MCPStreamableHTTPServer instance.
func NewMCPStreamableHTTPServer(mcpServer *server.MCPServer, config MCPStreamableHTTPServerConfig, opts ...server.StreamableHTTPOption) *MCPStreamableHTTPServer {
	streamableHTTPServer := server.NewStreamableHTTPServer(mcpServer, opts...)

	return &MCPStreamableHTTPServer{
		server: streamableHTTPServer,
		config: config,
	}
}

// Server returns the MCPStreamableHTTPServer underlying server.
func (s *MCPStreamableHTTPServer) Server() *server.StreamableHTTPServer {
	return s.server
}

// Config returns the MCPStreamableHTTPServer config.
func (s *MCPStreamableHTTPServer) Config() MCPStreamableHTTPServerConfig {
	return s.config
}

// Start starts the MCPStreamableHTTPServer.
func (s *MCPStreamableHTTPServer) Start(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msgf("starting MCP StreamableHTTP server on %s", s.config.Address)

	s.mutex.Lock()
	s.running = true
	s.mutex.Unlock()

	err := s.server.Start(s.config.Address)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to start MCP StreamableHTTP server")

		s.mutex.Lock()
		s.running = false
		s.mutex.Unlock()
	}

	return err
}

// Stop stops the MCPStreamableHTTPServer.
func (s *MCPStreamableHTTPServer) Stop(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msg("stopping MCP StreamableHTTP server")

	s.mutex.Lock()
	s.running = false
	s.mutex.Unlock()

	err := s.server.Shutdown(ctx)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to stop MCP StreamableHTTP server")
	}

	return err
}

// Running returns true if the MCPStreamableHTTPServer is running.
func (s *MCPStreamableHTTPServer) Running() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.running
}

// Info returns the MCPStreamableHTTPServer information.
func (s *MCPStreamableHTTPServer) Info() map[string]any {
	return map[string]any{
		"config": map[string]any{
			"address":             s.config.Address,
			"stateless":           s.config.Stateless,
			"base_path":           s.config.BasePath,
			"keep_alive":          s.config.KeepAlive,
			"keep_alive_interval": s.config.KeepAliveInterval.Seconds(),
		},
		"status": map[string]any{
			"running": s.Running(),
		},
	}
}
