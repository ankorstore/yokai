package server_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/prompt"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resource"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resourcetemplate"
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
			tool.NewSimpleTestTool(),
		},
		[]server.MCPServerPrompt{
			prompt.NewSimpleTestPrompt(),
		},
		[]server.MCPServerResource{
			resource.NewSimpleTestResource(),
		},
		[]server.MCPServerResourceTemplate{
			resourcetemplate.NewSimpleTestResourceTemplate(),
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
				"simple-test-tool": "github.com/ankorstore/yokai/fxmcpserver/testdata/tool.(*SimpleTestTool).Handle.func1",
			},
			Prompts: map[string]string{
				"simple-test-prompt": "github.com/ankorstore/yokai/fxmcpserver/testdata/prompt.(*SimpleTestPrompt).Handle.func1",
			},
			Resources: map[string]string{
				"simple-test-resource": "github.com/ankorstore/yokai/fxmcpserver/testdata/resource.(*SimpleTestResource).Handle.func1",
			},
			ResourceTemplates: map[string]string{
				"simple-test-resource-template": "github.com/ankorstore/yokai/fxmcpserver/testdata/resourcetemplate.(*SimpleTestResourceTemplate).Handle.func1",
			},
		},
	}

	assert.Equal(t, expectedInfo, reg.Info())
}
