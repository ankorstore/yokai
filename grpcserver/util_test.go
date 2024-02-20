package grpcserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/grpcserver"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Parallel()

	list := []string{
		"/foo",
		"/bar",
		"/baz",
	}

	assert.True(t, grpcserver.Contains(list, "/foo"))
	assert.True(t, grpcserver.Contains(list, "/bar"))
	assert.True(t, grpcserver.Contains(list, "/baz"))

	assert.False(t, grpcserver.Contains(list, "/fo"))
	assert.False(t, grpcserver.Contains(list, "/ba"))
	assert.False(t, grpcserver.Contains(list, "/invalid"))
}

func TestUnique(t *testing.T) {
	t.Parallel()

	list := []string{
		"/foo",
		"/bar",
		"/baz",
		"/bar",
		"/baz",
		"/fo",
		"/ba",
		"/ba",
	}

	uniqueList := grpcserver.Unique(list)

	assert.Len(t, uniqueList, 5)
	assert.Equal(
		t,
		[]string{
			"/foo",
			"/bar",
			"/baz",
			"/fo",
			"/ba",
		},
		uniqueList,
	)
}
