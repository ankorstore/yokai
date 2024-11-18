package uuidv6

import (
	uuidv6test "github.com/ankorstore/yokai/generate/generatetest/uuidv6"
	"github.com/ankorstore/yokai/generate/uuidv6"
	"go.uber.org/fx"
)

// FxTestUuidV6GeneratorFactoryParam is used to retrieve the provided generate-test-uuid-V6-value from Fx.
type FxTestUuidV6GeneratorFactoryParam struct {
	fx.In
	Value string `name:"generate-test-uuid-v6-value"`
}

// TestUuidGeneratorV6Factory is a [uuidv6.Ui] implementation.
type TestUuidGeneratorV6Factory struct {
	value string
}

// NewFxTestUuidV6GeneratorFactory returns a new [TestUuidGeneratorV6Factory], implementing [uuidv6.UuidV6GeneratorFactory].
func NewFxTestUuidV6GeneratorFactory(p FxTestUuidV6GeneratorFactoryParam) uuidv6.UuidV6GeneratorFactory {
	return &TestUuidGeneratorV6Factory{
		value: p.Value,
	}
}

// Create returns a new [uuidv6.UuidV6Generator].
func (f *TestUuidGeneratorV6Factory) Create() uuidv6.UuidV6Generator {
	generator, err := uuidv6test.NewTestUuidV6Generator(f.value)
	if err != nil {
		return nil
	}

	return generator
}
