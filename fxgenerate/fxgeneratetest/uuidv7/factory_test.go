package uuidv7_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgenerate"
	fxgeneratetestuuidv7 "github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuidv7"
	testuuidv7 "github.com/ankorstore/yokai/fxgenerate/testdata/uuidv7"
	"github.com/ankorstore/yokai/generate/uuidv7"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestTestUuidV7GeneratorSuccess(t *testing.T) {
	t.Parallel()

	var generator uuidv7.UuidV7Generator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Provide(
			fx.Annotate(
				func() string {
					return testuuidv7.TestUUIDV7
				},
				fx.ResultTags(`name:"generate-test-uuid-v7-value"`),
			),
		),
		fx.Decorate(fxgeneratetestuuidv7.NewFxTestUuidV7GeneratorFactory),
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	value, err := generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, testuuidv7.TestUUIDV7, value.String())
}

func TestTestUuidV7GeneratorError(t *testing.T) {
	t.Parallel()

	var generator uuidv7.UuidV7Generator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Provide(
			fx.Annotate(
				func() string {
					return "invalid"
				},
				fx.ResultTags(`name:"generate-test-uuid-v7-value"`),
			),
		),
		fx.Decorate(fxgeneratetestuuidv7.NewFxTestUuidV7GeneratorFactory),
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	assert.Nil(t, generator)
}
