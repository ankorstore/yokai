package uuidv6

import (
	uuidtest "github.com/ankorstore/yokai/generate/generatetest/uuidv6"
	"github.com/ankorstore/yokai/generate/uuidv6"
)

const TestUUIDV6 = "1efa5a47-e5d2-6d70-8953-776e81422ff3"

type TestStaticUuidV6GeneratorFactory struct{}

func NewTestStaticUuidV6GeneratorFactory() uuidv6.UuidV6GeneratorFactory {
	return &TestStaticUuidV6GeneratorFactory{}
}

func (f *TestStaticUuidV6GeneratorFactory) Create() uuidv6.UuidV6Generator {
	//nolint:errcheck
	generator, _ := uuidtest.NewTestUuidV6Generator(TestUUIDV6)

	return generator
}
