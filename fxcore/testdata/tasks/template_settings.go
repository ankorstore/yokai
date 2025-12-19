package tasks

import (
	"context"

	"github.com/ankorstore/yokai/fxcore"
)

var _ fxcore.Task = (*TemplateSettingsTask)(nil)
var _ fxcore.TaskWithTemplateSettings = (*TemplateSettingsTask)(nil)

type TemplateSettingsTask struct{}

func NewTemplateSettingsTask() *TemplateSettingsTask {
	return &TemplateSettingsTask{}
}

func (t *TemplateSettingsTask) Name() string {
	return "template-settings"
}

func (t *TemplateSettingsTask) Run(context.Context, []byte) fxcore.TaskResult {
	return fxcore.TaskResult{
		Success: true,
		Message: "template settings task",
	}
}

func (t *TemplateSettingsTask) TemplateSettings(settings fxcore.TaskTemplateSettings) fxcore.TaskTemplateSettings {
	settings.Placeholder = "Custom placeholder"
	settings.DefaultValue = "Default content"
	settings.Rows = 5

	return settings
}
