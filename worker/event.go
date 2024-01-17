package worker

import (
	"fmt"
	"time"
)

// WorkerExecutionEvent is an event happening during a [Worker] execution.
type WorkerExecutionEvent struct {
	executionId string
	message     string
	timestamp   time.Time
}

// NewWorkerExecutionEvent returns a new [WorkerExecutionEvent].
func NewWorkerExecutionEvent(executionId string, message string, timestamp time.Time) *WorkerExecutionEvent {
	return &WorkerExecutionEvent{
		executionId: executionId,
		message:     message,
		timestamp:   timestamp,
	}
}

// ExecutionId returns the worker execution id.
func (e *WorkerExecutionEvent) ExecutionId() string {
	return e.executionId
}

// Message returns the worker execution message.
func (e *WorkerExecutionEvent) Message() string {
	return e.message
}

// Timestamp returns the worker execution timestamp.
func (e *WorkerExecutionEvent) Timestamp() time.Time {
	return e.timestamp
}

// String returns a string representation of the [WorkerExecutionEvent].
func (e *WorkerExecutionEvent) String() string {
	return fmt.Sprintf("[%s - %s] %s", e.executionId, e.timestamp.Format(time.DateTime), e.message)
}
