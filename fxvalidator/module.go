package fxvalidator

import (
	"github.com/ankorstore/yokai/config"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

const (
	ModuleName = "validator"
	TagName    = "validate"
)

// FXValidatorModule is the [Fx] validator module.
//
// [Fx]: https://github.com/uber-go/fx
var FXValidatorModule = fx.Module(
	ModuleName,
	fx.Provide(
		ProvideValidator,
	),
)

// ProvideValidatorParams allows injection of the required dependencies in [ProvideValidator].
type ProvideValidatorParams struct {
	fx.In
	Config                      *config.Config
	AliasDefinitions            []AliasDefinition            `group:"validator-aliases"`
	ValidationDefinitions       []ValidationDefinition       `group:"validator-validations"`
	StructValidationDefinitions []StructValidationDefinition `group:"validator-struct-validations"`
	CustomTypeDefinitions       []CustomTypeDefinition       `group:"validator-custom-types"`
}

// ProvideValidator provides a new *validator.Validate instance.
func ProvideValidator(p ProvideValidatorParams) (*validator.Validate, error) {
	opts := []validator.Option{
		validator.WithRequiredStructEnabled(),
	}

	if p.Config.GetBool("modules.validator.private_fields") {
		opts = append(opts, validator.WithPrivateFieldValidation())
	}

	validate := validator.New(opts...)

	// tag name configuration
	tagName := p.Config.GetString("modules.validator.tag_name")
	if tagName == "" {
		tagName = TagName
	}

	validate.SetTagName(tagName)

	// aliases registration
	for _, def := range p.AliasDefinitions {
		validate.RegisterAlias(def.Alias(), def.Tags())
	}

	// var or struct field level validations registration
	for _, def := range p.ValidationDefinitions {
		err := validate.RegisterValidationCtx(def.Tag(), def.Fn(), def.CallEvenIfNull())
		if err != nil {
			return nil, err
		}
	}

	// struct level validations registration
	for _, def := range p.StructValidationDefinitions {
		validate.RegisterStructValidationCtx(def.Fn(), def.Types()...)
	}

	// custom types registration
	for _, def := range p.CustomTypeDefinitions {
		validate.RegisterCustomTypeFunc(def.Fn(), def.Types()...)
	}

	return validate, nil
}
