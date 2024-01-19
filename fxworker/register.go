package fxworker

import (
	"github.com/ankorstore/yokai/worker"
	"go.uber.org/fx"
)

// AsWorker registers a [worker.Worker] into Fx, with an optional list of [worker.WorkerExecutionOption].
func AsWorker(w any, options ...worker.WorkerExecutionOption) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				w,
				fx.As(new(worker.Worker)),
				fx.ResultTags(`group:"workers"`),
			),
		),
		fx.Supply(
			fx.Annotate(
				NewWorkerDefinition(GetReturnType(w), options...),
				fx.As(new(WorkerDefinition)),
				fx.ResultTags(`group:"workers-definitions"`),
			),
		),
	)
}
