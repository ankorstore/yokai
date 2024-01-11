package uuid_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgenerate"
	fxgeneratetestuuid "github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuid"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestTestUuidGenerator(t *testing.T) {
	t.Parallel()

	var generator uuid.UuidGenerator

	fxtest.New(
		t,
		fx.NopLogger,
		fxgenerate.FxGenerateModule,
		fx.Provide(
			fx.Annotate(
				func() string {
					return "some test value"
				},
				fx.ResultTags(`name:"generate-test-uuid-value"`),
			),
		),
		fx.Decorate(fxgeneratetestuuid.NewFxTestUuidGeneratorFactory),
		fx.Populate(&generator),
	).RequireStart().RequireStop()

	assert.Equal(t, "some test value", generator.Generate())
}
