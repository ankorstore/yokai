package fxcore_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/stretchr/testify/assert"
)

func TestAsCoreExtraInfo(t *testing.T) {
	t.Parallel()

	result := fxcore.AsCoreExtraInfo("foo", "bar")

	assert.Equal(t, "fx.supplyOption", fmt.Sprintf("%T", result))
}
