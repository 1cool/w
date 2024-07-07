package template

import "embed"

//go:embed tmpl/*.tmpl
var TemplateDir embed.FS
