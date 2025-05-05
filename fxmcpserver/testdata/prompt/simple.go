package prompt

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SimpleTestPrompt struct{}

func NewSimpleTestPrompt() *SimpleTestPrompt {
	return &SimpleTestPrompt{}
}

func (p *SimpleTestPrompt) Name() string {
	return "simple-test-prompt"
}

func (p *SimpleTestPrompt) Options() []mcp.PromptOption {
	return []mcp.PromptOption{
		mcp.WithPromptDescription("Simple test prompt."),
	}
}

func (p *SimpleTestPrompt) Handle() server.PromptHandlerFunc {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return mcp.NewGetPromptResult(
			"ok",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent("simple test prompt"),
				),
			},
		), nil
	}
}
