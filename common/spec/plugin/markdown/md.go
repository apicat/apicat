package markdown

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

func Encode(in *spec.Spec) ([]byte, error) {
	paths := in.CollectionsMap(true, 2)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# %s\n", in.Info.Title)
	buf.WriteString(in.Info.Description)
	buf.WriteString("\n\n## servers\n")
	for _, v := range in.Servers {
		fmt.Fprintf(&buf, "- **%s** `%s`\n", v.Description, v.URL)
	}
	buf.WriteString("\n## apis\n")
	buildtoc(&buf, 0, in.Collections)
	for path, items := range paths {
		for method, item := range items {
			rednerHttpPart(&buf, path, method, item)
		}
	}
	return buf.Bytes(), nil
}

func buildtoc(buf *bytes.Buffer, lvl int, cs []*spec.CollectItem) {
	for _, item := range cs {
		switch item.Type {
		case spec.ContentItemTypeDir:
			if len(item.Items) > 0 {
				fmt.Fprintf(buf, "%s- %s\n", strings.Repeat("/t", lvl), item.Title)
				buildtoc(buf, lvl+1, item.Items)
			}
		case spec.ContentItemTypeDoc:
			// todo
		case spec.ContentItemTypeHttp:
			fmt.Fprintf(buf, "%s- [%s](#%d)\n", strings.Repeat("/t", lvl), item.Title, item.ID)
		}
	}
}

var jsonschemaHeaderCols = []string{"name", "type", "required", "comment"}
var paramsHeaderCols = []string{"name", "in", "type", "required", "comment"}

func rednerHttpPart(buf *bytes.Buffer, path, method string, part spec.HTTPPart) {
	fmt.Fprintf(buf, "### <span id=\"%d\">%s</span>\n", part.ID, part.Title)
	fmt.Fprintf(buf, "- method: **`%s`**\n", method)
	fmt.Fprintf(buf, "- path: **`%s`**\n", path)
	fmt.Fprintf(buf, "\n>%s\n", "")
	fmt.Fprintf(buf, "\n**parameters**\n\n")
	renderTableHeader(buf, paramsHeaderCols)
	for k, v := range part.Parameters.Map() {
		for _, item := range v {
			buf.WriteString("| **")
			renderString(buf, item.Name)
			buf.WriteString("** | ")
			buf.WriteString(k)
			buf.WriteString(" | ")
			buf.WriteString(item.Schema.Type.Value()[0])
			buf.WriteString(" | ")
			if item.Required {
				buf.WriteString("yes")
			}
			buf.WriteString(" | ")
			renderString(buf, item.Schema.Description)
			buf.WriteByte('\n')
		}
	}
	fmt.Fprintf(buf, "\n**body**\n\n")
	for k, v := range part.Content {
		fmt.Fprintf(buf, "contentType `%s`\n", k)
		renderTableHeader(buf, jsonschemaHeaderCols)
		renderSchema(buf, "`root`", 0, true, v.Schema)
		if strings.Contains(k, "json") {
			b, _ := json.Marshal(v.Schema)
			if rx, err := datagen.JSONSchemaGen(b, &datagen.GenOption{DatagenKey: "x-apicat-mock"}); err == nil {
				buf.WriteString("example\n\n")
				buf.WriteString("\n```json\n")
				mockexample, _ := json.MarshalIndent(rx, "", "  ")
				buf.Write(mockexample)
				buf.WriteString("\n```\n\n")
			}
		}
	}

	fmt.Fprintf(buf, "\n**response**\n\n")
	for _, res := range part.Responses {
		fmt.Fprintf(buf, "%s\n\n", res.Name)
		fmt.Fprintf(buf, "- statusCode: `%d`\n", res.Code)
		for k, v := range res.Content {
			fmt.Fprintf(buf, "- contentType `%s`\n\n", k)
			renderTableHeader(buf, jsonschemaHeaderCols)
			renderSchema(buf, "`root`", 0, true, v.Schema)

			if strings.Contains(k, "json") {
				b, _ := json.Marshal(v.Schema)
				if rx, err := datagen.JSONSchemaGen(b, &datagen.GenOption{DatagenKey: "x-apicat-mock"}); err == nil {
					buf.WriteString("example\n\n")
					buf.WriteString("\n```json\n")
					mockexample, _ := json.MarshalIndent(rx, "", "  ")
					buf.Write(mockexample)
					buf.WriteString("\n```\n\n")
				}
			}
			break
		}
	}

	buf.WriteString("\n------------\n")
}

func renderSchema(buf *bytes.Buffer, name string, lvl int, required bool, s *jsonschema.Schema) {
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
		renderSchema(buf, "`item`", lvl, required, s.Items.Value())
	}
}

func renderSchemaItem(buf *bytes.Buffer, name, typ, desc string, lvl int, required bool) {
	buf.WriteByte('|')
	buf.WriteString(strings.Repeat("·", lvl*4))
	buf.WriteString("**")
	renderString(buf, name)
	buf.WriteString("**")
	buf.WriteByte('|')
	buf.WriteByte('`')
	buf.WriteString(typ)
	buf.WriteByte('`')
	buf.WriteByte('|')
	if required {
		buf.WriteString("yes")
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
