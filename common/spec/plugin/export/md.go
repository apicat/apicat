package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apicat/apicat/common/spec"
	"github.com/apicat/apicat/common/spec/jsonschema"
	"github.com/apicat/datagen"
	"golang.org/x/exp/slices"
)

func Markdown(in *spec.Spec) ([]byte, error) {
	paths := in.CollectionsMap(true, 2)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# %s\n", in.Info.Title)
	fmt.Fprintf(&buf, "%s\n\nversion: `%s`\n\n", in.Info.Description, in.Info.Version)

	if len(in.Servers) > 0 {
		buf.WriteString("\n\n## Servers\n")
		for _, v := range in.Servers {
			fmt.Fprintf(&buf, "- **%s** `%s`\n", v.Description, v.URL)
		}
	}

	list := make([][2]string, 0)
	for path, items := range paths {
		for method := range items {
			list = append(list, [2]string{path, method})
		}
	}

	buf.WriteString("\n## Table of APIs\n")
	for k, v := range list {
		item := paths[v[0]][v[1]]
		fmt.Fprintf(&buf, "%d. [%s](#api-%d)\n", k+1, item.Title, k+1)
	}

	buf.WriteString("\n\n")

	for k, v := range list {
		item := paths[v[0]][v[1]]
		rednerHttpPart(&buf, k+1, v[0], v[1], item)
	}

	return buf.Bytes(), nil
}

var jsonschemaHeaderCols = []string{"name", "type", "required", "comment"}
var paramsHeaderCols = []string{"name", "in", "type", "required", "comment"}

func rednerHttpPart(buf *bytes.Buffer, i int, path, method string, part spec.HTTPPart) {
	fmt.Fprintf(buf, "## <span id=\"api-%d\">%d. %s</span>\n", i, i, part.Title)
	fmt.Fprintf(buf, "### Path\n [%s](%s)\n", path, path)
	fmt.Fprintf(buf, "### Method\n %s\n", strings.ToUpper(method))
	fmt.Fprintf(buf, "### Parameters\n")
	renderTableHeader(buf, paramsHeaderCols)
	for k, v := range part.Parameters.Map() {
		for _, item := range v {
			buf.WriteString("|")
			renderString(buf, item.Name)
			buf.WriteString("|**")
			buf.WriteString(k)
			buf.WriteString("**|`")
			buf.WriteString(item.Schema.Type.Value()[0])
			buf.WriteString("` | ")
			if item.Required {
				buf.WriteString("*")
			}
			buf.WriteString(" | ")
			renderString(buf, item.Schema.Description)
			buf.WriteByte('\n')
		}
	}
	fmt.Fprintf(buf, "### Request Body\n")
	for k, v := range part.Content {
		fmt.Fprintf(buf, "ContentType `%s`\n", k)
		renderTableHeader(buf, jsonschemaHeaderCols)
		if v.Schema == nil {
			continue
		}
		renderSchema(buf, "`root`", 0, true, v.Schema)
		if strings.Contains(k, "json") {
			b, _ := json.Marshal(v.Schema)
			if rx, err := datagen.JSONSchemaGen(b, &datagen.GenOption{DatagenKey: "x-apicat-mock"}); err == nil {
				buf.WriteString("\n\nExample\n\n")
				buf.WriteString("\n```json\n")
				mockexample, _ := json.MarshalIndent(rx, "", "  ")
				buf.Write(mockexample)
				buf.WriteString("\n```\n\n")
			}
		}
	}

	fmt.Fprintf(buf, "### Responses\n")
	for _, res := range part.Responses {
		fmt.Fprintf(buf, "StatusCode `%d` \n> %s\n\n", res.Code, res.Description)
		for k, v := range res.Content {
			fmt.Fprintf(buf, " ContentType `%s`\n\n", k)
			renderTableHeader(buf, jsonschemaHeaderCols)
			if v.Schema == nil {
				continue
			}
			renderSchema(buf, "`root`", 0, true, v.Schema)
			if strings.Contains(k, "json") {
				b, _ := json.Marshal(v.Schema)
				if rx, err := datagen.JSONSchemaGen(b, &datagen.GenOption{DatagenKey: "x-apicat-mock"}); err == nil {
					buf.WriteString("\n\nExample\n\n")
					buf.WriteString("\n```json\n")
					mockexample, _ := json.MarshalIndent(rx, "", "  ")
					buf.Write(mockexample)
					buf.WriteString("\n```\n\n")
				}
			}
			break
		}
	}

	buf.WriteString("\n\n------------\n")
}

func renderSchema(buf *bytes.Buffer, name string, lvl int, required bool, s *jsonschema.Schema) {
	if s.Type == nil {
		return
	}
	typ := s.Type.Value()
	if len(typ) > 1 {
		renderSchemaItem(buf, name, "any", s.Description, lvl, required)
		return
	}
	renderSchemaItem(buf, name, typ[0], s.Description, lvl, required)
	switch typ[0] {
	case "object":
		for k, v := range s.Properties {
			renderSchema(buf, k, lvl+1, slices.Contains(s.Required, k), v)
		}
	case "array":
		renderSchema(buf, "`item`", lvl+1, required, s.Items.Value())
	}
}

func renderSchemaItem(buf *bytes.Buffer, name, typ, desc string, lvl int, required bool) {
	buf.WriteByte('|')
	buf.WriteString(strings.Repeat("·", lvl*4))
	buf.WriteString(name)
	buf.WriteByte('|')
	buf.WriteByte('`')
	buf.WriteString(typ)
	buf.WriteByte('`')
	buf.WriteByte('|')
	if required {
		buf.WriteString("*")
	}
	buf.WriteByte('|')
	renderString(buf, desc)
	buf.WriteByte('\n')
}

func renderTableHeader(buf *bytes.Buffer, fileds []string) {
	for _, v := range fileds {
		buf.WriteString("| ")
		renderString(buf, v)
	}
	buf.WriteByte('\n')
	buf.WriteString(strings.Repeat("|:---", len(fileds)))
	buf.WriteByte('\n')
}

func renderString(buf *bytes.Buffer, s string) {
	for _, b := range s {
		if needsEscaping(b) {
			buf.WriteByte('\\')
		}
		buf.WriteRune(b)
	}
}

func needsEscaping(text rune) bool {
	switch text {
	case '\\',
		'`',
		'*',
		'_',
		'{', '}',
		'[', ']',
		'(', ')',
		'#',
		'+',
		'-':
		return true
	case '!':
		return false
	case '<', '>':
		return true
	default:
		return false
	}
}
