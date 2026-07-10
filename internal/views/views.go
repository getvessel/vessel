package views

import (
	"embed"
	"html/template"
)

//go:embed emails/*.html
var emailFS embed.FS

var (
	HTMLTemplates *template.Template
)

func init() {
	HTMLTemplates = template.Must(template.ParseFS(emailFS, "emails/*.html"))
}
