package fxworker

import "github.com/ankorstore/yokai/worker"

// MiddlewareDefinition is the interface for middleware definitions.
type MiddlewareDefinition interface {
	ReturnType() string
}

type middlewareDefinition struct {
	returnType string
}

// NewMiddlewareDefinition returns a new [MiddlewareDefinition].
func NewMiddlewareDefinition(returnType string) MiddlewareDefinition {
	return &middlewareDefinition{
		returnType: returnType,
	}
}

// ReturnType returns the middleware definition return type.
func (m *middlewareDefinition) ReturnType() string {
	return m.returnType
}

// WorkerDefinition is the interface for workers definitions.
type WorkerDefinition interface {
	ReturnType() string
	Options() []worker.WorkerExecutionOption
	Middlewares() []MiddlewareDefinition
}

type workerDefinition struct {
	returnType  string
	options     []worker.WorkerExecutionOption
	middlewares []MiddlewareDefinition
}

// NewWorkerDefinition returns a new [WorkerDefinition].
func NewWorkerDefinition(returnType string, options ...worker.WorkerExecutionOption) WorkerDefinition {
	return &workerDefinition{
		returnType:  returnType,
		options:     options,
		middlewares: []MiddlewareDefinition{},
	}
}

// NewWorkerDefinitionWithMiddlewares returns a new [WorkerDefinition] with middlewares.
func NewWorkerDefinitionWithMiddlewares(returnType string, middlewares []MiddlewareDefinition, options ...worker.WorkerExecutionOption) WorkerDefinition {
	return &workerDefinition{
		returnType:  returnType,
		options:     options,
		middlewares: middlewares,
	}
}

// ReturnType returns the worker definition return type.
func (w *workerDefinition) ReturnType() string {
	return w.returnType
}

// Options returns the worker definition execution options.
func (w *workerDefinition) Options() []worker.WorkerExecutionOption {
	return w.options
}

// Middlewares returns the worker definition middlewares.
func (w *workerDefinition) Middlewares() []MiddlewareDefinition {
	return w.middlewares
}
