package fxmcpservertest

import (
	"context"
	"github.com/ankorstore/yokai/fxmcpserver/server/stream"
	"net/http/httptest"

	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPStreamableHTTPTestServer struct {
	config     *config.Config
	testServer *httptest.Server
}

func NewMCPStreamableHTTPTestServer(cfg *config.Config, srv *server.MCPServer, hdl stream.MCPStreamableHTTPServerContextHandler) *MCPStreamableHTTPTestServer {
	basePath := cfg.GetString("modules.mcp.server.transport.stream.base_path")
	if basePath == "" {
		basePath = stream.DefaultBasePath
	}

	testSrv := server.NewTestStreamableHTTPServer(
		srv,
		server.WithHTTPContextFunc(hdl.Handle()),
		server.WithEndpointPath(basePath),
	)

	return &MCPStreamableHTTPTestServer{
		config:     cfg,
		testServer: testSrv,
	}
}

func (s *MCPStreamableHTTPTestServer) Close() {
	s.testServer.Close()
}

func (s *MCPStreamableHTTPTestServer) StartClient(ctx context.Context, options ...transport.StreamableHTTPCOption) (*client.Client, error) {
	basePath := s.config.GetString("modules.mcp.server.transport.stream.base_path")
	if basePath == "" {
		basePath = stream.DefaultBasePath
	}

	baseURL := s.testServer.URL + basePath

	cli, err := client.NewStreamableHttpClient(baseURL, options...)
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
