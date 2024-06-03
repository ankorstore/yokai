package uuid_test

import (
	"testing"

	uuidv7test "github.com/ankorstore/yokai/generate/generatetest/uuidv7"
	"github.com/ankorstore/yokai/generate/uuidv7"
	"github.com/stretchr/testify/assert"
)

const (
	uuid1 = "018fdd7d-1576-7a21-900e-1399637bd1a1"
	uuid2 = "018fdd7d-1576-76ff-944b-39bd474b0ea9"
	uuid3 = "018fdd7d-1576-7b53-a364-7b96dcc158c9"
)

func TestNewTestUuidV7Generator(t *testing.T) {
	t.Parallel()

	generator, err := uuidv7test.NewTestUuidV7Generator(uuid1)
	assert.NoError(t, err)

	assert.IsType(t, &uuidv7test.TestUuidV7Generator{}, generator)
	assert.Implements(t, (*uuidv7.UuidV7Generator)(nil), generator)
}

func TestGenerateSuccess(t *testing.T) {
	t.Parallel()

	generator, err := uuidv7test.NewTestUuidV7Generator(uuid2)
	assert.NoError(t, err)

	value1, err := generator.Generate()
	assert.NoError(t, err)

	value2, err := generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, uuid2, value1.String())
	assert.Equal(t, uuid2, value2.String())

	err = generator.SetValue(uuid3)
	assert.NoError(t, err)

	value1, err = generator.Generate()
	assert.NoError(t, err)

	value2, err = generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, uuid3, value1.String())
	assert.Equal(t, uuid3, value2.String())
}

func TestGenerateFailure(t *testing.T) {
	t.Parallel()

	_, err := uuidv7test.NewTestUuidV7Generator("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID length: 7")

	generator, err := uuidv7test.NewTestUuidV7Generator(uuid1)
	assert.NoError(t, err)

	err = generator.SetValue("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID length: 7")
}
