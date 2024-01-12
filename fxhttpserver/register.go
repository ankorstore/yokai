package fxhttpserver

import (
	"go.uber.org/fx"
)

// MiddlewareRegistration is a middleware registration.
type MiddlewareRegistration struct {
	middleware any
	kind       MiddlewareKind
}

// NewMiddlewareRegistration returns a new [MiddlewareRegistration].
func NewMiddlewareRegistration(middleware any, kind MiddlewareKind) *MiddlewareRegistration {
	return &MiddlewareRegistration{
		middleware: middleware,
		kind:       kind,
	}
}

// Middleware returns the middleware.
func (m *MiddlewareRegistration) Middleware() any {
	return m.middleware
}

// Kind returns the middleware kind.
func (m *MiddlewareRegistration) Kind() MiddlewareKind {
	return m.kind
}

// AsMiddleware registers a middleware into Fx.
func AsMiddleware(middleware any, kind MiddlewareKind) fx.Option {
	return RegisterMiddleware(NewMiddlewareRegistration(middleware, kind))
}

// RegisterMiddleware registers a middleware registration into Fx.
func RegisterMiddleware(middlewareRegistration *MiddlewareRegistration) fx.Option {
	var providers []any

	var middlewareDef MiddlewareDefinition
	if !IsConcreteMiddleware(middlewareRegistration.Middleware()) {
		providers = append(
			providers,
			fx.Annotate(
				middlewareRegistration.Middleware(),
				fx.As(new(Middleware)),
				fx.ResultTags(`group:"httpserver-middlewares"`),
			),
		)

		middlewareDef = NewMiddlewareDefinition(GetReturnType(middlewareRegistration.Middleware()), middlewareRegistration.kind)
	} else {
		middlewareDef = NewMiddlewareDefinition(middlewareRegistration.Middleware(), middlewareRegistration.kind)
	}

	return fx.Options(
		fx.Provide(providers...),
		fx.Supply(
			fx.Annotate(
				middlewareDef,
				fx.As(new(MiddlewareDefinition)),
				fx.ResultTags(`group:"httpserver-middleware-definitions"`),
			),
		),
	)
}

// HandlerRegistration is a handler registration.
type HandlerRegistration struct {
	method      string
	path        string
	handler     any
	middlewares []any
}

