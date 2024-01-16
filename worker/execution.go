package worker

import (
	"sync"
	"time"
)

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

func (e *WorkerExecution) Id() string {
	return e.id
}

func (e *WorkerExecution) SetId(id string) *WorkerExecution {
	e.id = id

	return e
}

func (e *WorkerExecution) Name() string {
	return e.name
}

func (e *WorkerExecution) SetName(name string) *WorkerExecution {
	e.name = name

	return e
}

func (e *WorkerExecution) Status() WorkerStatus {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.status
}

func (e *WorkerExecution) SetStatus(status WorkerStatus) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.status = status

	return e
}

func (e *WorkerExecution) CurrentExecutionAttempt() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.currentExecutionAttempt
}

func (e *WorkerExecution) SetCurrentExecutionAttempt(current int) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.currentExecutionAttempt = current

	return e
}

func (e *WorkerExecution) MaxExecutionsAttempts() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.maxExecutionsAttempts
}

func (e *WorkerExecution) SetMaxExecutionsAttempts(max int) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.maxExecutionsAttempts = max

	return e
}

func (e *WorkerExecution) DeferredStartThreshold() float64 {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.deferredStartThreshold
}

func (e *WorkerExecution) SetDeferredStartThreshold(threshold float64) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.deferredStartThreshold = threshold

	return e
}

func (e *WorkerExecution) Events() []*WorkerExecutionEvent {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.events
}

func (e *WorkerExecution) AddEvent(message string) *WorkerExecution {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.events = append(e.events, NewWorkerExecutionEvent(e.id, message, time.Now()))

	return e
}

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
