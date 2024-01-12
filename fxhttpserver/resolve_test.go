package fxhttpserver_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var testMiddlewareFunc = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

var testHandlerFunc = func(next echo.Context) error {
	return fmt.Errorf("custom error")
}

func TestResolvedMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		middleware echo.MiddlewareFunc
		kind       fxhttpserver.MiddlewareKind
	}{
		{testMiddlewareFunc, fxhttpserver.GlobalUse},
		{testMiddlewareFunc, fxhttpserver.GlobalPre},
		{testMiddlewareFunc, fxhttpserver.Attached},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.kind.String(), func(t *testing.T) {
			t.Parallel()

			rm := fxhttpserver.NewResolvedMiddleware(tt.middleware, tt.kind)

			assert.Equal(t, "custom error", rm.Middleware()(testHandlerFunc)(nil).Error())
			assert.Equal(t, tt.kind, rm.Kind())
		})
	}
}

func TestResolvedHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method      string
		path        string
		handler     echo.HandlerFunc
		middlewares []echo.MiddlewareFunc
	}{
		{"GET", "/path1", testHandlerFunc, nil},
		{"POST", "/path2", testHandlerFunc, []echo.MiddlewareFunc{testMiddlewareFunc}},
		{"PUT", "/path3", testHandlerFunc, []echo.MiddlewareFunc{testMiddlewareFunc, testMiddlewareFunc}},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.method+tt.path, func(t *testing.T) {
			t.Parallel()

			rh := fxhttpserver.NewResolvedHandler(tt.method, tt.path, tt.handler, tt.middlewares...)

			assert.Equal(t, tt.method, rh.Method())
			assert.Equal(t, tt.path, rh.Path())
			assert.Equal(t, "custom error", rh.Handler()(nil).Error())
			assert.Equal(t, tt.middlewares, rh.Middlewares())
		})
	}
}

func TestResolvedHandlersGroup(t *testing.T) {
	t.Parallel()

	rh1 := fxhttpserver.NewResolvedHandler("GET", "/path1", testHandlerFunc)
	rh2 := fxhttpserver.NewResolvedHandler("POST", "/path2", testHandlerFunc, testMiddlewareFunc)

	tests := []struct {
		prefix      string
		handlers    []fxhttpserver.ResolvedHandler
		middlewares []echo.MiddlewareFunc
	}{
		{"/group/1", []fxhttpserver.ResolvedHandler{rh1}, nil},
		{"/group/2", []fxhttpserver.ResolvedHandler{rh1, rh2}, []echo.MiddlewareFunc{testMiddlewareFunc}},
		{"/group/3", []fxhttpserver.ResolvedHandler{rh1}, []echo.MiddlewareFunc{testMiddlewareFunc, testMiddlewareFunc}},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.prefix, func(t *testing.T) {
			t.Parallel()

			rg := fxhttpserver.NewResolvedHandlersGroup(tt.prefix, tt.handlers, tt.middlewares...)

			assert.Equal(t, tt.prefix, rg.Prefix())
			assert.Equal(t, tt.handlers, rg.Handlers())
			assert.Equal(t, tt.middlewares, rg.Middlewares())
		})
	}
}
