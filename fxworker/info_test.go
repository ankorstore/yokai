package fxworker_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

const testExecutionId = "test-execution-id"

func TestNewFxWorkerModuleInfo(t *testing.T) {
	t.Parallel()

	generator := uuid.NewTestUuidGenerator(testExecutionId)

	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithGenerator(generator),
		worker.WithWorker(workers.NewClassicWorker(), worker.WithDeferredStartThreshold(1)),
	)
	assert.NoError(t, err)

	err = pool.Start(context.Background())
	assert.NoError(t, err)

	info := fxworker.NewFxWorkerModuleInfo(pool)
	assert.Equal(t, fxworker.ModuleName, info.Name())

	assert.Equal(
		t,
		map[string]interface{}{
			"workers": map[string]interface{}{
				"ClassicWorker": map[string]interface{}{
					"status": worker.Unknown.String(),
					"events": []map[string]string(nil),
				},
			},
		},
		info.Data(),
	)
}
