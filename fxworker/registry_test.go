package fxworker_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

func TestNewCheckerProbeRegistry(t *testing.T) {
	t.Parallel()

	param := fxworker.FxWorkerRegistryParam{
		Workers:     []worker.Worker{},
		Definitions: []fxworker.WorkerDefinition{},
	}

	registry := fxworker.NewFxWorkerRegistry(param)

	assert.IsType(t, &fxworker.WorkerRegistry{}, registry)
}

func TestResolveCheckerProbesRegistrationsSuccess(t *testing.T) {
	t.Parallel()

	param := fxworker.FxWorkerRegistryParam{
		Workers: []worker.Worker{
			workers.NewClassicWorker(),
			workers.NewCancellableWorker(),
		},
		Definitions: []fxworker.WorkerDefinition{
			fxworker.NewWorkerDefinition("github.com/ankorstore/yokai/worker/testdata/workers.ClassicWorker"),
			fxworker.NewWorkerDefinition("github.com/ankorstore/yokai/worker/testdata/workers.CancellableWorker"),
		},
	}

	registry := fxworker.NewFxWorkerRegistry(param)

	registrations, err := registry.ResolveWorkersRegistrations()
	assert.NoError(t, err)

	assert.Len(t, registrations, 2)
	assert.IsType(t, &workers.ClassicWorker{}, registrations[0].Worker())
	assert.IsType(t, &workers.CancellableWorker{}, registrations[1].Worker())
}

func TestResolveCheckerProbesRegistrationsFailure(t *testing.T) {
	t.Parallel()

	param := fxworker.FxWorkerRegistryParam{
		Workers: []worker.Worker{
			workers.NewClassicWorker(),
		},
		Definitions: []fxworker.WorkerDefinition{
			fxworker.NewWorkerDefinition("invalid"),
		},
	}

	registry := fxworker.NewFxWorkerRegistry(param)

	_, err := registry.ResolveWorkersRegistrations()
	assert.Error(t, err)
	assert.Equal(t, "cannot find worker implementation for type invalid", err.Error())
}

func TestResolveWorkersRegistrationsWithMiddlewares(t *testing.T) {
	t.Parallel()

	// Create middleware definition
	middlewareDef := fxworker.NewMiddlewareDefinition("github.com/ankorstore/yokai/fxworker_test.TestMiddleware")

	// Create worker definition with middleware
	workerDef := fxworker.NewWorkerDefinitionWithMiddlewares(
		"github.com/ankorstore/yokai/worker/testdata/workers.ClassicWorker",
		[]fxworker.MiddlewareDefinition{middlewareDef},
	)

	// Create middleware instance
	middleware := &TestMiddleware{}

	// Create registry param
	param := fxworker.FxWorkerRegistryParam{
		Workers: []worker.Worker{
			workers.NewClassicWorker(),
		},
		Definitions: []fxworker.WorkerDefinition{
			workerDef,
		},
		Middlewares: []worker.Middleware{
			middleware,
		},
	}

	// Create registry
	registry := fxworker.NewFxWorkerRegistry(param)

	// Resolve worker registrations
	registrations, err := registry.ResolveWorkersRegistrations()
	assert.NoError(t, err)
	assert.Len(t, registrations, 1)

	// Get the worker registration
	registration := registrations[0]
	assert.IsType(t, &workers.ClassicWorker{}, registration.Worker())

	// Check that the registration has options (which should include the middleware)
	assert.Len(t, registration.Options(), 1)
}
