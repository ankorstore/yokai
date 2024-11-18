package uuidv6_test

import (
	"testing"

	"github.com/ankorstore/yokai/generate/uuidv6"
	googleuuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultUuidV6GeneratorFactory(t *testing.T) {
	t.Parallel()

	factory := uuidv6.NewDefaultUuidV6GeneratorFactory()

	assert.IsType(t, &uuidv6.DefaultUuidV6GeneratorFactory{}, factory)
	assert.Implements(t, (*uuidv6.UuidV6GeneratorFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	generator := uuidv6.NewDefaultUuidV6GeneratorFactory().Create()

	uuid1, err := generator.Generate()
	assert.NoError(t, err)

	uuid2, err := generator.Generate()
	assert.NoError(t, err)

	assert.NotEqual(t, uuid1, uuid2)

	parsedValue1, err := googleuuid.Parse(uuid1.String())
	assert.NoError(t, err)

	parsedValue2, err := googleuuid.Parse(uuid2.String())
	assert.NoError(t, err)

	assert.NotEqual(t, parsedValue1.String(), parsedValue2.String())

	assert.Equal(t, uuid1.String(), parsedValue1.String())
	assert.Equal(t, uuid2.String(), parsedValue2.String())
}
