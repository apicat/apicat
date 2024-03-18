package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates
var tplfs embed.FS

var templateFs = template.Must(template.New("").ParseFS(tplfs, "templates/*.tmpl"))

func createContent(templateName string, data contentData) fmt.Stringer {
	return &content{
		tplname: templateName,
		tpldata: data,
	}
}

type contentData struct {
	Link string
	Data any
}

type content struct {
	tplname string
	tpldata contentData
}

func (c *content) String() string {
	var buf bytes.Buffer
	if err := templateFs.ExecuteTemplate(&buf, c.tplname, c.tpldata); err != nil {
		return ""
	}
	return buf.String()
}
