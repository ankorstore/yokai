package worker_test

import (
	"testing"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

func TestWorkerPoolOptionsWithGenerator(t *testing.T) {
	t.Parallel()

	generator := uuid.NewDefaultUuidGenerator()

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithGenerator(generator)(&opt)

	assert.Equal(t, generator, opt.Generator)
}

func TestWorkerPoolOptionsWithMetrics(t *testing.T) {
	t.Parallel()

	metrics := worker.NewWorkerMetrics("foo", "bar")

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithMetrics(metrics)(&opt)

	assert.Equal(t, metrics, opt.Metrics)
}

func TestWorkerPoolOptionsWithWorker(t *testing.T) {
	t.Parallel()

	classicWorker := workers.NewClassicWorker()

	registration := worker.NewWorkerRegistration(classicWorker)

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithWorker(classicWorker)(&opt)

	assert.Equal(t, registration, opt.Registrations[classicWorker.Name()])
}

func TestWorkerPoolOptionsWithGlobalDeferredStartThresholds(t *testing.T) {
	t.Parallel()

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithGlobalDeferredStartThreshold(1.5)(&opt)

	assert.Equal(t, 1.5, opt.GlobalDeferredStartThreshold)
}

func TestWorkerPoolOptionsWithGlobalMaxExecutionsAttempts(t *testing.T) {
	t.Parallel()

	opt := worker.DefaultWorkerPoolOptions()
	worker.WithGlobalMaxExecutionsAttempts(2)(&opt)

	assert.Equal(t, 2, opt.GlobalMaxExecutionsAttempts)
}

func TestWorkerExecutionOptionsWithDeferredStartThreshold(t *testing.T) {
	t.Parallel()

	opt := worker.DefaultWorkerExecutionOptions()
	worker.WithDeferredStartThreshold(1.5)(&opt)

	assert.Equal(t, 1.5, opt.DeferredStartThreshold)
}

func TestWorkerExecutionOptionsWithMaxExecutionsAttempts(t *testing.T) {
	t.Parallel()

	opt := worker.DefaultWorkerExecutionOptions()
	worker.WithMaxExecutionsAttempts(2)(&opt)

	assert.Equal(t, 2, opt.MaxExecutionsAttempts)
}
