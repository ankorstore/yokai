package fxhttpserver

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

// GetType returns the type of a target.
func GetType(target any) string {
	return reflect.TypeOf(target).String()
}

// GetReturnType returns the return type of a target.
func GetReturnType(target any) string {
	return reflect.TypeOf(target).Out(0).String()
}

// IsConcreteMiddleware returns true if the middleware is a concrete [echo.MiddlewareFunc] implementation.
func IsConcreteMiddleware(middleware any) bool {
	return reflect.TypeOf(middleware).ConvertibleTo(reflect.TypeOf(echo.MiddlewareFunc(nil)))
}

// IsConcreteHandler returns true if the handler is a concrete [echo.HandlerFunc] implementation.
func IsConcreteHandler(handler any) bool {
	return reflect.TypeOf(handler).ConvertibleTo(reflect.TypeOf(echo.HandlerFunc(nil)))
}
