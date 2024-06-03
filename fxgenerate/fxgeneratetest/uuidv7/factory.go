package uuidv7

import (
	uuidv7test "github.com/ankorstore/yokai/generate/generatetest/uuidv7"
	"github.com/ankorstore/yokai/generate/uuidv7"
	"go.uber.org/fx"
)

// FxTestUuidV7GeneratorFactoryParam is used to retrieve the provided generate-test-uuid-V7-value from Fx.
type FxTestUuidV7GeneratorFactoryParam struct {
	fx.In
	Value string `name:"generate-test-uuid-v7-value"`
}

// TestUuidGeneratorV7Factory is a [uuidv7.Ui] implementation.
type TestUuidGeneratorV7Factory struct {
	value string
}

// NewFxTestUuidV7GeneratorFactory returns a new [TestUuidGeneratorV7Factory], implementing [uuidv7.UuidV7GeneratorFactory].
func NewFxTestUuidV7GeneratorFactory(p FxTestUuidV7GeneratorFactoryParam) uuidv7.UuidV7GeneratorFactory {
	return &TestUuidGeneratorV7Factory{
		value: p.Value,
	}
}

// Create returns a new [uuidv7.UuidV7Generator].
func (f *TestUuidGeneratorV7Factory) Create() uuidv7.UuidV7Generator {
	generator, err := uuidv7test.NewTestUuidV7Generator(f.value)
	if err != nil {
		return nil
	}

	return generator
}
