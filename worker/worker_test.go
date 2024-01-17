package worker_test

import (
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkerRegistration(t *testing.T) {
	t.Parallel()

	classicWorker := workers.NewClassicWorker()
	options := []worker.WorkerExecutionOption(nil)

	resolvedWorker := worker.NewWorkerRegistration(classicWorker, options...)

	assert.IsType(t, &worker.WorkerRegistration{}, resolvedWorker)
	assert.Equal(t, classicWorker, resolvedWorker.Worker())
	assert.Equal(t, options, resolvedWorker.Options())
}
