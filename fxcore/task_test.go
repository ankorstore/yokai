package fxcore_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxcore/testdata/tasks"
	"github.com/stretchr/testify/assert"
)

func TestTaskRegistry(t *testing.T) {
	t.Parallel()

	createRegistry := func(tb testing.TB) *fxcore.TaskRegistry {
		tb.Helper()

		cfg, err := config.NewDefaultConfigFactory().Create(
			config.WithFilePaths("./testdata/config"),
		)
		assert.NoError(tb, err)

		return fxcore.NewTaskRegistry(fxcore.TaskRegistryParams{
			Tasks: []fxcore.Task{
				tasks.NewSuccessTask(cfg),
				tasks.NewErrorTask(),
				tasks.NewTemplateSettingsTask(),
			},
		})
	}

	t.Run("test names", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)

		assert.Equal(t, []string{"error", "success", "template-settings"}, registry.Names())
	})

	t.Run("test run with success task", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)

		res := registry.Run(context.Background(), "success", []byte("test input"))

		assert.True(t, res.Success)
		assert.Equal(t, "task success", res.Message)
		assert.Equal(
			t,
			map[string]any{
				"app":   "core-app",
				"input": "test input",
			},
			res.Details,
		)
	})

	t.Run("test run with error task", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)

		res := registry.Run(context.Background(), "error", []byte("test input"))

		assert.False(t, res.Success)
		assert.Equal(t, "task error", res.Message)
		assert.Nil(t, res.Details)
	})

	t.Run("test run with invalid task", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)

		res := registry.Run(context.Background(), "invalid", []byte("test input"))

		assert.False(t, res.Success)
		assert.Equal(t, "task invalid not found", res.Message)
		assert.Nil(t, res.Details)
	})

	t.Run("test template settings", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)
		settings := registry.TemplateSettings()

		assert.Len(t, settings, 3)

		successSettings, ok := settings["success"]
		assert.True(t, ok)
		assert.Equal(t, "Optional input...", successSettings.Placeholder)
		assert.Equal(t, "", successSettings.DefaultValue)
		assert.Equal(t, 1, successSettings.Rows)
		assert.True(t, successSettings.EscapeContent)

		errorSettings, ok := settings["error"]
		assert.True(t, ok)
		assert.Equal(t, "Optional input...", errorSettings.Placeholder)
		assert.Equal(t, "", errorSettings.DefaultValue)
		assert.Equal(t, 1, errorSettings.Rows)
		assert.True(t, errorSettings.EscapeContent)

		templateSettings, ok := settings["template-settings"]
		assert.True(t, ok)
		assert.Equal(t, "Custom placeholder", templateSettings.Placeholder)
		assert.Equal(t, "Default content", templateSettings.DefaultValue)
		assert.Equal(t, 5, templateSettings.Rows)
		assert.False(t, templateSettings.EscapeContent)
	})
}
