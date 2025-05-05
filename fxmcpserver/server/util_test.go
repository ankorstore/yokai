package server_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/stretchr/testify/assert"
)

func TestFuncName(t *testing.T) {
	t.Parallel()

	fn := func() {}

	assert.Equal(t, "github.com/ankorstore/yokai/fxmcpserver/server_test.TestFuncName.func1", server.FuncName(fn))
}

func TestSanitize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "foo_bar", server.Sanitize("foo-bar"))
	assert.Equal(t, "foo_bar", server.Sanitize("foo bar"))
	assert.Equal(t, "foo_bar", server.Sanitize("Foo-Bar"))
	assert.Equal(t, "foo_bar", server.Sanitize("Foo Bar"))
}

func TestSplit(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"1", "2", "3"}, server.Split("1,2,3"))
	assert.Equal(t, []string{"1", "2", "3"}, server.Split(" 1,2,3"))
	assert.Equal(t, []string{"1", "2", "3"}, server.Split("1,2,3 "))
	assert.Equal(t, []string{"1", "2", "3"}, server.Split("1, 2, 3"))
	assert.Equal(t, []string{"1", "2", "3"}, server.Split(" 1, 2, 3 "))
}
