package fxmcpserver_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/hook"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/prompt"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resource"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resourcetemplate"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/tool"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestAsMCPServerTool(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerTool(tool.NewSimpleTestTool)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerTools(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerTools(tool.NewSimpleTestTool)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerPrompt(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerPrompt(prompt.NewSimpleTestPrompt)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerPrompts(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerPrompts(prompt.NewSimpleTestPrompt)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerResource(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerResource(resource.NewSimpleTestResource)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerResources(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerResources(resource.NewSimpleTestResource)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerResourceTemplate(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerResourceTemplate(resourcetemplate.NewSimpleTestResourceTemplate)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerResourceTemplates(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerResourceTemplates(resourcetemplate.NewSimpleTestResourceTemplate)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPSSEServerContextHook(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPSSEServerContextHook(hook.NewSimpleMCPSSEServerContextHook)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPSSEServerContextHooks(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPSSEServerContextHooks(hook.NewSimpleMCPSSEServerContextHook)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPStreamableHTTPServerContextHook(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPStreamableHTTPServerContextHook(hook.NewSimpleMCPStreamableHTTPServerContextHook)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPStreamableHTTPServerContextHooks(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPStreamableHTTPServerContextHooks(hook.NewSimpleMCPStreamableHTTPServerContextHook)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}
