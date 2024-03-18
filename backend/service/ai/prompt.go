package ai

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed prompts
var tplfs embed.FS

var templateFs = template.Must(template.New("").ParseFS(tplfs, "prompts/*.tmpl"))

type contentData struct {
	Lang string
	Data any
}

type content struct {
	tplname string
	tpldata contentData
}

func createContent(templateName string, data contentData) *content {
	return &content{
		tplname: templateName,
		tpldata: data,
	}
}

func (c *content) String() (string, error) {
	var buf bytes.Buffer
	if err := templateFs.ExecuteTemplate(&buf, c.tplname, c.tpldata); err != nil {
		return "", err
	}
	return buf.String(), nil
}
