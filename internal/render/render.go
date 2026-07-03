package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/apicat/apicat/v3/web"
	"github.com/yuin/goldmark"
)

var templates *template.Template

var md = goldmark.New()

func Init() error {
	funcMap := template.FuncMap{
		"methodLower": func(s string) string {
			s = strings.ToLower(s)
			if s == "delete" {
				return "del"
			}
			return s
		},
		"statusClass": statusClass,
		"joinEnum":    joinEnum,
		"add":         func(a, b int) int { return a + b },
		"depthClass": func(d int) string {
			return fmt.Sprintf("depth-%d", d)
		},
		"schemaTypeDisplay": schemaTypeDisplay,
		"markdown": func(s string) template.HTML {
			var buf bytes.Buffer
			if err := md.Convert([]byte(s), &buf); err != nil {
				return template.HTML(template.HTMLEscapeString(s))
			}
			return template.HTML(buf.String())
		},
	}

	var err error
	templates, err = template.New("").Funcs(funcMap).ParseFS(
		web.Content,
		"templates/*.html",
		"templates/partials/*.html",
	)
	return err
}

func Execute(w io.Writer, name string, data any) error {
	return templates.ExecuteTemplate(w, name, data)
}

func statusClass(code string) string {
	if code == "default" {
		return "status-default"
	}
	if len(code) > 0 {
		switch code[0] {
		case '2':
			return "status-2xx"
		case '3':
			return "status-3xx"
		case '4':
			return "status-4xx"
		case '5':
			return "status-5xx"
		}
	}
	return "status-default"
}

func joinEnum(values []string) string {
	return strings.Join(values, ", ")
}

func schemaTypeDisplay(typ, format string) string {
	if format != "" {
		return fmt.Sprintf("%s<%s>", typ, format)
	}
	return typ
}
