package tool

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opencensus.io/trace"
)

type AdvancedTestTool struct {
	config *config.Config
}

func NewAdvancedTestTool(config *config.Config) *AdvancedTestTool {
	return &AdvancedTestTool{
		config: config,
	}
}

func (t *AdvancedTestTool) Name() string {
	return "advanced-test-tool"
}

func (t *AdvancedTestTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Advanced test tool."),
		mcp.WithBoolean(
			"shouldFail",
			mcp.Description("If the tool call should fail or not."),
		),
	}
}

func (t *AdvancedTestTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, span := trace.StartSpan(ctx, "AdvancedTestTool.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("AdvancedTestTool.Handle")

		shouldFail := request.Params.Arguments["shouldFail"].(string)
		if shouldFail == "true" {
			return nil, fmt.Errorf("advanced tool test failure")
		}

		return mcp.NewToolResultText(t.config.AppName()), nil
	}
}
