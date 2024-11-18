package uuidv6_test

import (
	"testing"

	uuidv6test "github.com/ankorstore/yokai/generate/generatetest/uuidv6"
	"github.com/ankorstore/yokai/generate/uuidv6"
	"github.com/stretchr/testify/assert"
)

const (
	uuid1 = "1efa59f2-d438-6ec0-9d52-4da3ad16f2c6"
	uuid2 = "1efa59f2-d438-6ec1-8b52-6844309a22de"
	uuid3 = "1efa59f2-d438-6ec2-abd0-e36a967ab868"
)

func TestNewTestUuidV6Generator(t *testing.T) {
	t.Parallel()

	generator, err := uuidv6test.NewTestUuidV6Generator(uuid1)
	assert.NoError(t, err)

	assert.IsType(t, &uuidv6test.TestUuidV6Generator{}, generator)
	assert.Implements(t, (*uuidv6.UuidV6Generator)(nil), generator)
}

func TestGenerateSuccess(t *testing.T) {
	t.Parallel()

	generator, err := uuidv6test.NewTestUuidV6Generator(uuid2)
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

	_, err := uuidv6test.NewTestUuidV6Generator("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID length: 7")

	generator, err := uuidv6test.NewTestUuidV6Generator(uuid1)
	assert.NoError(t, err)

	err = generator.SetValue("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID length: 7")
}
