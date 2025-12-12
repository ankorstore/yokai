package worker

import "context"

// Middleware is the interface to implement to provide worker middlewares.
type Middleware interface {
	Name() string
	Handle() MiddlewareFunc
}

// MiddlewareFunc wraps handlers in the middleware chain to perform operations before and after execution.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// HandlerFunc executes a worker's Run method or another middleware in the chain with context information.
type HandlerFunc func(ctx context.Context) error
