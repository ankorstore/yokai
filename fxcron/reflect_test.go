package fxcron_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/tracker"
	"github.com/stretchr/testify/assert"
)

func TestGetType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target   any
		expected string
	}{
		{123, "int"},
		{"test", "string"},
		{tracker.NewCronExecutionTracker(), "*tracker.CronExecutionTracker"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			got := fxcron.GetType(tt.target)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetReturnType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target   any
		expected string
	}{
		{func() string { return "test" }, "string"},
		{func() int { return 123 }, "int"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			got := fxcron.GetReturnType(tt.target)
			assert.Equal(t, tt.expected, got)
		})
	}
}
