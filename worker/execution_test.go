package worker_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkerExecution(t *testing.T) {
	t.Parallel()

	execution := worker.NewWorkerExecution("id", "name", worker.DefaultWorkerExecutionOptions())

	assert.IsType(t, &worker.WorkerExecution{}, execution)
	assert.Equal(t, "id", execution.Id())
	assert.Equal(t, "name", execution.Name())
	assert.Equal(t, worker.Unknown, execution.Status())
	assert.Equal(t, 0, execution.CurrentExecutionAttempt())
	assert.Equal(t, worker.DefaultMaxExecutionsAttempts, execution.MaxExecutionsAttempts())
	assert.Equal(t, float64(worker.DefaultDeferredStartThreshold), execution.DeferredStartThreshold())
	assert.Len(t, execution.Events(), 0)
}

func TestWorkerExecutionSetters(t *testing.T) {
	t.Parallel()

	execution := worker.NewWorkerExecution("id", "name", worker.DefaultWorkerExecutionOptions())

	execution.
		SetId("new id").
		SetName("new name").
		SetStatus(worker.Success).
		SetCurrentExecutionAttempt(2).
		SetMaxExecutionsAttempts(3).
		SetDeferredStartThreshold(4)

	assert.Equal(t, "new id", execution.Id())
	assert.Equal(t, "new name", execution.Name())
	assert.Equal(t, worker.Success, execution.Status())
	assert.Equal(t, 2, execution.CurrentExecutionAttempt())
	assert.Equal(t, 3, execution.MaxExecutionsAttempts())
	assert.Equal(t, float64(4), execution.DeferredStartThreshold())
}

func TestWorkerExecutionEventsLifecycle(t *testing.T) {
	t.Parallel()

	execution := worker.NewWorkerExecution("id", "name", worker.DefaultWorkerExecutionOptions())

	assert.Len(t, execution.Events(), 0)

	execution.
		AddEvent("event-1").
		AddEvent("event-2")

	assert.Len(t, execution.Events(), 2)

	assert.True(t, execution.HasEvent("event-1"))
	assert.True(t, execution.HasEvent("event-2"))
	assert.False(t, execution.HasEvent("event-3"))

	i := 1
	for _, event := range execution.Events() {
		assert.Equal(t, fmt.Sprintf("event-%d", i), event.Message())
		i++
	}
}
