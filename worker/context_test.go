package worker_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
)

func TestCtxEmptyCtxWorkerName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", worker.CtxWorkerName(context.Background()))
}

func TestCtxEmptyCtxWorkerExecutionId(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", worker.CtxWorkerExecutionId(context.Background()))
}
