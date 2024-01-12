package fxcore

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// DashboardRenderer is the core dashboard template renderer, based on [template.Template].
type DashboardRenderer struct {
	engine *template.Template
}

// NewDashboardRenderer returns a new DashboardRenderer.
func NewDashboardRenderer(fs embed.FS, tpl string) *DashboardRenderer {
	return &DashboardRenderer{
		engine: template.Must(template.ParseFS(fs, tpl)),
	}
}

// Render renders the core dashboard template.
func (r *DashboardRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.engine.ExecuteTemplate(w, name, data)
}
