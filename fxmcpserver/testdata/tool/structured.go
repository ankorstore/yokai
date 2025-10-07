package tool

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type StructuredTestTool struct{}

func NewStructuredTestTool() *StructuredTestTool {
	return &StructuredTestTool{}
}

func (t *StructuredTestTool) Name() string {
	return "structured-test-tool"
}

func (t *StructuredTestTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Structured test tool."),
		mcp.WithInputSchema[StructuredTestToolRequest](),
		mcp.WithOutputSchema[StructuredTestToolResult](),
	}
}

type StructuredTestToolRequest struct {
	Input string `json:"input" jsonschema_description:"Test input" jsonschema:"required"`
}

type StructuredTestToolResult struct {
	Output string `json:"output" jsonschema_description:"Test output"`
}

func (t *StructuredTestTool) Handle() server.ToolHandlerFunc {
	return mcp.NewStructuredToolHandler(
		func(ctx context.Context, request mcp.CallToolRequest, args StructuredTestToolRequest) (StructuredTestToolResult, error) {
			output := fmt.Sprintf("input: %s", args.Input)

			return StructuredTestToolResult{
				Output: output,
			}, nil
		},
	)
}
