package web

import "embed"

// Files contains the templates and static assets needed at runtime.
//
//go:embed templates/*.html static/css/*.css
var Files embed.FS
