package fxworker

import "github.com/ankorstore/yokai/worker"

// WorkerDefinition is the interface for workers definitions.
type WorkerDefinition interface {
	ReturnType() string
	Options() []worker.WorkerExecutionOption
}

type workerDefinition struct {
	returnType string
	options    []worker.WorkerExecutionOption
}

// NewWorkerDefinition returns a new [WorkerDefinition].
func NewWorkerDefinition(returnType string, options ...worker.WorkerExecutionOption) WorkerDefinition {
	return &workerDefinition{
		returnType: returnType,
		options:    options,
	}
}

// ReturnType returns the worker definition return type.
func (w *workerDefinition) ReturnType() string {
	return w.returnType
}

// Options returns the worker definition execution options.
func (w *workerDefinition) Options() []worker.WorkerExecutionOption {
	return w.options
}
