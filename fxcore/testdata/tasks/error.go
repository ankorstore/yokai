package tasks

import (
	"context"

	"github.com/ankorstore/yokai/fxcore"
)

var _ fxcore.Task = (*ErrorTask)(nil)

type ErrorTask struct{}

func NewErrorTask() *ErrorTask {
	return &ErrorTask{}
}

func (t *ErrorTask) Name() string {
	return "error"
}

func (t *ErrorTask) Run(context.Context, []byte) fxcore.TaskResult {
	return fxcore.TaskResult{
		Success: false,
		Message: "task error",
	}
}
