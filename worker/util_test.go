package worker_test

import (
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "foo_bar", worker.Sanitize("foo-bar"))
	assert.Equal(t, "foo_bar", worker.Sanitize("foo bar"))
	assert.Equal(t, "foo_bar", worker.Sanitize("Foo-Bar"))
	assert.Equal(t, "foo_bar", worker.Sanitize("Foo Bar"))
}
