package resourcetemplate

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SimpleTestResourceTemplate struct{}

func NewSimpleTestResourceTemplate() *SimpleTestResourceTemplate {
	return &SimpleTestResourceTemplate{}
}

func (r *SimpleTestResourceTemplate) Name() string {
	return "simple-test-resource-template"
}

func (r *SimpleTestResourceTemplate) URI() string {
	return "simple-test://resources/{id}"
}

func (r *SimpleTestResourceTemplate) Options() []mcp.ResourceTemplateOption {
	return []mcp.ResourceTemplateOption{
		mcp.WithTemplateDescription("Simple test resource template."),
	}
}

func (r *SimpleTestResourceTemplate) Handle() server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     "ok",
			},
		}, nil
	}
}
