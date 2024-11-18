package uuidv6_test

import (
	"testing"

	"github.com/ankorstore/yokai/generate/uuidv6"
	googleuuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultUuidV6Generator(t *testing.T) {
	t.Parallel()

	generator := uuidv6.NewDefaultUuidV6Generator()

	assert.IsType(t, &uuidv6.DefaultUuidV6Generator{}, generator)
	assert.Implements(t, (*uuidv6.UuidV6Generator)(nil), generator)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	generator := uuidv6.NewDefaultUuidV6Generator()

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
