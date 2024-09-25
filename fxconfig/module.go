package fxconfig

import (
	"os"

	"github.com/ankorstore/yokai/config"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "config"

// FxConfigModule is the [Fx] config module.
//
// [Fx]: https://github.com/uber-go/fx
var FxConfigModule = fx.Module(
	ModuleName,
	fx.Provide(
		config.NewDefaultConfigFactory,
		NewFxConfig,
	),
)

// FxConfigParam allows injection of the required dependencies in [NewFxConfig].
type FxConfigParam struct {
	fx.In
	Factory     config.ConfigFactory
	ConfigPaths []string `group:"config-paths"`
}

// NewFxConfig returns a [config.Config].
func NewFxConfig(p FxConfigParam) (*config.Config, error) {
	configFilePaths := append([]string{os.Getenv("APP_CONFIG_PATH")}, p.ConfigPaths...)

	return p.Factory.Create(
		config.WithFileName("config"),
		config.WithFilePaths(configFilePaths...),
	)
}
