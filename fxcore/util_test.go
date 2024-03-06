package fxcore_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "foo_bar", fxcore.Sanitize("foo-bar"))
	assert.Equal(t, "foo_bar", fxcore.Sanitize("foo bar"))
	assert.Equal(t, "foo_bar", fxcore.Sanitize("Foo-Bar"))
	assert.Equal(t, "foo_bar", fxcore.Sanitize("Foo Bar"))
}

func TestSplit(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"1", "2", "3"}, fxcore.Split("1,2,3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxcore.Split(" 1,2,3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxcore.Split("1,2,3 "))
	assert.Equal(t, []string{"1", "2", "3"}, fxcore.Split("1, 2, 3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxcore.Split(" 1, 2, 3 "))
}
