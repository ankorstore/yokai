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

	uuid1, err := generator.Generate()
	assert.NoError(t, err)

	uuid2, err := generator.Generate()
	assert.NoError(t, err)

	assert.NotEqual(t, uuid1.String(), uuid2.String())

	parsedUuid1, err := googleuuid.Parse(uuid1.String())
	assert.NoError(t, err)

	parsedUuid2, err := googleuuid.Parse(uuid2.String())
	assert.NoError(t, err)

	assert.NotEqual(t, parsedUuid1.String(), parsedUuid2.String())

	assert.Equal(t, uuid1.String(), parsedUuid1.String())
	assert.Equal(t, uuid2.String(), parsedUuid2.String())
}
