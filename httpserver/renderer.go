package httpserver

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// HtmlTemplateRenderer allows to render HTML templates, based on [html/template].
//
// [html/template]: https://pkg.go.dev/html/template
type HtmlTemplateRenderer struct {
	engine *template.Template
}

// NewHtmlTemplateRenderer returns a [HtmlTemplateRenderer], for a file pattern.
func NewHtmlTemplateRenderer(pattern string) *HtmlTemplateRenderer {
	return &HtmlTemplateRenderer{
		engine: template.Must(template.ParseGlob(pattern)),
	}
}

// Render executes a named template, with provided data, and write the result to the provided [io.Writer].
func (r *HtmlTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.engine.ExecuteTemplate(w, name, data)
}
