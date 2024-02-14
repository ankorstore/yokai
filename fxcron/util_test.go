package fxcron_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Parallel()

	list := []string{
		"/foo",
		"/bar",
		"/baz",
	}

	assert.True(t, fxcron.Contains(list, "/foo"))
	assert.True(t, fxcron.Contains(list, "/bar"))
	assert.True(t, fxcron.Contains(list, "/baz"))

	assert.False(t, fxcron.Contains(list, "/fo"))
	assert.False(t, fxcron.Contains(list, "/ba"))
	assert.False(t, fxcron.Contains(list, "/invalid"))
}

func TestSanitize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "foo_bar", fxcron.Sanitize("foo-bar"))
	assert.Equal(t, "foo_bar", fxcron.Sanitize("foo bar"))
	assert.Equal(t, "foo_bar", fxcron.Sanitize("Foo-Bar"))
	assert.Equal(t, "foo_bar", fxcron.Sanitize("Foo Bar"))
}
