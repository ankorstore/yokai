package fxhttpserver

// MiddlewareDefinition is the interface for middlewares definitions.
type MiddlewareDefinition interface {
	Concrete() bool
	Middleware() any
	Kind() MiddlewareKind
}

type middlewareDefinition struct {
	middleware any
	kind       MiddlewareKind
}

// NewMiddlewareDefinition returns a new [MiddlewareDefinition].
func NewMiddlewareDefinition(middleware any, kind MiddlewareKind) MiddlewareDefinition {
	return &middlewareDefinition{
		middleware: middleware,
		kind:       kind,
	}
}

// Concrete returns true if the middleware is a [echo.MiddlewareFunc] concrete implementation.
func (d *middlewareDefinition) Concrete() bool {
	return IsConcreteMiddleware(d.middleware)
}

// Middleware returns the middleware.
func (d *middlewareDefinition) Middleware() any {
	return d.middleware
}

// Kind returns the middleware kind.
func (d *middlewareDefinition) Kind() MiddlewareKind {
	return d.kind
}

// HandlerDefinition is the interface for handlers definitions.
type HandlerDefinition interface {
	Concrete() bool
	Method() string
	Path() string
	Handler() any
	Middlewares() []MiddlewareDefinition
}

type handlerDefinition struct {
	method      string
	path        string
	handler     any
	middlewares []MiddlewareDefinition
}

// NewHandlerDefinition returns a new [HandlerDefinition].
func NewHandlerDefinition(method string, path string, handler any, middlewares []MiddlewareDefinition) HandlerDefinition {
	return &handlerDefinition{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

// Concrete returns true if the handler is a [echo.HandlerFunc] concrete implementation.
func (d *handlerDefinition) Concrete() bool {
	return IsConcreteHandler(d.handler)
}

// Method returns the handler http method.
func (d *handlerDefinition) Method() string {
	return d.method
}

// Path returns the handler http path.
func (d *handlerDefinition) Path() string {
	return d.path
}

// Handler returns the handler.
func (d *handlerDefinition) Handler() any {
	return d.handler
}

// Middlewares returns the handler associated middlewares.
func (d *handlerDefinition) Middlewares() []MiddlewareDefinition {
	return d.middlewares
}

// HandlersGroupDefinition is the interface for handlers groups definitions.
type HandlersGroupDefinition interface {
	Prefix() string
	Handlers() []HandlerDefinition
	Middlewares() []MiddlewareDefinition
}

type handlersGroupDefinition struct {
	prefix      string
	handlers    []HandlerDefinition
	middlewares []MiddlewareDefinition
}

// NewHandlersGroupDefinition returns a new [HandlersGroupDefinition].
func NewHandlersGroupDefinition(prefix string, handlers []HandlerDefinition, middlewares []MiddlewareDefinition) HandlersGroupDefinition {
	return &handlersGroupDefinition{
		prefix:      prefix,
		handlers:    handlers,
		middlewares: middlewares,
	}
}

// Prefix returns the handlers group http path prefix.
func (h *handlersGroupDefinition) Prefix() string {
	return h.prefix
}

// Handlers returns the handlers group associated handlers.
func (h *handlersGroupDefinition) Handlers() []HandlerDefinition {
	return h.handlers
}

// Middlewares returns the handlers group associated middlewares.
func (h *handlersGroupDefinition) Middlewares() []MiddlewareDefinition {
	return h.middlewares
}
