package fxworker_test

import (
	"context"
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

const testExecutionId = "test-execution-id"

func TestNewFxWorkerModuleInfo(t *testing.T) {
	t.Parallel()

	// wait until the beginning of the next second to avoid flaky test
	now := time.Now()
	nextSecond := now.Truncate(time.Second).Add(time.Second)
	sleepTime := nextSecond.Sub(now)
	time.Sleep(sleepTime)

	generator := uuid.NewTestUuidGenerator(testExecutionId)

	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithGenerator(generator),
		worker.WithWorker(workers.NewClassicWorker()),
	)
	assert.NoError(t, err)

	info := fxworker.NewFxWorkerModuleInfo(pool)
	assert.Equal(t, fxworker.ModuleName, info.Name())

	err = pool.Start(context.Background())
	assert.NoError(t, err)

	time.Sleep(5 * time.Millisecond)

	err = pool.Stop()
	assert.NoError(t, err)

	assert.Equal(
		t,
		map[string]interface{}{
			"workers": map[string]interface{}{
				"ClassicWorker": map[string]interface{}{
					"status": worker.Success.String(),
					"events": []map[string]string{
						{
							"execution": "test-execution-id",
							"message":   "starting execution attempt 1/1",
							"time":      nextSecond.Format(time.DateTime),
						},
						{
							"execution": "test-execution-id",
							"message":   "stopping execution attempt 1/1 with success",
							"time":      nextSecond.Format(time.DateTime),
						},
					},
				},
			},
		},
		info.Data(),
	)
}
