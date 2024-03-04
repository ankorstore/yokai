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

func TestFlip(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		map[string]string{},
		fxhttpclient.Flip(map[string]string{}),
	)

	m := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
	}

	assert.Equal(
		t,
		map[string]string{
			"1": "one",
			"2": "two",
			"3": "three",
		},
		fxhttpclient.Flip(m),
	)
}
