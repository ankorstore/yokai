package healthcheck_test

import (
	"context"
	"testing"
	"time"

	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/healthcheck"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

func TestWorkerProbe(t *testing.T) {
	t.Parallel()

	t.Run("empty pool", func(t *testing.T) {
		t.Parallel()

		pool := worker.NewWorkerPool()

		probe := healthcheck.NewWorkerProbe(pool)

		res := probe.Check(context.Background())

		assert.True(t, res.Success)
		assert.Empty(t, res.Message)
	})

	t.Run("success pool", func(t *testing.T) {
		t.Parallel()

		pool, err := worker.NewDefaultWorkerPoolFactory().Create(
			worker.WithWorker(workers.NewClassicWorker()),
		)
		assert.NoError(t, err)

		probe := healthcheck.NewWorkerProbe(pool)

		err = pool.Start(context.Background())
		assert.NoError(t, err)

		time.Sleep(15 * time.Millisecond)

		res := probe.Check(context.Background())

		assert.True(t, res.Success)
		assert.Equal(t, "ClassicWorker: success", res.Message)
	})

	t.Run("error pool", func(t *testing.T) {
		t.Parallel()

		pool, err := worker.NewDefaultWorkerPoolFactory().Create(
			worker.WithWorker(workers.NewClassicWorker()),
			worker.WithWorker(workers.NewErrorWorker()),
		)
		assert.NoError(t, err)

		probe := healthcheck.NewWorkerProbe(pool)

		err = pool.Start(context.Background())
		assert.NoError(t, err)

		time.Sleep(15 * time.Millisecond)

		res := probe.Check(context.Background())

		assert.False(t, res.Success)
		assert.Equal(t, "ClassicWorker: success, ErrorWorker: error", res.Message)
	})
}
