package httpserver_test

import (
	"strings"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/stretchr/testify/assert"
)

func TestHtmlTemplateRenderer(t *testing.T) {
	t.Parallel()

	var builder strings.Builder

	renderer := httpserver.NewHtmlTemplateRenderer("testdata/templates/*.html")

	err := renderer.Render(
		&builder,
		"test.html",
		map[string]interface{}{
			"value": "some test value",
		},
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(t, "Result: some test value", builder.String())
}
