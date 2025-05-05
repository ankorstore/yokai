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
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resourcetemplate"
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
			tool.NewSimpleTestTool(),
		},
		[]fs.MCPServerPrompt{
			prompt.NewSimpleTestPrompt(),
		},
		[]fs.MCPServerResource{
			resource.NewSimpleTestResource(),
		},
		[]fs.MCPServerResourceTemplate{
			resourcetemplate.NewSimpleTestResourceTemplate(),
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
				"simple-test-tool": "github.com/ankorstore/yokai/fxmcpserver/testdata/tool.(*SimpleTestTool).Handle.func1",
			},
			"prompts": map[string]string{
				"simple-test-prompt": "github.com/ankorstore/yokai/fxmcpserver/testdata/prompt.(*SimpleTestPrompt).Handle.func1",
			},
			"resources": map[string]string{
				"simple-test-resource": "github.com/ankorstore/yokai/fxmcpserver/testdata/resource.(*SimpleTestResource).Handle.func1",
			},
			"resourceTemplates": map[string]string{
				"simple-test-resource-template": "github.com/ankorstore/yokai/fxmcpserver/testdata/resourcetemplate.(*SimpleTestResourceTemplate).Handle.func1",
			},
		},
	}

	assert.Equal(t, expectedData, info.Data())
}
