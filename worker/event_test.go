package worker_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkerExecutionEvent(t *testing.T) {
	t.Parallel()

	now := time.Now()

	event := worker.NewWorkerExecutionEvent("id", "message", now)

	assert.IsType(t, &worker.WorkerExecutionEvent{}, event)
	assert.Equal(t, "id", event.ExecutionId())
	assert.Equal(t, "message", event.Message())
	assert.Equal(t, now, event.Timestamp())
	assert.Equal(t, fmt.Sprintf("[id - %s] message", now.Format(time.DateTime)), event.String())
}
