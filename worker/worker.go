package worker

import "context"

// Worker is the interface to implement to provide workers.
type Worker interface {
	Name() string
	Run(ctx context.Context) error
}

// WorkerRegistration is a [Worker] registration, with optional [WorkerExecutionOption].
type WorkerRegistration struct {
	worker  Worker
	options []WorkerExecutionOption
}

// NewWorkerRegistration returns a new [WorkerRegistration] for a given [Worker] and an optional list of [WorkerRegistration].
func NewWorkerRegistration(worker Worker, options ...WorkerExecutionOption) *WorkerRegistration {
	return &WorkerRegistration{
		worker:  worker,
		options: options,
	}
}

// Worker returns the [Worker] of the [WorkerRegistration].
func (r *WorkerRegistration) Worker() Worker {
	return r.worker
}

// Options returns the list of [WorkerExecutionOption] of the [WorkerRegistration].
func (r *WorkerRegistration) Options() []WorkerExecutionOption {
	return r.options
}
