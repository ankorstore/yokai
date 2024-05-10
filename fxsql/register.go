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
