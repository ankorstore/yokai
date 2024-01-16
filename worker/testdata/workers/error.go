package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/ankorstore/yokai/worker"
)

type ErrorWorker struct{}

func NewErrorWorker() *ErrorWorker {
	return &ErrorWorker{}
}

func (w *ErrorWorker) Name() string {
	return "ErrorWorker"
}

func (w *ErrorWorker) Run(ctx context.Context) error {
	logger := worker.CtxLogger(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			logger.Info().Msg("running")

			time.Sleep(10 * time.Millisecond) // simulate work

			return fmt.Errorf("custom error")
		}
	}
}
