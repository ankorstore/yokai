package fxsql

import (
	"github.com/ankorstore/yokai/sql"
	"go.uber.org/fx"
)

// AsSQLHook registers a [sql.Hook] into Fx.
func AsSQLHook(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(sql.Hook)),
			fx.ResultTags(`group:"sql-hooks"`),
		),
	)
}

// AsSQLHooks registers a list of [sql.Hook] into Fx.
func AsSQLHooks(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsSQLHook(constructor))
	}

	return fx.Options(options...)
}

// AsSQLSeed registers a [Seed] into Fx.
func AsSQLSeed(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(Seed)),
			fx.ResultTags(`group:"sql-seeds"`),
		),
	)
}

// AsSQLSeeds registers a list of [Seed] into Fx.
func AsSQLSeeds(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsSQLSeed(constructor))
	}

	return fx.Options(options...)
}
