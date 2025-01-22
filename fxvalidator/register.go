package fxvalidator

import (
	"github.com/go-playground/validator/v10"
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

// AsValidation registers a var or struct field level validation.
func AsValidation(tag string, fn validator.FuncCtx, callEvenIfNull bool) fx.Option {
	return fx.Supply(
		fx.Annotate(
			NewValidationDefinition(tag, fn, callEvenIfNull),
			fx.As(new(ValidationDefinition)),
			fx.ResultTags(`group:"validator-validations"`),
		),
	)
}

// AsStructValidation registers a struct level validation.
func AsStructValidation(fn validator.StructLevelFuncCtx, types ...any) fx.Option {
	return fx.Supply(
		fx.Annotate(
			NewStructValidationDefinition(fn, types...),
			fx.As(new(StructValidationDefinition)),
			fx.ResultTags(`group:"validator-struct-validations"`),
		),
	)
}

// AsCustomType registers a custom type.
func AsCustomType(fn validator.CustomTypeFunc, types ...any) fx.Option {
	return fx.Supply(
		fx.Annotate(
			NewCustomTypeDefinition(fn, types...),
			fx.As(new(CustomTypeDefinition)),
			fx.ResultTags(`group:"validator-custom-types"`),
		),
	)
}
