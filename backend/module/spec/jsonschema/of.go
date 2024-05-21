package jsonschema

type AllOf []*Schema
type AnyOf []*Schema
type OneOf []*Schema

func (a AllOf) Merge() AllOf {
	helper := NewMergeHelper(&Schema{})
	helper.Merge(a)
	return []*Schema{helper.result}
}

type mergeHelper struct {
	result *Schema
}

func NewMergeHelper(s *Schema) *mergeHelper {
	return &mergeHelper{
		result: s,
	}
}

func (m *mergeHelper) Merge(froms []*Schema) *Schema {
	for _, from := range froms {
		if from.CheckAllOf() {
			helper := NewMergeHelper(m.result)
			helper.Merge(from.AllOf)
		} else {
			m.merge(from)
		}
	}
	return m.result
}

func (m *mergeHelper) merge(from *Schema) {
	m.mergeType(from)
	m.mergeProperties(from)
	m.mergeOthers(from)
}

func (m *mergeHelper) mergeType(from *Schema) {
	if m.result.Type == nil && from.Type != nil {
		m.result.Type = from.Type
	}
}

func (m *mergeHelper) mergeProperties(from *Schema) {
	if from.Properties == nil {
		return
	}

	if m.result.Properties == nil {
		m.result.Properties = make(map[string]*Schema)
	}

	for name, prop := range from.Properties {
		m.result.Properties[name] = prop
	}
}

func (m *mergeHelper) mergeOthers(from *Schema) {
	if len(m.result.Enum) == 0 && len(from.Enum) > 0 {
		m.result.Enum = from.Enum
	}
	if m.result.Pattern == "" && from.Pattern != "" {
		m.result.Pattern = from.Pattern
	}
	if m.result.MinLength == nil && from.MinLength != nil {
		m.result.MinLength = from.MinLength
	}
	if m.result.MaxLength == nil && from.MaxLength != nil {
		m.result.MaxLength = from.MaxLength
	}
	if m.result.ExclusiveMaximum == nil && from.ExclusiveMaximum != nil {
		m.result.ExclusiveMaximum = from.ExclusiveMaximum
	}
	if m.result.MultipleOf == nil && from.MultipleOf != nil {
		m.result.MultipleOf = from.MultipleOf
	}
	if m.result.ExclusiveMinimum == nil && from.ExclusiveMinimum != nil {
		m.result.ExclusiveMinimum = from.ExclusiveMinimum
	}
	if m.result.Maximum == nil && from.Maximum != nil {
		m.result.Maximum = from.Maximum
	}
	if m.result.Minimum == nil && from.Minimum != nil {
		m.result.Minimum = from.Minimum
	}
	if m.result.MaxProperties == nil && from.MaxProperties != nil {
		m.result.MaxProperties = from.MaxProperties
	}
	if m.result.MinProperties == nil && from.MinProperties != nil {
		m.result.MinProperties = from.MinProperties
	}
	if m.result.Required == nil && from.Required != nil {
		m.result.Required = from.Required
	}
	if m.result.MaxItems == nil && from.MaxItems != nil {
		m.result.MaxItems = from.MaxItems
	}
	if m.result.MinItems == nil && from.MinItems != nil {
		m.result.MinItems = from.MinItems
	}
	if m.result.UniqueItems == nil && from.UniqueItems != nil {
		m.result.UniqueItems = from.UniqueItems
	}
	if m.result.Format == "" && from.Format != "" {
		m.result.Format = from.Format
	}
}
