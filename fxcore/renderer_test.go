package fxcore_test

import (
	"embed"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/templates/*
var testTemplatesFS embed.FS

func TestDashboardRenderer(t *testing.T) {
	t.Parallel()

	var builder strings.Builder

	renderer := fxcore.NewDashboardRenderer(testTemplatesFS, "testdata/templates/dashboard.html")

	err := renderer.Render(
		&builder,
		"dashboard.html",
		map[string]interface{}{
			"value": "some test value",
		},
		nil,
	)
	assert.NoError(t, err)
	assert.Contains(t, "Result: some test value", builder.String())
}
