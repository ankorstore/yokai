package fxcore

import (
	"context"
	"fmt"
	"sort"

	"go.uber.org/fx"
)

// TaskResult is a Task execution result.
type TaskResult struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// Task is an interface for tasks implementations.
type Task interface {
	Name() string
	Run(ctx context.Context, input []byte) TaskResult
}

// TaskRegistry is a registry of Task implementations.
type TaskRegistry struct {
	tasks map[string]Task
}

// TaskRegistryParams is used to inject dependencies in NewTaskRegistry.
type TaskRegistryParams struct {
	fx.In
	Tasks []Task `group:"core-tasks"`
}

// NewTaskRegistry returns a new TaskRegistry instance.
func NewTaskRegistry(p TaskRegistryParams) *TaskRegistry {
	tasks := make(map[string]Task)

	for _, task := range p.Tasks {
		tasks[task.Name()] = task
	}

	return &TaskRegistry{
		tasks: tasks,
	}
}

// Names returns all registered Task names.
func (r *TaskRegistry) Names() []string {
	var names []string

	for name := range r.tasks {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Run runs a specific Task.
func (r *TaskRegistry) Run(ctx context.Context, name string, input []byte) TaskResult {
	task, ok := r.tasks[name]
	if !ok {
		return TaskResult{
			Success: false,
			Message: fmt.Sprintf("task %s not found", name),
		}
	}

	return task.Run(ctx, input)
}
