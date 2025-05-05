package tool

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SimpleTestTool struct{}

func NewSimpleTestTool() *SimpleTestTool {
	return &SimpleTestTool{}
}

func (t *SimpleTestTool) Name() string {
	return "simple-test-tool"
}

func (t *SimpleTestTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Simple test tool."),
	}
}

func (t *SimpleTestTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("ok"), nil
	}
}
