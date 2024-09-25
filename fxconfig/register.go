package fxconfig

import (
	"go.uber.org/fx"
)

// AsConfigPath registers an additional config files lookup path.
func AsConfigPath(path string) fx.Option {
	return fx.Supply(
		fx.Annotate(
			path,
			fx.ResultTags(`group:"config-paths"`),
		),
	)
}
