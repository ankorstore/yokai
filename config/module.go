package config

import (
	"os"

	"go.uber.org/fx"
)

// ModuleName is the config module name.
const ModuleName = "config"

// ConfigModule is the Yokai config module.
var ConfigModule = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(
			NewDefaultConfigFactory,
			fx.As(new(ConfigFactory)),
		),
		ProvideConfig,
	),
)

// ProvideConfigParams allows injection of the required dependencies in [ProvideConfig].
type ProvideConfigParams struct {
	fx.In
	Factory ConfigFactory
}

// ProvideConfig provides a [Config] instance.
func ProvideConfig(p ProvideConfigParams) (*Config, error) {
	return p.Factory.Create(
		WithFileName("config"),
		WithFilePaths(os.Getenv("APP_CONFIG_PATH")),
	)
}
