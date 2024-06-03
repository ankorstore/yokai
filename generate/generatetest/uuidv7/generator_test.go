package uuid_test

import (
	"testing"

	uuidv7test "github.com/ankorstore/yokai/generate/generatetest/uuidv7"
	"github.com/ankorstore/yokai/generate/uuidv7"
	"github.com/stretchr/testify/assert"
)

func TestNewTestUuidV7Generator(t *testing.T) {
	t.Parallel()

	generator := uuidv7test.NewTestUuidV7Generator("random")

	assert.IsType(t, &uuidv7test.TestUuidV7Generator{}, generator)
	assert.Implements(t, (*uuidv7.UuidV7Generator)(nil), generator)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	generator := uuidv7test.NewTestUuidV7Generator("test")

	value1, err := generator.Generate()
	assert.NoError(t, err)

	value2, err := generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, "test", value1)
	assert.Equal(t, "test", value2)

	generator.SetValue("other test")

	value1, err = generator.Generate()
	assert.NoError(t, err)

	value2, err = generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, "other test", value1)
	assert.Equal(t, "other test", value2)
}
