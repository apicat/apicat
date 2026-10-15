package render

import (
	"fmt"

	"github.com/apicat/apicat/v3/internal/model"
)

type schemaTableView struct {
	Rows       []schemaRow
	Parameters bool
}

type schemaRow struct {
	*model.SchemaView
	Label        string
	Level        int
	Guides       []int
	Synthetic    bool
	ShowRequired bool
	HasChildren  bool
	In           string
}

func schemaTable(schema *model.SchemaView) schemaTableView {
	table := schemaTableView{}
	table.appendSchema(schema, "Root", 0, true)
	return table
}

func parameterTable(parameters []model.Parameter) schemaTableView {
	table := schemaTableView{Parameters: true}
	for _, parameter := range parameters {
		schema := model.SchemaView{}
		if parameter.Schema != nil {
			schema = *parameter.Schema
		}
		schema.Name = parameter.Name
		schema.Description = parameter.Description
		schema.Required = parameter.Required
		schema.Deprecated = schema.Deprecated || parameter.Deprecated
		index := len(table.Rows)
		table.appendSchema(&schema, parameter.Name, 0, false)
		table.Rows[index].In = parameter.In
	}
	return table
}

// Flatten once in preorder so every field aligns to the same table columns.
// The converter terminates recursive references and bounds schema depth.
func (table *schemaTableView) appendSchema(schema *model.SchemaView, label string, level int, synthetic bool) {
	if schema == nil {
		return
	}
	index := len(table.Rows)
	table.Rows = append(table.Rows, schemaRow{
		SchemaView: schema, Label: label, Level: level, Guides: make([]int, level),
		Synthetic: synthetic, ShowRequired: !synthetic,
	})
	if schema.Recursive || schema.Truncated {
		return
	}
	for i := range schema.Properties {
		child := &schema.Properties[i]
		table.appendSchema(child, child.Name, level+1, false)
	}
	table.appendSchema(schema.Items, "items", level+1, true)
	for _, group := range []struct {
		name     string
		children []model.SchemaView
	}{{"allOf", schema.AllOf}, {"oneOf", schema.OneOf}, {"anyOf", schema.AnyOf}} {
		for i := range group.children {
			table.appendSchema(&group.children[i], fmt.Sprintf("%s · %d", group.name, i+1), level+1, true)
		}
	}
	table.Rows[index].HasChildren = len(table.Rows) > index+1
}
