package worker

import "context"

type Worker interface {
	Name() string
	Run(ctx context.Context) error
}

type WorkerRegistration struct {
	worker  Worker
	options []WorkerExecutionOption
}

func NewWorkerRegistration(worker Worker, options ...WorkerExecutionOption) *WorkerRegistration {
	return &WorkerRegistration{
		worker:  worker,
		options: options,
	}
}

func (r *WorkerRegistration) Worker() Worker {
	return r.worker
}

func (r *WorkerRegistration) Options() []WorkerExecutionOption {
	return r.options
}
