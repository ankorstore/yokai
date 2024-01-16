package workers

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/worker"
)

type OneShotWorker struct{}

func NewOneShotWorker() *OneShotWorker {
	return &OneShotWorker{}
}

func (w *OneShotWorker) Name() string {
	return "OneShotWorker"
}

func (w *OneShotWorker) Run(ctx context.Context) error {
	ctx, span := worker.CtxTracer(ctx).Start(ctx, "one shot span")
	defer span.End()

	logger := worker.CtxLogger(ctx)

	logger.Info().Msgf("running worker %s [id %s]", worker.CtxWorkerName(ctx), worker.CtxWorkerExecutionId(ctx))

	time.Sleep(10 * time.Millisecond) // simulate work

	logger.Info().Msg("stopped")

	return nil
}
