package fxcore_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxcore/testdata/tasks"
	"github.com/stretchr/testify/assert"
)

func TestAsCoreExtraInfo(t *testing.T) {
	t.Parallel()

	result := fxcore.AsCoreExtraInfo("foo", "bar")

	assert.Equal(t, "fx.supplyOption", fmt.Sprintf("%T", result))
}

func TestAsTask(t *testing.T) {
	t.Parallel()

	result := fxcore.AsTask(tasks.NewErrorTask)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", result))
}

func TestAsTasks(t *testing.T) {
	t.Parallel()

	result := fxcore.AsTasks(tasks.NewErrorTask)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
