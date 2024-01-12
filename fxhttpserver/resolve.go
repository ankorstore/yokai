package fxhttpserver

import "github.com/labstack/echo/v4"

// ResolvedMiddleware is an interface for the resolved middlewares.
type ResolvedMiddleware interface {
	Middleware() echo.MiddlewareFunc
	Kind() MiddlewareKind
}

type resolvedMiddleware struct {
	middleware echo.MiddlewareFunc
	kind       MiddlewareKind
}

// NewResolvedMiddleware returns a new [ResolvedMiddleware].
func NewResolvedMiddleware(middleware echo.MiddlewareFunc, kind MiddlewareKind) ResolvedMiddleware {
	return &resolvedMiddleware{
		middleware: middleware,
		kind:       kind,
	}
}

// Middleware return the resolved middleware as [echo.MiddlewareFunc].
func (r *resolvedMiddleware) Middleware() echo.MiddlewareFunc {
	return r.middleware
}

// Kind return the resolved middleware kind.
func (r *resolvedMiddleware) Kind() MiddlewareKind {
	return r.kind
}

// ResolvedHandler is an interface for the resolved handlers.
type ResolvedHandler interface {
	Method() string
	Path() string
	Handler() echo.HandlerFunc
	Middlewares() []echo.MiddlewareFunc
}

type resolvedHandler struct {
	method      string
	path        string
	handler     echo.HandlerFunc
	middlewares []echo.MiddlewareFunc
}

// NewResolvedHandler returns a new [ResolvedHandler].
func NewResolvedHandler(method string, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) ResolvedHandler {
	return &resolvedHandler{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

// Method return the resolved handler http method.
func (r *resolvedHandler) Method() string {
	return r.method
}

// Path return the resolved handler http path.
func (r *resolvedHandler) Path() string {
	return r.path
}

// Handler return the resolved handler as [echo.HandlerFunc].
func (r *resolvedHandler) Handler() echo.HandlerFunc {
	return r.handler
}

// Middlewares return the resolved handler associated middlewares as a list of [echo.MiddlewareFunc].
func (r *resolvedHandler) Middlewares() []echo.MiddlewareFunc {
	return r.middlewares
}

// ResolvedHandlersGroup is an interface for the resolved handlers groups.
type ResolvedHandlersGroup interface {
	Prefix() string
	Handlers() []ResolvedHandler
	Middlewares() []echo.MiddlewareFunc
}

type resolvedHandlersGroup struct {
	prefix      string
	handlers    []ResolvedHandler
	middlewares []echo.MiddlewareFunc
}

// NewResolvedHandlersGroup returns a new [ResolvedHandlersGroup].
func NewResolvedHandlersGroup(prefix string, handlers []ResolvedHandler, middlewares ...echo.MiddlewareFunc) ResolvedHandlersGroup {
	return &resolvedHandlersGroup{
		prefix:      prefix,
		handlers:    handlers,
		middlewares: middlewares,
	}
}

// Prefix return the resolved handlers group http path prefix.
func (r *resolvedHandlersGroup) Prefix() string {
	return r.prefix
}

// Handlers return the resolved handlers group associated handlers as a list of [ResolvedHandler].
func (r *resolvedHandlersGroup) Handlers() []ResolvedHandler {
	return r.handlers
}

// Middlewares return the resolved handlers group associated middlewares as a list of [echo.MiddlewareFunc].
func (r *resolvedHandlersGroup) Middlewares() []echo.MiddlewareFunc {
	return r.middlewares
}
