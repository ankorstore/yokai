package fxworker

import (
	"reflect"

	"github.com/ankorstore/yokai/worker"
	"go.uber.org/fx"
)

// AsWorker registers a [worker.Worker] into Fx, with an optional list of [worker.WorkerExecutionOption] and middlewares.
func AsWorker(w any, options ...any) fx.Option {
	var providers []any
	var workerExecutionOptions []worker.WorkerExecutionOption
	var middlewareDefs []MiddlewareDefinition

	for _, o := range options {
		if option, ok := o.(worker.WorkerExecutionOption); ok {
			workerExecutionOptions = append(workerExecutionOptions, option)

			continue
		}

		if IsConcreteMiddleware(o) {
			middlewareDefs = append(middlewareDefs, NewMiddlewareDefinition(GetType(o)))
		} else {
			providers = append(
				providers,
				fx.Annotate(
					o,
					fx.As(new(worker.Middleware)),
					fx.ResultTags(`group:"worker-middlewares"`),
				),
			)

			middlewareDefs = append(middlewareDefs, NewMiddlewareDefinition(GetReturnType(o)))
		}
	}

	return fx.Options(
		fx.Provide(append(
			[]any{
				fx.Annotate(
					w,
					fx.As(new(worker.Worker)),
					fx.ResultTags(`group:"workers"`),
				),
			},
			providers...,
		)...),
		fx.Supply(
			fx.Annotate(
				NewWorkerDefinitionWithMiddlewares(GetReturnType(w), middlewareDefs, workerExecutionOptions...),
				fx.As(new(WorkerDefinition)),
				fx.ResultTags(`group:"workers-definitions"`),
			),
		),
	)
}

// IsConcreteMiddleware returns true if the middleware is a concrete [worker.MiddlewareFunc] implementation.
func IsConcreteMiddleware(middleware any) bool {
	return reflect.TypeOf(middleware).ConvertibleTo(reflect.TypeOf(worker.MiddlewareFunc(nil)))
}
