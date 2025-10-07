package tool

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type TypedTestTool struct{}

func NewTypedTestTool() *TypedTestTool {
	return &TypedTestTool{}
}

func (t *TypedTestTool) Name() string {
	return "typed-test-tool"
}

func (t *TypedTestTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Typed test tool."),
	}
}

type TypedTestToolRequest struct {
	Input string `json:"input" jsonschema_description:"Test input" jsonschema:"required"`
}

func (t *TypedTestTool) Handle() server.ToolHandlerFunc {
	return mcp.NewTypedToolHandler(
		func(ctx context.Context, request mcp.CallToolRequest, args TypedTestToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText(fmt.Sprintf("input: %s", args.Input)), nil
		},
	)
}
