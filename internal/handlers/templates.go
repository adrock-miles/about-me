package handlers

import (
	"html/template"
	"io/fs"
)

// ParseTemplates loads every template under "templates/*.html" inside the
// provided fs.FS. Both the page and contact handlers share a single
// *template.Template so cross-file `{{template "..."}}` lookups resolve.
func ParseTemplates(templates fs.FS) (*template.Template, error) {
	funcs := template.FuncMap{
		"add":  func(a, b int) int { return a + b },
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}
	return template.New("").Funcs(funcs).ParseFS(templates, "templates/*.html")
}
