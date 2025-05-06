package resource

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SimpleTestResource struct{}

func NewSimpleTestResource() *SimpleTestResource {
	return &SimpleTestResource{}
}

func (r *SimpleTestResource) Name() string {
	return "simple-test-resource"
}

func (r *SimpleTestResource) URI() string {
	return "simple-test://resources"
}

func (r *SimpleTestResource) Options() []mcp.ResourceOption {
	return []mcp.ResourceOption{
		mcp.WithResourceDescription("Simple test resource."),
	}
}

func (r *SimpleTestResource) Handle() server.ResourceHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     "simple test resource",
			},
		}, nil
	}
}
