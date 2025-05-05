package tool

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opencensus.io/trace"
)

type TestTool struct{}

func NewTestTool() *TestTool {
	return &TestTool{}
}

func (t *TestTool) Name() string {
	return "test-tool"
}

func (t *TestTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Test tool."),
	}
}

func (t *TestTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, span := trace.StartSpan(ctx, "TestTool.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("TestTool.Handle")

		return mcp.NewToolResultText("ok"), nil
	}
}
