package fxworker_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkerDefinition(t *testing.T) {
	t.Parallel()

	definition := fxworker.NewWorkerDefinition("*TestWorker")

	assert.Implements(t, (*fxworker.WorkerDefinition)(nil), definition)
	assert.Equal(t, "*TestWorker", definition.ReturnType())
	assert.Equal(t, []worker.WorkerExecutionOption(nil), definition.Options())
}
