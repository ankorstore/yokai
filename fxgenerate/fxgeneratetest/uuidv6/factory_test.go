package uuidv6_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgenerate"
	fxgeneratetestuuidv6 "github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuidv6"
	testuuidv6 "github.com/ankorstore/yokai/fxgenerate/testdata/uuidv6"
	"github.com/ankorstore/yokai/generate/uuidv6"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestTestUuidV6GeneratorSuccess(t *testing.T) {
	t.Parallel()

	var generator uuidv6.UuidV6Generator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Provide(
			fx.Annotate(
				func() string {
					return testuuidv6.TestUUIDV6
				},
				fx.ResultTags(`name:"generate-test-uuid-v6-value"`),
			),
		),
		fx.Decorate(fxgeneratetestuuidv6.NewFxTestUuidV6GeneratorFactory),
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	value, err := generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, testuuidv6.TestUUIDV6, value.String())
}

func TestTestUuidV6GeneratorError(t *testing.T) {
	t.Parallel()

	var generator uuidv6.UuidV6Generator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Provide(
			fx.Annotate(
				func() string {
					return "invalid"
				},
				fx.ResultTags(`name:"generate-test-uuid-v6-value"`),
			),
		),
		fx.Decorate(fxgeneratetestuuidv6.NewFxTestUuidV6GeneratorFactory),
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	assert.Nil(t, generator)
}
