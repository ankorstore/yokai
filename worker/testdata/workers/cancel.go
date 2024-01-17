package workers

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/worker"
)

type CancellableWorker struct{}

func NewCancellableWorker() *CancellableWorker {
	return &CancellableWorker{}
}

func (w *CancellableWorker) Name() string {
	return "CancellableWorker"
}

func (w *CancellableWorker) Run(ctx context.Context) error {
	logger := worker.CtxLogger(ctx)

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("stopping")

			time.Sleep(10 * time.Millisecond) // simulate work

			logger.Info().Msg("stopped")

			return nil
		default:
			logger.Info().Msg("running")

			time.Sleep(10 * time.Millisecond) // simulate work
		}
	}
}
