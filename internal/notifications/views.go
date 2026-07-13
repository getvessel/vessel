package notifications

import (
	"embed"
	"html/template"
)

//go:embed emails/*.tmpl
var emailFS embed.FS

var HTMLTemplates *template.Template

func init() {
	HTMLTemplates = template.Must(template.ParseFS(emailFS, "emails/*.tmpl"))
}
