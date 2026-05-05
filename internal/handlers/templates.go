package handlers

import (
	"html/template"
	"io/fs"
)

// ParseTemplates loads every template under "templates/*.html" inside the
// provided fs.FS. Both the page and contact handlers share a single
// *template.Template so cross-file `{{template "..."}}` lookups resolve.
//
// assetVersion is appended as a `?v=` query parameter by the `asset`
// template func — pass web.AssetVersion at the call site so static
// asset URLs bust browser + CDN caches whenever the embedded files
// change.
func ParseTemplates(templates fs.FS, assetVersion string) (*template.Template, error) {
	funcs := template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"safe":  func(s string) template.HTML { return template.HTML(s) },
		"asset": func(path string) string { return path + "?v=" + assetVersion },
	}
	return template.New("").Funcs(funcs).ParseFS(templates, "templates/*.html")
}
