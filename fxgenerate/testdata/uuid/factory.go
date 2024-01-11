package uuid

import (
	uuidtest "github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/generate/uuid"
)

type TestStaticUuidGeneratorFactory struct{}

func NewTestStaticUuidGeneratorFactory() uuid.UuidGeneratorFactory {
	return &TestStaticUuidGeneratorFactory{}
}

func (f *TestStaticUuidGeneratorFactory) Create() uuid.UuidGenerator {
	return uuidtest.NewTestUuidGenerator("static")
}
