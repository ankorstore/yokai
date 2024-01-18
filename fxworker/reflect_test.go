package fxworker_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
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
		{workers.NewClassicWorker(), "*workers.ClassicWorker"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			got := fxworker.GetType(tt.target)
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

			got := fxworker.GetReturnType(tt.target)
			assert.Equal(t, tt.expected, got)
		})
	}
}
