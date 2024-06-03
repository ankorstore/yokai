package fxgenerate

import (
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/generate/uuidv7"
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
		fx.Annotate(
			uuid.NewDefaultUuidGeneratorFactory,
			fx.As(new(uuid.UuidGeneratorFactory)),
		),
		fx.Annotate(
			uuidv7.NewDefaultUuidV7GeneratorFactory,
			fx.As(new(uuidv7.UuidV7GeneratorFactory)),
		),
		NewFxUuidGenerator,
		NewFxUuidV7Generator,
	),
)

// FxUuidGeneratorParam allows injection of the required dependencies in [NewFxUuidGenerator].
type FxUuidGeneratorParam struct {
	fx.In
	Factory uuid.UuidGeneratorFactory
}

// NewFxUuidGenerator returns a [uuid.UuidGenerator].
func NewFxUuidGenerator(p FxUuidGeneratorParam) uuid.UuidGenerator {
	return p.Factory.Create()
}

// FxUuidV7GeneratorParam allows injection of the required dependencies in [NewFxUuidV7Generator].
type FxUuidV7GeneratorParam struct {
	fx.In
	Factory uuidv7.UuidV7GeneratorFactory
}

// NewFxUuidV7Generator returns a [uuidv7.UuidV7Generator].
func NewFxUuidV7Generator(p FxUuidV7GeneratorParam) uuidv7.UuidV7Generator {
	return p.Factory.Create()
}
