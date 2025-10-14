package templates

import "embed"

// FS exposes embedded HTML templates for the server package.
//
//go:embed *.tmpl
var FS embed.FS
