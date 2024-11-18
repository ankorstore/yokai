//nolint:dupl
package fxgenerate_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgenerate"
	testuuid "github.com/ankorstore/yokai/fxgenerate/testdata/uuid"
	testuuidv6 "github.com/ankorstore/yokai/fxgenerate/testdata/uuidv6"
	testuuidv7 "github.com/ankorstore/yokai/fxgenerate/testdata/uuidv7"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/generate/uuidv6"
	"github.com/ankorstore/yokai/generate/uuidv7"
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

func TestModuleUuidV6Generator(t *testing.T) {
	t.Parallel()

	var generator uuidv6.UuidV6Generator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	value1, err := generator.Generate()
	assert.NoError(t, err)

	value2, err := generator.Generate()
	assert.NoError(t, err)

	assert.NotEqual(t, value1, value2)

	parsedValue1, err := googleuuid.Parse(value1.String())
	assert.NoError(t, err)

	parsedValue2, err := googleuuid.Parse(value2.String())
	assert.NoError(t, err)

	assert.NotEqual(t, parsedValue1.String(), parsedValue2.String())

	assert.Equal(t, value1.String(), parsedValue1.String())
	assert.Equal(t, value2.String(), parsedValue2.String())
}

func TestModuleUuidV7Generator(t *testing.T) {
	t.Parallel()

	var generator uuidv7.UuidV7Generator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	value1, err := generator.Generate()
	assert.NoError(t, err)

	value2, err := generator.Generate()
	assert.NoError(t, err)

	assert.NotEqual(t, value1, value2)

	parsedValue1, err := googleuuid.Parse(value1.String())
	assert.NoError(t, err)

	parsedValue2, err := googleuuid.Parse(value2.String())
	assert.NoError(t, err)

	assert.NotEqual(t, parsedValue1.String(), parsedValue2.String())

	assert.Equal(t, value1.String(), parsedValue1.String())
	assert.Equal(t, value2.String(), parsedValue2.String())
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

func TestModuleUuidV6GeneratorDecoration(t *testing.T) {
	t.Parallel()

	var generator uuidv6.UuidV6Generator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Decorate(testuuidv6.NewTestStaticUuidV6GeneratorFactory),
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	value, err := generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, testuuidv6.TestUUIDV6, value.String())
}

func TestModuleUuidV7GeneratorDecoration(t *testing.T) {
	t.Parallel()

	var generator uuidv7.UuidV7Generator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Decorate(testuuidv7.NewTestStaticUuidV7GeneratorFactory),
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	value, err := generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, testuuidv7.TestUUIDV7, value.String())
}
