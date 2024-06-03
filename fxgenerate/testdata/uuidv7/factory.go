package uuidv7

import (
	uuidtest "github.com/ankorstore/yokai/generate/generatetest/uuidv7"
	"github.com/ankorstore/yokai/generate/uuidv7"
)

const TestUUIDV7 = "018fdd68-1b41-7eb0-afad-57f45297c7c1"

type TestStaticUuidV7GeneratorFactory struct{}

func NewTestStaticUuidV7GeneratorFactory() uuidv7.UuidV7GeneratorFactory {
	return &TestStaticUuidV7GeneratorFactory{}
}

func (f *TestStaticUuidV7GeneratorFactory) Create() uuidv7.UuidV7Generator {
	//nolint:errcheck
	generator, _ := uuidtest.NewTestUuidV7Generator(TestUUIDV7)

	return generator
}
