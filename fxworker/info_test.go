package fxworker_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

func TestNewFxWorkerModuleInfo(t *testing.T) {
	t.Parallel()

	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(workers.NewClassicWorker(), worker.WithDeferredStartThreshold(0.2)),
		worker.WithWorker(workers.NewCancellableWorker(), worker.WithMaxExecutionsAttempts(2)),
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
				"CancellableWorker": map[string]interface{}{
					"status": worker.Unknown.String(),
					"events": []map[string]string(nil),
				},
			},
		},
		info.Data(),
	)
}
