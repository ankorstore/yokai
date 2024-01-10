package uuid_test

import (
	"testing"

	"github.com/ankorstore/yokai/generate/uuid"
	googleuuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultUuidGeneratorFactory(t *testing.T) {
	t.Parallel()

	factory := uuid.NewDefaultUuidGeneratorFactory()

	assert.IsType(t, &uuid.DefaultUuidGeneratorFactory{}, factory)
	assert.Implements(t, (*uuid.UuidGeneratorFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	generator := uuid.NewDefaultUuidGeneratorFactory().Create()

	value1 := generator.Generate()
	value2 := generator.Generate()

	assert.NotEqual(t, value1, value2)

	parsedValue1, err := googleuuid.Parse(value1)
	assert.NoError(t, err)

	parsedValue2, err := googleuuid.Parse(value2)
	assert.NoError(t, err)

	assert.NotEqual(t, parsedValue1.String(), parsedValue2.String())

	assert.Equal(t, value1, parsedValue1.String())
	assert.Equal(t, value2, parsedValue2.String())
}
