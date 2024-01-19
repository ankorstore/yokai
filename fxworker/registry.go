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
}

// FxWorkerRegistryParam allows injection of the required dependencies in [NewFxWorkerRegistry].
type FxWorkerRegistryParam struct {
	fx.In
	Workers     []worker.Worker    `group:"workers"`
	Definitions []WorkerDefinition `group:"workers-definitions"`
}

// NewFxWorkerRegistry returns as new [WorkerRegistry].
func NewFxWorkerRegistry(p FxWorkerRegistryParam) *WorkerRegistry {
	return &WorkerRegistry{
		workers:     p.Workers,
		definitions: p.Definitions,
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

		registrations = append(
			registrations,
			worker.NewWorkerRegistration(implementation, definition.Options()...),
		)
	}

	return registrations, nil
}

func (r *WorkerRegistry) lookupRegisteredWorker(returnType string) (worker.Worker, error) {
	for _, implementation := range r.workers {
		if GetType(implementation) == returnType {
			return implementation, nil
		}
	}

	return nil, fmt.Errorf("cannot find worker implementation for type %s", returnType)
}
