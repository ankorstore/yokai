package fxvalidator

import (
	"go.uber.org/fx"
)

// AsAlias registers a validation alias.
func AsAlias(alias string, tags string) fx.Option {
	return fx.Supply(
		fx.Annotate(
			NewAliasDefinition(alias, tags),
			fx.As(new(AliasDefinition)),
			fx.ResultTags(`group:"validator-aliases"`),
		),
	)
}
