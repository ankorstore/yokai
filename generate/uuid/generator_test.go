package uuid_test

import (
	"testing"

	"github.com/ankorstore/yokai/generate/uuid"
	googleuuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultUuidGenerator(t *testing.T) {
	t.Parallel()

	generator := uuid.NewDefaultUuidGenerator()

	assert.IsType(t, &uuid.DefaultUuidGenerator{}, generator)
	assert.Implements(t, (*uuid.UuidGenerator)(nil), generator)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	generator := uuid.NewDefaultUuidGenerator()

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
