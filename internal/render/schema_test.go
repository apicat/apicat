package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/apicat/apicat/v3/internal/model"
)

func TestSchemaTableHierarchy(t *testing.T) {
	schema := &model.SchemaView{Type: "object", Properties: []model.SchemaView{
		{Name: "data", Type: "array", Required: true, Items: &model.SchemaView{
			Type: "object", Properties: []model.SchemaView{{Name: "title", Type: "string"}},
		}},
		{Name: "choice", OneOf: []model.SchemaView{{Type: "string"}, {Type: "integer"}}},
	}}
	rows := schemaTable(schema).Rows
	wantNames := []string{"Root", "data", "items", "title", "choice", "oneOf · 1", "oneOf · 2"}
	wantDepths := []int{0, 1, 2, 3, 1, 2, 2}
	if len(rows) != len(wantNames) {
		t.Fatalf("rows = %d", len(rows))
	}
	for i, row := range rows {
		if row.Label != wantNames[i] || row.Level != wantDepths[i] {
			t.Errorf("row %d: %+v", i, row)
		}
		if row.HasChildren != (i == 0 || i == 1 || i == 2 || i == 4) {
			t.Errorf("incorrect toggle for %s", row.Label)
		}
	}
	if rows[0].ShowRequired || rows[2].ShowRequired || !rows[1].ShowRequired {
		t.Error("structural rows must not imply optional properties")
	}
}

func TestParameterTableMetadataAndEscaping(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	zero := int64(0)
	parameters := []model.Parameter{{
		Name: "filter", In: "query", Description: "<script>alert(1)</script>", Required: true, Deprecated: true,
		Schema: &model.SchemaView{Type: "object", Default: "{}", Example: "<img src=x>", MinLength: &zero,
			Properties: []model.SchemaView{{Name: "value", Type: "string", Enum: []string{"a", "b"}}}},
	}, {Name: "empty", In: "header"}}
	var html bytes.Buffer
	if err := Execute(&html, "parameter_table", parameters); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Default</th>", "Example</th>", "Description</th>", "deprecated", "minLength: 0", "enum: [a, b]", "&lt;script&gt;", "&lt;img src=x&gt;", `data-depth="1"`} {
		if !strings.Contains(html.String(), want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(html.String(), "<script>") {
		t.Fatal("unescaped schema text")
	}
	if strings.Count(html.String(), "<table ") != 1 {
		t.Fatal("schema fields should share one table")
	}
}
