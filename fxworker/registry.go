package fxworker

import (
	"fmt"

	"github.com/ankorstore/yokai/worker"
	"go.uber.org/fx"
)

// WorkerRegistry is the registry collecting workers and their definitions.
type WorkerRegistry struct {
	workers     []worker.Worker
	definitions []WorkerDefinition
	middlewares []worker.Middleware
}

// FxWorkerRegistryParam allows injection of the required dependencies in [NewFxWorkerRegistry].
type FxWorkerRegistryParam struct {
	fx.In
	Workers     []worker.Worker     `group:"workers"`
	Definitions []WorkerDefinition  `group:"workers-definitions"`
	Middlewares []worker.Middleware `group:"worker-middlewares"`
}

// NewFxWorkerRegistry returns as new [WorkerRegistry].
func NewFxWorkerRegistry(p FxWorkerRegistryParam) *WorkerRegistry {
	return &WorkerRegistry{
		workers:     p.Workers,
		definitions: p.Definitions,
		middlewares: p.Middlewares,
	}
}

// ResolveWorkersRegistrations resolves a list of [worker.WorkerRegistration] from their definitions.
func (r *WorkerRegistry) ResolveWorkersRegistrations() ([]*worker.WorkerRegistration, error) {
	registrations := []*worker.WorkerRegistration{}

	for _, definition := range r.definitions {
		implementation, err := r.lookupRegisteredWorker(definition.ReturnType())
		if err != nil {
			return nil, err
		}

		options := definition.Options()

		// Extract middlewares from definition
		if len(definition.Middlewares()) > 0 {
			resolvedMiddlewares, err := r.resolveMiddlewares(definition.Middlewares())
			if err != nil {
				return nil, err
			}

			options = append(options, worker.WithMiddlewares(resolvedMiddlewares...))
		}

		registrations = append(
			registrations,
			worker.NewWorkerRegistration(implementation, options...),
		)
	}

	return registrations, nil
}

// resolveMiddlewares resolves middleware instances from their definitions.
func (r *WorkerRegistry) resolveMiddlewares(definitions []MiddlewareDefinition) ([]worker.Middleware, error) {
	resolvedMiddlewares := []worker.Middleware{}

	for _, definition := range definitions {
		// Try to find middleware by return type
		middleware, err := r.lookupRegisteredMiddleware(definition.ReturnType())
		if err != nil {
			return nil, err
		}

		resolvedMiddlewares = append(resolvedMiddlewares, middleware)
	}

	return resolvedMiddlewares, nil
}

// lookupRegisteredMiddleware finds a middleware by its return type.
func (r *WorkerRegistry) lookupRegisteredMiddleware(returnType string) (worker.Middleware, error) {
	for _, middleware := range r.middlewares {
		if GetType(middleware) == returnType {
			return middleware, nil
		}
	}

	return nil, fmt.Errorf("cannot find middleware implementation for type %s", returnType)
}

func (r *WorkerRegistry) lookupRegisteredWorker(returnType string) (worker.Worker, error) {
	for _, implementation := range r.workers {
		if GetType(implementation) == returnType {
			return implementation, nil
		}
	}

	return nil, fmt.Errorf("cannot find worker implementation for type %s", returnType)
}
