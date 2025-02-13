package fxhttpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Middleware is the interface for middlewares.
type Middleware interface {
	Handle() echo.MiddlewareFunc
}

// Handler is the interface for handlers.
type Handler interface {
	Handle() echo.HandlerFunc
}

// ErrorHandler is the interface for error handlers.
type ErrorHandler interface {
	Handle() echo.HTTPErrorHandler
}

// HttpServerRegistry is the registry collecting middlewares, handlers, handlers groups and their definitions.
type HttpServerRegistry struct {
	middlewares              []Middleware
	middlewareDefinitions    []MiddlewareDefinition
	handlers                 []Handler
	handlerDefinitions       []HandlerDefinition
	handlersGroupDefinitions []HandlersGroupDefinition
	errorHandlers            []ErrorHandler
}

// FxHttpServerRegistryParam allows injection of the required dependencies in [NewFxHttpServerRegistry].
type FxHttpServerRegistryParam struct {
	fx.In
	Middlewares              []Middleware              `group:"httpserver-middlewares"`
	MiddlewareDefinitions    []MiddlewareDefinition    `group:"httpserver-middleware-definitions"`
	Handlers                 []Handler                 `group:"httpserver-handlers"`
	HandlerDefinitions       []HandlerDefinition       `group:"httpserver-handler-definitions"`
	HandlersGroupDefinitions []HandlersGroupDefinition `group:"httpserver-handlers-group-definitions"`
	ErrorHandlers            []ErrorHandler            `group:"httpserver-error-handlers"`
}

// NewFxHttpServerRegistry returns as new [HttpServerRegistry].
func NewFxHttpServerRegistry(p FxHttpServerRegistryParam) *HttpServerRegistry {
	return &HttpServerRegistry{

		middlewares:              p.Middlewares,
		middlewareDefinitions:    p.MiddlewareDefinitions,
		handlers:                 p.Handlers,
		handlerDefinitions:       p.HandlerDefinitions,
		handlersGroupDefinitions: p.HandlersGroupDefinitions,
		errorHandlers:            p.ErrorHandlers,
	}
}

// ResolveMiddlewares resolves a list of [ResolvedMiddleware] from their definitions.
func (r *HttpServerRegistry) ResolveMiddlewares() ([]ResolvedMiddleware, error) {
	var resolvedMiddlewares []ResolvedMiddleware

	for _, middlewareDef := range r.middlewareDefinitions {
		if middlewareDef.Kind() != Attached {
			resMiddleware, err := r.resolveMiddlewareDefinition(middlewareDef)
			if err != nil {
				return nil, err
			}

			resolvedMiddlewares = append(resolvedMiddlewares, resMiddleware)
		}
	}

	return resolvedMiddlewares, nil
}

// ResolveHandlers resolves a list of [ResolvedHandler] from their definitions.
func (r *HttpServerRegistry) ResolveHandlers() ([]ResolvedHandler, error) {
	var resolvedHandlers []ResolvedHandler

	for _, handlerDef := range r.handlerDefinitions {
		var handlerMiddlewares []echo.MiddlewareFunc

		for _, middlewareDef := range handlerDef.Middlewares() {
			handlerMiddleware, err := r.resolveMiddlewareDefinition(middlewareDef)
			if err != nil {
				return nil, err
			}

			handlerMiddlewares = append(handlerMiddlewares, handlerMiddleware.Middleware())
		}

		resHandler, err := r.resolveHandlerDefinition(handlerDef, handlerMiddlewares)
		if err != nil {
			return nil, err
		}

		resolvedHandlers = append(resolvedHandlers, resHandler)
	}

	return resolvedHandlers, nil
}

