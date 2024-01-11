package fxgenerate

import (
	"github.com/ankorstore/yokai/generate/uuid"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "generate"

// FxGenerateModule is the [Fx] generate module.
//
// [Fx]: https://github.com/uber-go/fx
var FxGenerateModule = fx.Module(
	ModuleName,
	fx.Provide(
		uuid.NewDefaultUuidGeneratorFactory,
		NewFxUuidGenerator,
	),
)

// FxUuidGeneratorParam allows injection of the required dependencies in [NewFxUuidGenerator].
type FxUuidGeneratorParam struct {
	fx.In
	Factory uuid.UuidGeneratorFactory
}

// NewFxUuidGenerator returns a [uuid.UuidGenerator].
func NewFxUuidGenerator(p FxUuidGeneratorParam) (uuid.UuidGenerator, error) {
	return p.Factory.Create(), nil
}
