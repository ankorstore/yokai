package fxgenerate_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgenerate"
	testuuid "github.com/ankorstore/yokai/fxgenerate/testdata/uuid"
	"github.com/ankorstore/yokai/generate/uuid"
	googleuuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModuleUuidGenerator(t *testing.T) {
	t.Parallel()

	var generator uuid.UuidGenerator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Populate(&generator),
	).RequireStart().RequireStop()

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

func TestModuleUuidGeneratorDecoration(t *testing.T) {
	t.Parallel()

	var generator uuid.UuidGenerator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Decorate(testuuid.NewTestStaticUuidGeneratorFactory),
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	assert.Equal(t, "static", generator.Generate())
}