// ResolveHandlersGroups resolves a list of [ResolvedHandlersGroup] from their definitions.
func (r *HttpServerRegistry) ResolveHandlersGroups() ([]ResolvedHandlersGroup, error) {
	var resolvedHandlersGroups []ResolvedHandlersGroup

	for _, handlerGroupDef := range r.handlersGroupDefinitions {
		var groupMiddlewares []echo.MiddlewareFunc

		for _, middlewareDef := range handlerGroupDef.Middlewares() {
			groupMiddleware, err := r.resolveMiddlewareDefinition(middlewareDef)
			if err != nil {
				return nil, err
			}

			groupMiddlewares = append(groupMiddlewares, groupMiddleware.Middleware())
		}

		var groupHandlers []ResolvedHandler

		for _, handlerDef := range handlerGroupDef.Handlers() {
			var resolvedHandlerMiddlewares []echo.MiddlewareFunc

			for _, middlewareDef := range handlerDef.Middlewares() {
				resolvedHandlerMiddleware, err := r.resolveMiddlewareDefinition(middlewareDef)
				if err != nil {
					return nil, err
				}

				resolvedHandlerMiddlewares = append(resolvedHandlerMiddlewares, resolvedHandlerMiddleware.Middleware())
			}

			groupHandler, err := r.resolveHandlerDefinition(handlerDef, resolvedHandlerMiddlewares)
			if err != nil {
				return nil, err
			}

			groupHandlers = append(groupHandlers, groupHandler)
		}

		resolvedHandlersGroups = append(
			resolvedHandlersGroups,
			NewResolvedHandlersGroup(
				handlerGroupDef.Prefix(),
				groupHandlers,
				groupMiddlewares...,
			),
		)
	}

	return resolvedHandlersGroups, nil
}

// ResolveErrorHandlers resolves resolves a list of [ErrorHandler].
func (r *HttpServerRegistry) ResolveErrorHandlers() []ErrorHandler {
	return r.errorHandlers
}

func (r *HttpServerRegistry) resolveMiddlewareDefinition(middlewareDefinition MiddlewareDefinition) (ResolvedMiddleware, error) {
	if middlewareDefinition.Concrete() {
		if castMiddleware, ok := middlewareDefinition.Middleware().(func(echo.HandlerFunc) echo.HandlerFunc); ok {
			return NewResolvedMiddleware(castMiddleware, middlewareDefinition.Kind()), nil
		} else if castMiddleware, ok = middlewareDefinition.Middleware().(echo.MiddlewareFunc); ok {
			return NewResolvedMiddleware(castMiddleware, middlewareDefinition.Kind()), nil
		} else {
			return nil, fmt.Errorf("cannot cast middleware definition as MiddlewareFunc")
		}
	}

	registeredMiddleware, err := r.lookupRegisteredMiddleware(middlewareDefinition.Middleware().(string))
	if err != nil {
		return nil, fmt.Errorf("cannot lookup registered middleware")
	}

	return NewResolvedMiddleware(
		registeredMiddleware.Handle(),
		middlewareDefinition.Kind(),
	), nil
}

func (r *HttpServerRegistry) resolveHandlerDefinition(handlerDefinition HandlerDefinition, handlerMiddlewares []echo.MiddlewareFunc) (ResolvedHandler, error) {
	if handlerDefinition.Concrete() {
		if castHandler, ok := handlerDefinition.Handler().(func(echo.Context) error); ok {
			return NewResolvedHandler(
				handlerDefinition.Method(),
				handlerDefinition.Path(),
				castHandler,
				handlerMiddlewares...,
			), nil
		} else if castHandler, ok = handlerDefinition.Handler().(echo.HandlerFunc); ok {
			return NewResolvedHandler(
				handlerDefinition.Method(),
				handlerDefinition.Path(),
				castHandler,
				handlerMiddlewares...,
			), nil
		} else {
			return nil, fmt.Errorf("cannot cast handler definition as HandlerFunc")
		}
	}

	registeredHandler, err := r.lookupRegisteredHandler(handlerDefinition.Handler().(string))
	if err != nil {
		return nil, fmt.Errorf("cannot lookup registered handler")
	}

	return NewResolvedHandler(
		handlerDefinition.Method(),
		handlerDefinition.Path(),
		registeredHandler.Handle(),
		handlerMiddlewares...,
	), nil
}

func (r *HttpServerRegistry) lookupRegisteredMiddleware(middleware string) (Middleware, error) {
	for _, m := range r.middlewares {
		if GetType(m) == middleware {
			return m, nil
		}
	}

	return nil, fmt.Errorf("cannot find middleware for type %s", middleware)
}

func (r *HttpServerRegistry) lookupRegisteredHandler(handler string) (Handler, error) {
	for _, h := range r.handlers {
		if GetType(h) == handler {
			return h, nil
		}
	}

	return nil, fmt.Errorf("cannot find handler for type %s", handler)
}
