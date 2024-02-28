package fxhttpclient_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxhttpclient"
	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "foo_bar", fxhttpclient.Sanitize("foo-bar"))
	assert.Equal(t, "foo_bar", fxhttpclient.Sanitize("foo bar"))
	assert.Equal(t, "foo_bar", fxhttpclient.Sanitize("Foo-Bar"))
	assert.Equal(t, "foo_bar", fxhttpclient.Sanitize("Foo Bar"))
}
