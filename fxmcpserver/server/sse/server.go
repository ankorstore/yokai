package sse

import (
	"context"
	"sync"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
)

// MCPSSEServerConfig is the MCP SSE server configuration.
type MCPSSEServerConfig struct {
	Address           string
	BaseURL           string
	BasePath          string
	SSEEndpoint       string
	MessageEndpoint   string
	KeepAlive         bool
	KeepAliveInterval time.Duration
}

// MCPSSEServer is the MCP SSE server.
type MCPSSEServer struct {
	server  *server.SSEServer
	config  MCPSSEServerConfig
	mutex   sync.RWMutex
	running bool
}

// NewMCPSSEServer returns a new MCPSSEServer instance.
func NewMCPSSEServer(mcpServer *server.MCPServer, config MCPSSEServerConfig, opts ...server.SSEOption) *MCPSSEServer {
	return &MCPSSEServer{
		server: server.NewSSEServer(mcpServer, opts...),
		config: config,
	}
}

// Server returns the MCPSSEServer underlying server.
func (s *MCPSSEServer) Server() *server.SSEServer {
	return s.server
}

// Config returns the MCPSSEServer config.
func (s *MCPSSEServer) Config() MCPSSEServerConfig {
	return s.config
}

// Start starts the MCPSSEServer.
func (s *MCPSSEServer) Start(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msgf("starting MCP SSE server on %s", s.config.Address)

	s.mutex.Lock()
	s.running = true
	s.mutex.Unlock()

	err := s.server.Start(s.config.Address)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to start MCP SSE server")

		s.mutex.Lock()
		s.running = false
		s.mutex.Unlock()
	}

	return err
}

// Stop stops the MCPSSEServer.
func (s *MCPSSEServer) Stop(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msg("stopping MCP SSE server")

	s.mutex.Lock()
	s.running = false
	s.mutex.Unlock()

	err := s.server.Shutdown(ctx)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to stop MCP SSE server")
	}

	return err
}

// Running returns true if the MCPSSEServer is running.
func (s *MCPSSEServer) Running() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.running
}

// Info returns the MCPSSEServer information.
func (s *MCPSSEServer) Info() map[string]any {
	return map[string]any{
		"config": map[string]any{
			"address":             s.config.Address,
			"base_url":            s.config.BaseURL,
			"base_path":           s.config.BasePath,
			"sse_endpoint":        s.config.SSEEndpoint,
			"message_endpoint":    s.config.MessageEndpoint,
			"keep_alive":          s.config.KeepAlive,
			"keep_alive_interval": s.config.KeepAliveInterval.Seconds(),
		},
		"status": map[string]any{
			"running": s.Running(),
		},
	}
}
