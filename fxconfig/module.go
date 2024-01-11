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
	Factory config.ConfigFactory
}

// NewFxConfig returns a [config.Config].
func NewFxConfig(p FxConfigParam) (*config.Config, error) {
	return p.Factory.Create(
		config.WithFileName("config"),
		config.WithFilePaths(
			".",
			"./configs",
			os.Getenv("APP_CONFIG_PATH"),
		),
	)
}
