package worker

import (
	"sync"
	"time"
)

// WorkerExecution represents a [Worker] execution within the [WorkerPool].
type WorkerExecution struct {
	mutex                   sync.Mutex
	id                      string
	name                    string
	status                  WorkerStatus
	currentExecutionAttempt int
	maxExecutionsAttempts   int
	deferredStartThreshold  float64
	events                  []*WorkerExecutionEvent
}

// NewWorkerExecution returns a new [WorkerExecution].
func NewWorkerExecution(id string, name string, options ExecutionOptions) *WorkerExecution {
	return &WorkerExecution{
		id:                      id,
		name:                    name,
		status:                  Unknown,
		currentExecutionAttempt: 0,
		maxExecutionsAttempts:   options.MaxExecutionsAttempts,
		deferredStartThreshold:  options.DeferredStartThreshold,
		events:                  []*WorkerExecutionEvent{},
	}
}

// Id returns the [WorkerExecution] id.
func (e *WorkerExecution) Id() string {
	return e.id
}

// SetId sets the [WorkerExecution] id.
func (e *WorkerExecution) SetId(id string) *WorkerExecution {
	e.id = id

	return e
}

// Name returns the [WorkerExecution] name.
func (e *WorkerExecution) Name() string {
	return e.name
}

// SetName sets the [WorkerExecution] name.
func (e *WorkerExecution) SetName(name string) *WorkerExecution {
	e.name = name

	return e
}

// Status returns the [WorkerExecution] status.
func (e *WorkerExecution) Status() WorkerStatus {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.status
}

// SetStatus sets the [WorkerExecution] status.
func (e *WorkerExecution) SetStatus(status WorkerStatus) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.status = status

	return e
}

// CurrentExecutionAttempt returns the [WorkerExecution] current execution attempt.
func (e *WorkerExecution) CurrentExecutionAttempt() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.currentExecutionAttempt
}

// SetCurrentExecutionAttempt sets the [WorkerExecution] current execution attempt.
func (e *WorkerExecution) SetCurrentExecutionAttempt(current int) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.currentExecutionAttempt = current

	return e
}

// MaxExecutionsAttempts returns the [WorkerExecution] max execution attempts.
func (e *WorkerExecution) MaxExecutionsAttempts() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.maxExecutionsAttempts
}

// SetMaxExecutionsAttempts sets the [WorkerExecution] max execution attempts.
func (e *WorkerExecution) SetMaxExecutionsAttempts(max int) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.maxExecutionsAttempts = max

	return e
}

// DeferredStartThreshold returns the [WorkerExecution] max deferred start threshold, in seconds.
func (e *WorkerExecution) DeferredStartThreshold() float64 {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.deferredStartThreshold
}

// SetDeferredStartThreshold sets the [WorkerExecution] max deferred start threshold, in seconds.
func (e *WorkerExecution) SetDeferredStartThreshold(threshold float64) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.deferredStartThreshold = threshold

	return e
}

// Events returns the [WorkerExecution] list of [WorkerExecutionEvent].
func (e *WorkerExecution) Events() []*WorkerExecutionEvent {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.events
}

// AddEvent adds a [WorkerExecutionEvent] to the [WorkerExecution].
func (e *WorkerExecution) AddEvent(message string) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.events = append(e.events, NewWorkerExecutionEvent(e.id, message, time.Now()))

	return e
}

// HasEvent returns true if a [WorkerExecutionEvent] was found for a given message.
func (e *WorkerExecution) HasEvent(message string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, event := range e.events {
		if event.Message() == message {
			return true
		}
	}

	return false
}
