package hook

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

type SimpleMCPSSEServerContextHook struct{}

func NewSimpleMCPSSEServerContextHook() *SimpleMCPSSEServerContextHook {
	return &SimpleMCPSSEServerContextHook{}
}

func (p *SimpleMCPSSEServerContextHook) Handle() server.SSEContextFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, "foo", "bar")
	}
}
