package workers

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/worker"
)

type PanicWorker struct{}

func NewPanicWorker() *PanicWorker {
	return &PanicWorker{}
}

func (w *PanicWorker) Name() string {
	return "PanicWorker"
}

func (w *PanicWorker) Run(ctx context.Context) error {
	logger := worker.CtxLogger(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			logger.Info().Msg("running")

			time.Sleep(10 * time.Millisecond) // simulate work

			panic("custom panic")
		}
	}
}
