package fxcore

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type Task interface {
	Name() string
	Run(ctx context.Context, input []byte) ([]byte, error)
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

func (r *TaskRegistry) Run(ctx context.Context, name string, input []byte) ([]byte, error) {
	task, ok := r.tasks[name]
	if !ok {
		return nil, fmt.Errorf("task %s not found in registry", name)
	}

	return task.Run(ctx, input)

}
