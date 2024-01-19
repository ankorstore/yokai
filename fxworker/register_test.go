package fxworker_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

func TestAsWorker(t *testing.T) {
	t.Parallel()

	result := fxworker.AsWorker(workers.NewClassicWorker, worker.WithMaxExecutionsAttempts(2))

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
