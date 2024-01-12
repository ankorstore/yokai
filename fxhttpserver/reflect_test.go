package fxhttpserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target   any
		expected string
	}{
		{123, "int"},
		{"test", "string"},
		{echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}), "echo.MiddlewareFunc"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			got := fxhttpserver.GetType(tt.target)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetReturnType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target   any
		expected string
	}{
		{func() string { return "test" }, "string"},
		{func() int { return 123 }, "int"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			got := fxhttpserver.GetReturnType(tt.target)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestIsConcreteMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		middleware any
		expected   bool
	}{
		{echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc { return next }), true},
		{123, false},
		{"test", false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(fxhttpserver.GetType(tt.middleware), func(t *testing.T) {
			t.Parallel()

			got := fxhttpserver.IsConcreteMiddleware(tt.middleware)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestIsConcreteHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		handler  any
		expected bool
	}{
		{echo.HandlerFunc(func(c echo.Context) error { return nil }), true},
		{123, false},
		{"test", false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(fxhttpserver.GetType(tt.handler), func(t *testing.T) {
			t.Parallel()

			got := fxhttpserver.IsConcreteHandler(tt.handler)
			assert.Equal(t, tt.expected, got)
		})
	}
}
