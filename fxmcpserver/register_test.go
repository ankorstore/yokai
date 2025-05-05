package fxmcpserver_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/prompt"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/resource"
	"github.com/ankorstore/yokai/fxmcpserver/testdata/tool"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestAsMCPServerTool(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerTool(tool.NewTestTool)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerTools(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerTools(tool.NewTestTool)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerPrompt(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerPrompt(prompt.NewTestPrompt)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerPrompts(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerPrompts(prompt.NewTestPrompt)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerResource(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerResource(resource.NewTestResource)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerResources(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerResources(resource.NewTestResource)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerResourceTemplate(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerResourceTemplate(resource.NewTestResourceTemplate)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}

func TestAsMCPServerResourceTemplates(t *testing.T) {
	t.Parallel()

	reg := fxmcpserver.AsMCPServerResourceTemplates(resource.NewTestResourceTemplate)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", reg))
	assert.Implements(t, (*fx.Option)(nil), reg)
}
