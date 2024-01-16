package worker_test

import (
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
)

func TestDefaultWorkerPoolFactory(t *testing.T) {
	t.Parallel()

	factory := worker.NewDefaultWorkerPoolFactory()

	assert.IsType(t, &worker.DefaultWorkerPoolFactory{}, factory)
	assert.Implements(t, (*worker.WorkerPoolFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	oneShotWorker := workers.NewOneShotWorker()

	generator := uuid.NewDefaultUuidGenerator()

	metrics := worker.NewWorkerMetrics("foo", "bar")

	factory := worker.NewDefaultWorkerPoolFactory()

	pool, err := factory.Create(
		worker.WithWorker(oneShotWorker),
		worker.WithGenerator(generator),
		worker.WithMetrics(metrics),
		worker.WithGlobalMaxExecutionsAttempts(1),
		worker.WithGlobalDeferredStartThreshold(2),
	)
	assert.NoError(t, err)

	assert.IsType(t, &worker.WorkerPool{}, pool)

	options := pool.Options()
	assert.Equal(t, generator, options.Generator)
	assert.Equal(t, metrics, options.Metrics)
	assert.Equal(t, 1, options.GlobalMaxExecutionsAttempts)
	assert.Equal(t, float64(2), options.GlobalDeferredStartThreshold)
	assert.Len(t, options.Registrations, 1)
	assert.Equal(t, oneShotWorker, options.Registrations[oneShotWorker.Name()].Worker())
}
