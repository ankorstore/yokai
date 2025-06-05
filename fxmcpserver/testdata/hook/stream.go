package hook

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

type SimpleMCPStreamableHTTPServerContextHook struct{}

func NewSimpleMCPStreamableHTTPServerContextHook() *SimpleMCPStreamableHTTPServerContextHook {
	return &SimpleMCPStreamableHTTPServerContextHook{}
}

func (p *SimpleMCPStreamableHTTPServerContextHook) Handle() server.HTTPContextFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, "foo", "bar")
	}
}
