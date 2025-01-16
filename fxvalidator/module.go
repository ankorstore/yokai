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

var FXValidatorModule = fx.Module(
	ModuleName,
	fx.Provide(
		ProvideValidator,
	),
)

type ProvideValidatorParams struct {
	fx.In
	Config           *config.Config
	AliasDefinitions []AliasDefinition `group:"validator-aliases"`
}

func ProvideValidator(p ProvideValidatorParams) *validator.Validate {
	opts := []validator.Option{
		validator.WithRequiredStructEnabled(),
	}

	if p.Config.GetBool("modules.validator.private_fields") {
		opts = append(opts, validator.WithPrivateFieldValidation())
	}

	validate := validator.New(opts...)

	// tag name
	tagName := p.Config.GetString("modules.validator.tag_name")
	if tagName == "" {
		tagName = TagName
	}

	validate.SetTagName(tagName)

	// aliases
	for _, def := range p.AliasDefinitions {
		validate.RegisterAlias(def.Alias(), def.Tags())
	}

	return validate
}
