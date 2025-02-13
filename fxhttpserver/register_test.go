package fxhttpserver_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxhttpserver/testdata/errorhandler"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareRegistration(t *testing.T) {
	t.Parallel()

	type exampleMiddleware struct {
		name string
	}

	mw := exampleMiddleware{name: "test"}

	tests := []struct {
		middleware any
		kind       fxhttpserver.MiddlewareKind
	}{
		{mw, fxhttpserver.GlobalUse},
		{mw, fxhttpserver.GlobalPre},
		{mw, fxhttpserver.Attached},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.kind.String(), func(t *testing.T) {
			t.Parallel()

			mr := fxhttpserver.NewMiddlewareRegistration(tt.middleware, tt.kind)

			assert.Equal(t, tt.middleware, mr.Middleware())
			assert.Equal(t, tt.kind, mr.Kind())
		})
	}
}

func TestHandlerRegistration(t *testing.T) {
	t.Parallel()

	type exampleHandler struct {
		name string
	}
	type exampleMiddleware struct {
		name string
	}

	handler := exampleHandler{name: "handler-test"}
	mw1 := exampleMiddleware{name: "middleware1"}
	mw2 := exampleMiddleware{name: "middleware2"}

	tests := []struct {
		method      string
		path        string
		handler     any
		middlewares []any
	}{
		{"GET", "/path1", handler, nil},
		{"POST", "/path2", handler, []any{mw1}},
		{"PUT", "/path3", handler, []any{mw1, mw2}},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.method+tt.path, func(t *testing.T) {
			t.Parallel()

			hr := fxhttpserver.NewHandlerRegistration(tt.method, tt.path, tt.handler, tt.middlewares...)

			assert.Equal(t, tt.method, hr.Method())
			assert.Equal(t, tt.path, hr.Path())
			assert.Equal(t, tt.handler, hr.Handler())
			assert.Equal(t, tt.middlewares, hr.Middlewares())
		})
	}
}

func TestHandlersGroupRegistration(t *testing.T) {
	t.Parallel()

	type exampleHandler struct {
		name string
	}
	type exampleMiddleware struct {
		name string
	}

	handler1 := exampleHandler{name: "handler-test1"}
	handler2 := exampleHandler{name: "handler-test2"}
	mw1 := exampleMiddleware{name: "middleware1"}
	mw2 := exampleMiddleware{name: "middleware2"}

	hr1 := fxhttpserver.NewHandlerRegistration("GET", "/path1", handler1)
	hr2 := fxhttpserver.NewHandlerRegistration("POST", "/path2", handler2)

	tests := []struct {
		prefix                string
		handlersRegistrations []*fxhttpserver.HandlerRegistration
		middlewares           []any
	}{
		{"/group/1", []*fxhttpserver.HandlerRegistration{hr1}, nil},
		{"/group/2", []*fxhttpserver.HandlerRegistration{hr1, hr2}, []any{mw1}},
		{"/group/3", []*fxhttpserver.HandlerRegistration{hr1}, []any{mw1, mw2}},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.prefix, func(t *testing.T) {
			t.Parallel()

			hgr := fxhttpserver.NewHandlersGroupRegistration(tt.prefix, tt.handlersRegistrations, tt.middlewares...)

			assert.Equal(t, tt.prefix, hgr.Prefix())
			assert.Equal(t, tt.handlersRegistrations, hgr.HandlersRegistrations())
			assert.Equal(t, tt.middlewares, hgr.Middlewares())
		})
	}
}

func TestErrorHandlerRegistration(t *testing.T) {
	t.Parallel()

	result := fxhttpserver.AsErrorHandler(errorhandler.NewTestErrorHandler)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", result))
}