// NewHandlerRegistration returns a new [HandlerRegistration].
func NewHandlerRegistration(method string, path string, handler any, middlewares ...any) *HandlerRegistration {
	return &HandlerRegistration{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

// Method returns the handler http method.
func (h *HandlerRegistration) Method() string {
	return h.method
}

// Path returns the handler http path.
func (h *HandlerRegistration) Path() string {
	return h.path
}

// Handler returns the handler.
func (h *HandlerRegistration) Handler() any {
	return h.handler
}

// Middlewares returns the handler associated middlewares.
func (h *HandlerRegistration) Middlewares() []any {
	return h.middlewares
}

// AsHandler registers a handler into Fx.
func AsHandler(method string, path string, handler any, middlewares ...any) fx.Option {
	return RegisterHandler(NewHandlerRegistration(method, path, handler, middlewares...))
}

// RegisterHandler registers a handler registration into Fx.
func RegisterHandler(handlerRegistration *HandlerRegistration) fx.Option {
	var providers []any

	var middlewareDefs []MiddlewareDefinition
	for _, middleware := range handlerRegistration.Middlewares() {
		if !IsConcreteMiddleware(middleware) {
			providers = append(
				providers,
				fx.Annotate(
					middleware,
					fx.As(new(Middleware)),
					fx.ResultTags(`group:"httpserver-middlewares"`),
				),
			)

			middlewareDefs = append(middlewareDefs, NewMiddlewareDefinition(GetReturnType(middleware), Attached))
		} else {
			middlewareDefs = append(middlewareDefs, NewMiddlewareDefinition(middleware, Attached))
		}
	}

	var handlerDef HandlerDefinition
	if !IsConcreteHandler(handlerRegistration.Handler()) {
		providers = append(
			providers,
			fx.Annotate(
				handlerRegistration.Handler(),
				fx.As(new(Handler)),
				fx.ResultTags(`group:"httpserver-handlers"`),
			),
		)
		handlerDef = NewHandlerDefinition(
			handlerRegistration.Method(),
			handlerRegistration.Path(),
			GetReturnType(handlerRegistration.Handler()),
			middlewareDefs,
		)
	} else {
		handlerDef = NewHandlerDefinition(
			handlerRegistration.Method(),
			handlerRegistration.Path(),
			handlerRegistration.Handler(),
			middlewareDefs,
		)
	}

	return fx.Options(
		fx.Provide(providers...),
		fx.Supply(
			fx.Annotate(
				handlerDef,
				fx.As(new(HandlerDefinition)),
				fx.ResultTags(`group:"httpserver-handler-definitions"`),
			),
		),
	)
}

// HandlersGroupRegistration is a handlers group registration.
type HandlersGroupRegistration struct {
	prefix                string
	handlersRegistrations []*HandlerRegistration
	middlewares           []any
}

// NewHandlersGroupRegistration returns a new [HandlersGroupRegistration].
func NewHandlersGroupRegistration(prefix string, handlersRegistrations []*HandlerRegistration, middlewares ...any) *HandlersGroupRegistration {
	return &HandlersGroupRegistration{
		prefix:                prefix,
		handlersRegistrations: handlersRegistrations,
		middlewares:           middlewares,
	}
}

// Prefix returns the handlers group http path prefix.
func (h *HandlersGroupRegistration) Prefix() string {
	return h.prefix
}

// HandlersRegistrations returns the handlers group associated handlers registrations.
func (h *HandlersGroupRegistration) HandlersRegistrations() []*HandlerRegistration {
	return h.handlersRegistrations
}

// Middlewares returns the handlers group associated middlewares.
func (h *HandlersGroupRegistration) Middlewares() []any {
	return h.middlewares
}

// AsHandlersGroup registers a handlers group into Fx.
func AsHandlersGroup(prefix string, handlersRegistrations []*HandlerRegistration, middlewares ...any) fx.Option {
	return RegisterHandlersGroup(NewHandlersGroupRegistration(prefix, handlersRegistrations, middlewares...))
}

// RegisterHandlersGroup registers a handlers group registration into Fx.
func RegisterHandlersGroup(handlersGroupRegistration *HandlersGroupRegistration) fx.Option {
	var providers []any

	var groupMiddlewareDefs []MiddlewareDefinition
	for _, middleware := range handlersGroupRegistration.Middlewares() {
		if !IsConcreteMiddleware(middleware) {
			providers = append(
				providers,
				fx.Annotate(
					middleware,
					fx.As(new(Middleware)),
					fx.ResultTags(`group:"httpserver-middlewares"`),
				),
			)

			groupMiddlewareDefs = append(groupMiddlewareDefs, NewMiddlewareDefinition(GetReturnType(middleware), Attached))
		} else {
			groupMiddlewareDefs = append(groupMiddlewareDefs, NewMiddlewareDefinition(middleware, Attached))
		}
	}

	var groupHandlerDefs []HandlerDefinition
	for _, handlerRegistration := range handlersGroupRegistration.HandlersRegistrations() {
		var handlerDef HandlerDefinition
		var middlewareDefs []MiddlewareDefinition

		for _, middleware := range handlerRegistration.Middlewares() {
			if !IsConcreteMiddleware(middleware) {
				providers = append(
					providers,
					fx.Annotate(
						middleware,
						fx.As(new(Middleware)),
						fx.ResultTags(`group:"httpserver-middlewares"`),
					),
				)

				middlewareDefs = append(middlewareDefs, NewMiddlewareDefinition(GetReturnType(middleware), Attached))
			} else {
				middlewareDefs = append(middlewareDefs, NewMiddlewareDefinition(middleware, Attached))
			}
		}

		if !IsConcreteHandler(handlerRegistration.Handler()) {
			providers = append(
				providers,
				fx.Annotate(
					handlerRegistration.Handler(),
					fx.As(new(Handler)),
					fx.ResultTags(`group:"httpserver-handlers"`),
				),
			)
			handlerDef = NewHandlerDefinition(
				handlerRegistration.Method(),
				handlerRegistration.Path(),
				GetReturnType(handlerRegistration.Handler()),
				middlewareDefs,
			)
		} else {
			handlerDef = NewHandlerDefinition(
				handlerRegistration.Method(),
				handlerRegistration.Path(),
				handlerRegistration.Handler(),
				middlewareDefs,
			)
		}

		groupHandlerDefs = append(groupHandlerDefs, handlerDef)
	}

	handlersGroupDef := NewHandlersGroupDefinition(
		handlersGroupRegistration.Prefix(),
		groupHandlerDefs,
		groupMiddlewareDefs,
	)

	return fx.Options(
		fx.Provide(providers...),
		fx.Supply(
			fx.Annotate(
				handlersGroupDef,
				fx.As(new(HandlersGroupDefinition)),
				fx.ResultTags(`group:"httpserver-handlers-group-definitions"`),
			),
		),
	)
}
