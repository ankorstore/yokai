package prompt

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opencensus.io/trace"
)

type TestPrompt struct{}

func NewTestPrompt() *TestPrompt {
	return &TestPrompt{}
}

func (p *TestPrompt) Name() string {
	return "test-prompt"
}

func (p *TestPrompt) Options() []mcp.PromptOption {
	return []mcp.PromptOption{
		mcp.WithPromptDescription("Test prompt."),
	}
}

func (p *TestPrompt) Handle() server.PromptHandlerFunc {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		ctx, span := trace.StartSpan(ctx, "TestPrompt.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("TestPrompt.Handle")

		return mcp.NewGetPromptResult(
			"ok",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent("test content"),
				),
			},
		), nil
	}
}
