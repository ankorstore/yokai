package fxsql_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxsql/testdata/hook"
	"github.com/ankorstore/yokai/fxsql/testdata/seed"
	"github.com/stretchr/testify/assert"
)

func TestAsSQLHook(t *testing.T) {
	t.Parallel()

	result := fxsql.AsSQLHook(hook.NewTestHook)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", result))
}

func TestAsSQLHooks(t *testing.T) {
	t.Parallel()

	result := fxsql.AsSQLHooks(hook.NewTestHook)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}

func TestAsSQLSeed(t *testing.T) {
	t.Parallel()

	result := fxsql.AsSQLSeed(seed.NewTestSeed)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", result))
}

func TestAsSQLSeeds(t *testing.T) {
	t.Parallel()

	result := fxsql.AsSQLSeeds(seed.NewTestSeed)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
