package fxhttpserver_test

import (
	"net/http"
	"testing"

	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxhttpserver/testdata/handler"
	"github.com/ankorstore/yokai/fxhttpserver/testdata/middleware"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareDefinition(t *testing.T) {
	t.Parallel()

	kind := fxhttpserver.GlobalUse

	md := fxhttpserver.NewMiddlewareDefinition(middleware.NewTestGlobalMiddleware, kind)

	assert.False(t, md.Concrete())
	assert.Equal(t, kind, md.Kind())
}

func TestHandlerDefinition(t *testing.T) {
	t.Parallel()

	method := http.MethodGet
	path := "/test"
	hand := handler.NewTestBarHandler
	middlewares := []fxhttpserver.MiddlewareDefinition{
		fxhttpserver.NewMiddlewareDefinition(middleware.NewTestGlobalMiddleware, fxhttpserver.GlobalUse),
	}

	hd := fxhttpserver.NewHandlerDefinition(method, path, hand, middlewares)

	assert.False(t, hd.Concrete())
	assert.Equal(t, method, hd.Method())
	assert.Equal(t, path, hd.Path())
	assert.Equal(t, middlewares, hd.Middlewares())
}

func TestHandlersGroupDefinition(t *testing.T) {
	t.Parallel()

	prefix := "/group"
	hand := handler.NewTestBarHandler
	handlers := []fxhttpserver.HandlerDefinition{
		fxhttpserver.NewHandlerDefinition(http.MethodGet, "/test", hand, nil),
	}
	middlewares := []fxhttpserver.MiddlewareDefinition{
		fxhttpserver.NewMiddlewareDefinition(middleware.NewTestGlobalMiddleware, fxhttpserver.GlobalUse),
	}

	hgd := fxhttpserver.NewHandlersGroupDefinition(prefix, handlers, middlewares)

	assert.Equal(t, prefix, hgd.Prefix())
	assert.Equal(t, handlers, hgd.Handlers())
	assert.Equal(t, middlewares, hgd.Middlewares())
}
