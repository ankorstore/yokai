package fxcore_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/stretchr/testify/assert"
)

func TestRootDir(t *testing.T) {
	t.Parallel()

	dir := fxcore.RootDir(99)
	assert.Equal(t, "..", dir)
}
