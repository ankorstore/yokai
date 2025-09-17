package fxhttpserver

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

// fullTypeID builds a stable identifier for a type in the form "<pkgpath>.<TypeName>".
func fullTypeID(t reflect.Type) string {
	if t == nil {
		return ""
	}

	// Unwrap pointers to get the underlying named type (if any).
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// For named types, PkgPath() + Name() gives a unique and stable identity.
	if t.Name() != "" && t.PkgPath() != "" {
		return t.PkgPath() + "." + t.Name()
	}

	// Fallback for non-named kinds (slices, maps, func, etc.).
	return t.String()
}

// GetType returns a stable identifier for the given target’s type.
func GetType(target any) string {
	return fullTypeID(reflect.TypeOf(target))
}

// GetReturnType returns a stable identifier for the return type of constructor-like target.
// If a target is a function, we examine its first return value (index 0), unwrap pointers, and
// build an identifier for that named type. For non-function or empty-return cases, we return "".
func GetReturnType(target any) string {
	t := reflect.TypeOf(target)
	if t == nil || t.Kind() != reflect.Func || t.NumOut() == 0 {
		return ""
	}

	return fullTypeID(t.Out(0))
}

// IsConcreteMiddleware returns true if the middleware is a concrete [echo.MiddlewareFunc] implementation.
func IsConcreteMiddleware(middleware any) bool {
	return reflect.TypeOf(middleware).ConvertibleTo(reflect.TypeOf(echo.MiddlewareFunc(nil)))
}

// IsConcreteHandler returns true if the handler is a concrete [echo.HandlerFunc] implementation.
func IsConcreteHandler(handler any) bool {
	return reflect.TypeOf(handler).ConvertibleTo(reflect.TypeOf(echo.HandlerFunc(nil)))
}
