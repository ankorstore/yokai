package fxvalidator

import "github.com/go-playground/validator/v10"

// AliasDefinition is the interface for aliases definitions.
type AliasDefinition interface {
	Alias() string
	Tags() string
}

type aliasDefinition struct {
	alias string
	tags  string
}

// NewAliasDefinition returns a new AliasDefinition.
func NewAliasDefinition(alias string, tags string) AliasDefinition {
	return &aliasDefinition{
		alias: alias,
		tags:  tags,
	}
}

// Alias returns the alias definition alias.
func (d *aliasDefinition) Alias() string {
	return d.alias
}

// Tags returns the alias definition tags.
func (d *aliasDefinition) Tags() string {
	return d.tags
}

// ValidationDefinition is the interface for validations definitions.
type ValidationDefinition interface {
	Tag() string
	Fn() validator.FuncCtx
	CallEvenIfNull() bool
}

type validationDefinition struct {
	tag            string
	fn             validator.FuncCtx
	callEvenIfNull bool
}

// NewValidationDefinition returns a new ValidationDefinition.
func NewValidationDefinition(tag string, fn validator.FuncCtx, callEvenIfNull bool) ValidationDefinition {
	return &validationDefinition{
		tag:            tag,
		fn:             fn,
		callEvenIfNull: callEvenIfNull,
	}
}

// Tag returns the validation definition tag.
func (d *validationDefinition) Tag() string {
	return d.tag
}

// Fn returns the validation definition validator.FuncCtx.
func (d *validationDefinition) Fn() validator.FuncCtx {
	return d.fn
}

// CallEvenIfNull is true if the validation definition must be called even if null.
func (d *validationDefinition) CallEvenIfNull() bool {
	return d.callEvenIfNull
}

// StructValidationDefinition is the interface for struct validations definitions.
type StructValidationDefinition interface {
	Fn() validator.StructLevelFuncCtx
	Types() []any
}

type structValidationDefinition struct {
	fn    validator.StructLevelFuncCtx
	types []any
}

// NewStructValidationDefinition returns a new StructValidationDefinition.
func NewStructValidationDefinition(fn validator.StructLevelFuncCtx, types ...any) StructValidationDefinition {
	return &structValidationDefinition{
		fn:    fn,
		types: types,
	}
}

// Fn returns the struct validation definition validator.StructLevelFuncCtx.
func (d *structValidationDefinition) Fn() validator.StructLevelFuncCtx {
	return d.fn
}

// Types returns the struct validation definition types.
func (d *structValidationDefinition) Types() []any {
	return d.types
}

// CustomTypeDefinition is the interface for custom types definitions.
type CustomTypeDefinition interface {
	Fn() validator.CustomTypeFunc
	Types() []any
}

type customTypeDefinition struct {
	fn    validator.CustomTypeFunc
	types []any
}

// NewCustomTypeDefinition returns a new CustomTypeDefinition.
func NewCustomTypeDefinition(fn validator.CustomTypeFunc, types ...any) CustomTypeDefinition {
	return &customTypeDefinition{
		fn:    fn,
		types: types,
	}
}

// Fn returns the custom type definition validator.StructLevelFuncCtx.
func (d *customTypeDefinition) Fn() validator.CustomTypeFunc {
	return d.fn
}

// Types returns the custom type definition types.
func (d *customTypeDefinition) Types() []any {
	return d.types
}
