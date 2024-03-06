package fxhttpserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "foo_bar", fxhttpserver.Sanitize("foo-bar"))
	assert.Equal(t, "foo_bar", fxhttpserver.Sanitize("foo bar"))
	assert.Equal(t, "foo_bar", fxhttpserver.Sanitize("Foo-Bar"))
	assert.Equal(t, "foo_bar", fxhttpserver.Sanitize("Foo Bar"))
}

func TestSplit(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"1", "2", "3"}, fxhttpserver.Split("1,2,3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxhttpserver.Split(" 1,2,3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxhttpserver.Split("1,2,3 "))
	assert.Equal(t, []string{"1", "2", "3"}, fxhttpserver.Split("1, 2, 3"))
	assert.Equal(t, []string{"1", "2", "3"}, fxhttpserver.Split(" 1, 2, 3 "))
}
