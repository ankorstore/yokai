package httpserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/stretchr/testify/assert"
)

func TestMatchPrefix(t *testing.T) {
	t.Parallel()

	prefixes := []string{
		"/foo",
		"/bar",
	}

	assert.True(t, httpserver.MatchPrefix(prefixes, "/foo"))
	assert.True(t, httpserver.MatchPrefix(prefixes, "/bar"))
	assert.True(t, httpserver.MatchPrefix(prefixes, "/fooo"))
	assert.True(t, httpserver.MatchPrefix(prefixes, "/barr"))
	assert.True(t, httpserver.MatchPrefix(prefixes, "/foo/foo"))
	assert.True(t, httpserver.MatchPrefix(prefixes, "/foo/bar"))
	assert.True(t, httpserver.MatchPrefix(prefixes, "/bar/bar"))
	assert.True(t, httpserver.MatchPrefix(prefixes, "/bar/foo"))

	assert.False(t, httpserver.MatchPrefix(prefixes, "/fo"))
	assert.False(t, httpserver.MatchPrefix(prefixes, "/ba"))
	assert.False(t, httpserver.MatchPrefix(prefixes, "/fo/foo"))
	assert.False(t, httpserver.MatchPrefix(prefixes, "/ba/bar"))
	assert.False(t, httpserver.MatchPrefix(prefixes, "/baz"))
}
