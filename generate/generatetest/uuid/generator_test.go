package uuid_test

import (
	"testing"

	uuidtest "github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTestUuidGenerator(t *testing.T) {
	t.Parallel()

	generator := uuidtest.NewTestUuidGenerator("random")

	assert.IsType(t, &uuidtest.TestUuidGenerator{}, generator)
	assert.Implements(t, (*uuid.UuidGenerator)(nil), generator)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	generator := uuidtest.NewTestUuidGenerator("test")

	value1 := generator.Generate()
	value2 := generator.Generate()

	assert.Equal(t, "test", value1)
	assert.Equal(t, "test", value2)

	generator.SetValue("other test")

	value1 = generator.Generate()
	value2 = generator.Generate()

	assert.Equal(t, "other test", value1)
	assert.Equal(t, "other test", value2)
}
