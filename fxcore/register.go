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
