package fxgrpcserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxhealthcheck/testdata/probes"
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
		{probes.NewSuccessProbe(), "*probes.SuccessProbe"},
		{probes.NewFailureProbe(), "*probes.FailureProbe"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			got := fxhealthcheck.GetType(tt.target)
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

			got := fxhealthcheck.GetReturnType(tt.target)
			assert.Equal(t, tt.expected, got)
		})
	}
}
