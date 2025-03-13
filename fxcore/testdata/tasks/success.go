package tasks

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcore"
)

var _ fxcore.Task = (*SuccessTask)(nil)

type SuccessTask struct {
	config *config.Config
}

func NewSuccessTask(config *config.Config) *SuccessTask {
	return &SuccessTask{
		config: config,
	}
}

func (t *SuccessTask) Name() string {
	return "success"
}

func (t *SuccessTask) Run(ctx context.Context, input []byte) fxcore.TaskResult {
	return fxcore.TaskResult{
		Success: true,
		Message: "task success",
		Details: map[string]any{
			"app":   t.config.AppName(),
			"input": string(input),
		},
	}
}
