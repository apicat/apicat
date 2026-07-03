package web

import "embed"

//go:embed templates templates/partials static/css static/js
var Content embed.FS
