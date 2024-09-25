package fxconfig_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
)

func TestAsConfigPath(t *testing.T) {
	t.Parallel()

	result := fxconfig.AsConfigPath("foo")

	assert.Equal(t, "fx.supplyOption", fmt.Sprintf("%T", result))
}
