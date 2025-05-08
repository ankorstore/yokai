package middleware

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

type SimpleMCPSSEServerMiddleware struct{}

func NewSimpleMCPSSEServerMiddleware() *SimpleMCPSSEServerMiddleware {
	return &SimpleMCPSSEServerMiddleware{}
}

func (p *SimpleMCPSSEServerMiddleware) Handle() server.SSEContextFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, "foo", "bar")
	}
}
