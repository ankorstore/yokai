package uuid

import (
	uuidtest "github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/generate/uuid"
	"go.uber.org/fx"
)

// FxTestUuidGeneratorFactoryParam is used to retrieve the provided generate-test-uuid-value from Fx.
type FxTestUuidGeneratorFactoryParam struct {
	fx.In
	Value string `name:"generate-test-uuid-value"`
}

// TestUuidGeneratorFactory is a [uuid.UuidGeneratorFactory] implementation.
type TestUuidGeneratorFactory struct {
	value string
}

// NewFxTestUuidGeneratorFactory returns a new [TestUuidGeneratorFactory], implementing [uuid.UuidGeneratorFactory].
func NewFxTestUuidGeneratorFactory(p FxTestUuidGeneratorFactoryParam) uuid.UuidGeneratorFactory {
	return &TestUuidGeneratorFactory{
		value: p.Value,
	}
}

// Create returns a new [uuid.UuidGenerator].
func (f *TestUuidGeneratorFactory) Create() uuid.UuidGenerator {
	return uuidtest.NewTestUuidGenerator(f.value)
}
