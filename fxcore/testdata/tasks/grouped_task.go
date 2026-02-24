package tasks

import (
	"context"

	"github.com/ankorstore/yokai/fxcore"
)

var _ fxcore.Task = (*GroupedTask)(nil)
var _ fxcore.GroupedTask = (*GroupedTask)(nil)

type GroupedTask struct {
	name  string
	group string
}

func NewGroupedTask(name, group string) *GroupedTask {
	return &GroupedTask{name: name, group: group}
}

func (t *GroupedTask) Name() string {
	return t.name
}

func (t *GroupedTask) Group() string {
	return t.group
}

func (t *GroupedTask) Run(context.Context, []byte) fxcore.TaskResult {
	return fxcore.TaskResult{
		Success: true,
		Message: "grouped task " + t.name,
	}
}
