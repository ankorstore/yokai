package fxmcpservertest

import (
	"context"
	"net/http/httptest"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPSSETestServer struct {
	config     *config.Config
	testServer *httptest.Server
}

func NewMCPSSETestServer(cfg *config.Config, srv *server.MCPServer, hdl sse.MCPSSEServerContextHandler) *MCPSSETestServer {
	sseEndpoint := cfg.GetString("modules.mcp.server.transport.sse.sse_endpoint")
	if sseEndpoint == "" {
		sseEndpoint = sse.DefaultSSEEndpoint
	}

	messageEndpoint := cfg.GetString("modules.mcp.server.transport.sse.message_endpoint")
	if messageEndpoint == "" {
		messageEndpoint = sse.DefaultMessageEndpoint
	}

	testSrv := server.NewTestServer(
		srv,
		server.WithSSEContextFunc(hdl.Handle()),
		server.WithSSEEndpoint(sseEndpoint),
		server.WithMessageEndpoint(messageEndpoint),
	)

	return &MCPSSETestServer{
		config:     cfg,
		testServer: testSrv,
	}
}

func (s *MCPSSETestServer) Close() {
	s.testServer.Close()
}

func (s *MCPSSETestServer) StartClient(ctx context.Context, options ...transport.ClientOption) (*client.Client, error) {
	sseEndpoint := s.config.GetString("modules.mcp.server.transport.sse.sse_endpoint")
	if sseEndpoint == "" {
		sseEndpoint = sse.DefaultSSEEndpoint
	}

	baseURL := s.testServer.URL + sseEndpoint

	cli, err := client.NewSSEMCPClient(baseURL, options...)
	if err != nil {
		return nil, err
	}

	err = cli.Start(ctx)
	if err != nil {
		return nil, err
	}

	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initReq)
	if err != nil {
		return nil, err
	}

	return cli, nil
}
