package server_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/prompt"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resource"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/tool"
	"github.com/stretchr/testify/assert"
)

func TestMCPServerRegistry_Info(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../testdata/config"),
	)
	assert.NoError(t, err)

	reg := server.NewMCPServerRegistry(
		cfg,
		[]server.MCPServerTool{
			tool.NewTestTool(),
		},
		[]server.MCPServerPrompt{
			prompt.NewTestPrompt(),
		},
		[]server.MCPServerResource{
			resource.NewTestResource(),
		},
		[]server.MCPServerResourceTemplate{
			resource.NewTestResourceTemplate(),
		},
	)

	expectedInfo := server.MCPServerRegistryInfo{
		Capabilities: struct {
			Tools     bool
			Prompts   bool
			Resources bool
		}{
			Tools:     true,
			Prompts:   true,
			Resources: true,
		},
		Registrations: struct {
			Tools             map[string]string
			Prompts           map[string]string
			Resources         map[string]string
			ResourceTemplates map[string]string
		}{
			Tools: map[string]string{
				"test-tool": "github.com/ankorstore/yokai/fxmcpserver/testdata/tool.(*TestTool).Handle.func1",
			},
			Prompts: map[string]string{
				"test-prompt": "github.com/ankorstore/yokai/fxmcpserver/testdata/prompt.(*TestPrompt).Handle.func1",
			},
			Resources: map[string]string{
				"test-resource": "github.com/ankorstore/yokai/fxmcpserver/testdata/resource.(*TestResource).Handle.func1",
			},
			ResourceTemplates: map[string]string{
				"test-template": "github.com/ankorstore/yokai/fxmcpserver/testdata/resource.(*TestResourceTemplate).Handle.func1",
			},
		},
	}

	assert.Equal(t, expectedInfo, reg.Info())
}
