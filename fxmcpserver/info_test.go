package fxmcpserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver"
	fs "github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/fxmcpserver/server/stdio"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/prompt"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resource"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/tool"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestMCPServerModuleInfo(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("testdata/config"),
	)
	assert.NoError(t, err)

	reg := fs.NewMCPServerRegistry(
		cfg,
		[]fs.MCPServerTool{
			tool.NewTestTool(),
		},
		[]fs.MCPServerPrompt{
			prompt.NewTestPrompt(),
		},
		[]fs.MCPServerResource{
			resource.NewTestResource(),
		},
		[]fs.MCPServerResourceTemplate{
			resource.NewTestResourceTemplate(),
		},
	)

	mcpSrv := server.NewMCPServer("test-server", "1.0.0")

	sseSrv := sse.NewDefaultMCPSSEServerFactory(cfg).Create(mcpSrv)
	stdioSrv := stdio.NewDefaultMCPStdioServerFactory().Create(mcpSrv)

	info := fxmcpserver.NewMCPServerModuleInfo(cfg, reg, sseSrv, stdioSrv)

	assert.Equal(t, info.Name(), fxmcpserver.ModuleName)

	expectedData := map[string]any{
		"transports": map[string]any{
			"sse": map[string]any{
				"config": map[string]any{
					"address":             ":0",
					"base_url":            sse.DefaultBaseURL,
					"base_path":           sse.DefaultBasePath,
					"sse_endpoint":        sse.DefaultSSEEndpoint,
					"message_endpoint":    sse.DefaultMessageEndpoint,
					"keep_alive":          true,
					"keep_alive_interval": sse.DefaultKeepAliveInterval.Seconds(),
				},
				"status": map[string]any{
					"running": false,
				},
			},
			"stdio": map[string]any{
				"status": map[string]any{
					"running": false,
				},
			},
		},
		"capabilities": map[string]any{
			"tools":     true,
			"prompts":   true,
			"resources": true,
		},
		"registrations": map[string]any{
			"tools": map[string]string{
				"test-tool": "github.com/ankorstore/yokai/fxmcpserver/testdata/tool.(*TestTool).Handle.func1",
			},
			"prompts": map[string]string{
				"test-prompt": "github.com/ankorstore/yokai/fxmcpserver/testdata/prompt.(*TestPrompt).Handle.func1",
			},
			"resources": map[string]string{
				"test-resource": "github.com/ankorstore/yokai/fxmcpserver/testdata/resource.(*TestResource).Handle.func1",
			},
			"resourceTemplates": map[string]string{
				"test-template": "github.com/ankorstore/yokai/fxmcpserver/testdata/resource.(*TestResourceTemplate).Handle.func1",
			},
		},
	}

	assert.Equal(t, expectedData, info.Data())
}
