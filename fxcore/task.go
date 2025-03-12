package fxcore

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type TaskResult struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type Task interface {
	Name() string
	Run(ctx context.Context, input []byte) TaskResult
}

type TaskRegistry struct {
	tasks map[string]Task
}

type TaskRegistryParams struct {
	fx.In
	Tasks []Task `group:"core-tasks"`
}

func NewTaskRegistry(p TaskRegistryParams) *TaskRegistry {
	tasks := make(map[string]Task)

	for _, task := range p.Tasks {
		tasks[task.Name()] = task
	}

	return &TaskRegistry{
		tasks: tasks,
	}
}

func (r *TaskRegistry) Names() []string {
	var names []string

	for name := range r.tasks {
		names = append(names, name)
	}
	return names
}

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
