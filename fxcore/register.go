package fxcore

import (
	"go.uber.org/fx"
)

// AsCoreExtraInfo registers extra information in the core.
func AsCoreExtraInfo(name string, value string) fx.Option {
	return fx.Supply(
		fx.Annotate(
			NewFxExtraInfo(name, value),
			fx.As(new(FxExtraInfo)),
			fx.ResultTags(`group:"core-extra-infos"`),
		),
	)
}

// AsTask registers a task in the core.
func AsTask(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(Task)),
			fx.ResultTags(`group:"core-tasks"`),
		),
	)
}

// AsTasks registers several tasks in the core.
func AsTasks(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsTask(constructor))
	}

	return fx.Options(options...)
}
