package resource

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opencensus.io/trace"
)

type TestResource struct{}

func NewTestResource() *TestResource {
	return &TestResource{}
}

func (r *TestResource) Name() string {
	return "test-resource"
}

func (r *TestResource) URI() string {
	return "test://resources"
}

func (r *TestResource) Options() []mcp.ResourceOption {
	return []mcp.ResourceOption{
		mcp.WithResourceDescription("Test resource."),
	}
}

func (r *TestResource) Handle() server.ResourceHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		ctx, span := trace.StartSpan(ctx, "TestResource.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("TestResource.Handle")

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     "ok",
			},
		}, nil
	}
}
