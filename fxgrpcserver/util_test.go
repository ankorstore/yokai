package fxgrpcserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "foo_bar", fxgrpcserver.Sanitize("foo-bar"))
	assert.Equal(t, "foo_bar", fxgrpcserver.Sanitize("foo bar"))
	assert.Equal(t, "foo_bar", fxgrpcserver.Sanitize("Foo-Bar"))
	assert.Equal(t, "foo_bar", fxgrpcserver.Sanitize("Foo Bar"))
}

func TestSplit(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"1", "2", "3"}, fxgrpcserver.Split("1,2,3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxgrpcserver.Split(" 1,2,3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxgrpcserver.Split("1,2,3 "))
	assert.Equal(t, []string{"1", "2", "3"}, fxgrpcserver.Split("1, 2, 3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxgrpcserver.Split(" 1, 2, 3 "))
}
