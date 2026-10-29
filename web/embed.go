package web

import "embed"

// Content contains the production React bundle. Rebuild it before Go releases.
//go:embed all:dist
var Content embed.FS
