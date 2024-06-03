package uuidv7_test

import (
	"testing"

	"github.com/ankorstore/yokai/generate/uuidv7"
	googleuuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultUuidV7Generator(t *testing.T) {
	t.Parallel()

	generator := uuidv7.NewDefaultUuidV7Generator()

	assert.IsType(t, &uuidv7.DefaultUuidV7Generator{}, generator)
	assert.Implements(t, (*uuidv7.UuidV7Generator)(nil), generator)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	generator := uuidv7.NewDefaultUuidV7Generator()

	value1, err := generator.Generate()
	assert.NoError(t, err)

	value2, err := generator.Generate()
	assert.NoError(t, err)

	assert.NotEqual(t, value1, value2)

	parsedValue1, err := googleuuid.Parse(value1)
	assert.NoError(t, err)

	parsedValue2, err := googleuuid.Parse(value2)
	assert.NoError(t, err)

	assert.NotEqual(t, parsedValue1.String(), parsedValue2.String())

	assert.Equal(t, value1, parsedValue1.String())
	assert.Equal(t, value2, parsedValue2.String())
}
