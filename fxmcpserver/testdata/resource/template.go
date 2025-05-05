package resource

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opencensus.io/trace"
)

type TestResourceTemplate struct{}

func NewTestResourceTemplate() *TestResourceTemplate {
	return &TestResourceTemplate{}
}

func (r *TestResourceTemplate) Name() string {
	return "test-template"
}

func (r *TestResourceTemplate) URI() string {
	return "test://resources/{id}"
}

func (r *TestResourceTemplate) Options() []mcp.ResourceTemplateOption {
	return []mcp.ResourceTemplateOption{
		mcp.WithTemplateDescription("Test resource template."),
	}
}

func (r *TestResourceTemplate) Handle() server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		ctx, span := trace.StartSpan(ctx, "TestResourceTemplate.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("TestResourceTemplate.Handle")

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     "ok",
			},
		}, nil
	}
}
