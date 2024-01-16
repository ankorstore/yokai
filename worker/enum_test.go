package worker_test

import (
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
)

func TestWorkerStatusAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    worker.WorkerStatus
		expected string
	}{
		{worker.Unknown, "unknown"},
		{worker.Deferred, "deferred"},
		{worker.Running, "running"},
		{worker.Success, "success"},
		{worker.Error, "error"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			actual := tt.input.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
