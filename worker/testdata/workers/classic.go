package workers

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/worker"
)

type ClassicWorker struct{}

func NewClassicWorker() *ClassicWorker {
	return &ClassicWorker{}
}

func (w *ClassicWorker) Name() string {
	return "ClassicWorker"
}

func (w *ClassicWorker) Run(ctx context.Context) error {
	ctx, span := worker.CtxTracer(ctx).Start(ctx, "one shot span")
	defer span.End()

	logger := worker.CtxLogger(ctx)

	logger.Info().Msgf("running worker %s [id %s]", worker.CtxWorkerName(ctx), worker.CtxWorkerExecutionId(ctx))

	time.Sleep(10 * time.Millisecond) // simulate work

	logger.Info().Msg("stopped")

	return nil
}
